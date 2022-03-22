package server

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"

	chef "github.com/go-chef/chef"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	secrets "github.com/chef/automate/api/external/secrets"
	"github.com/chef/automate/api/interservice/infra_proxy/request"
	"github.com/chef/automate/api/interservice/infra_proxy/response"
	"github.com/chef/automate/components/infra-proxy-service/service"
	"github.com/chef/automate/components/infra-proxy-service/storage"
	"github.com/chef/automate/components/infra-proxy-service/validation"
)

//GetAutomateInfraServerUsersList: Fetches the list of automate infra server users from the DB
func (s *Server) GetAutomateInfraServerUsersList(ctx context.Context, req *request.AutomateInfraServerUsers) (*response.AutomateInfraServerUsers, error) {
	// Check server exists in automate or not
	err := s.isServerExistInAutomate(ctx, req.ServerId)
	if err != nil {
		return nil, service.ParseStorageError(err, req.ServerId, "server")
	}

	usersList, err := s.service.Storage.GetAutomateInfraServerUsers(ctx, req.ServerId)
	if err != nil {
		return nil, service.ParseStorageError(err, *req, "user")
	}

	return &response.AutomateInfraServerUsers{
		Users: fromStorageToListAutomateInfraServerUsers(usersList),
	}, nil
}

// fromStorageAutomateInfraServerUser: Create a response.AutomateInfraServerUsersListItem from a storage.User
func fromStorageAutomateInfraServerUser(u storage.User) *response.AutomateInfraServerUsersListItem {
	return &response.AutomateInfraServerUsersListItem{
		Id:                  u.ID,
		ServerId:            u.ServerID,
		InfraServerUsername: u.InfraServerUsername,
		Connector:           u.Connector,
		AutomateUserId:      u.AutomateUserID,
		IsServerAdmin:       u.IsServerAdmin,
	}
}

//fromStorageToListAutomateInfraServerUsers: Create a response.AutomateInfraServerUsersListItem from an array of storage.User
func fromStorageToListAutomateInfraServerUsers(ul []storage.User) []*response.AutomateInfraServerUsersListItem {
	tl := make([]*response.AutomateInfraServerUsersListItem, len(ul))

	for i, user := range ul {
		tl[i] = fromStorageAutomateInfraServerUser(user)
	}

	return tl
}

//isServerExistInAutomate: Check whether server exist in automate or not
func (s *Server) isServerExistInAutomate(ctx context.Context, serverId string) error {
	_, err := s.service.Storage.GetServer(ctx, serverId)
	if err != nil {
		return err
	}
	return nil
}

// ResetInfraServerUserKey updates the public key on the Chef Server and returns the private key
func (s *Server) ResetInfraServerUserKey(ctx context.Context, req *request.ResetInfraServerUserKeyReq) (*response.ResetInfraServerUserKeyRes, error) {
	err := validation.New(validation.Options{
		Target:  "user",
		Request: *req,
		Rules: validation.Rules{
			"ServerId": []string{"required"},
			"Name":     []string{"required"},
		},
	}).Validate()

	if err != nil {
		return nil, err
	}

	server, err := s.service.Storage.GetServer(ctx, req.ServerId)
	if err != nil {
		return nil, err
	}
	if server.CredentialID == "" {
		return nil, errors.New("webui key is not available with server")
	}
	// Get web ui key from secrets service
	secret, err := s.service.Secrets.Read(ctx, &secrets.Id{Id: server.CredentialID})
	if err != nil {
		return nil, err
	}
	c, err := s.createChefServerClient(ctx, req.ServerId, GetAdminKeyFrom(secret), "pivotal", true)
	if err != nil {
		return nil, err
	}

	key := "default"

	// Deletes the existing key
	_, err = c.client.Users.DeleteKey(req.UserName, key)
	chefError, _ := chef.ChefError(err)
	if err != nil && chefError.StatusCode() != 404 {
		return nil, ParseAPIError(err)
	}

	// Add new key to existing client
	body, err := chef.JSONReader(AccessKeyReq{
		Name:           key,
		ExpirationDate: "infinity",
		CreateKey:      true,
	})
	if err != nil {
		return nil, ParseAPIError(err)
	}

	var chefKey chef.ChefKey
	addReq, err := c.client.NewRequest("POST", fmt.Sprintf("users/%s/keys", req.UserName), body)

	if err != nil {
		return nil, ParseAPIError(err)
	}

	res, err := c.client.Do(addReq, &chefKey)
	if res != nil {
		defer res.Body.Close() //nolint:errcheck
	}

	if err != nil {
		return nil, ParseAPIError(err)
	}

	// pubKey, privKey, err := GenerateKeys()
	// if err != nil {
	// 	return nil, err
	// }

	// // Update the existing key
	// _, err = client.client.Users.UpdateKey(req.UserName, "", chef.AccessKey{
	// 	Name:           req.UserName,
	// 	PublicKey:      pubKey,
	// 	ExpirationDate: "infinity",
	// })

	return &response.ResetInfraServerUserKeyRes{
		PrivateKey: chefKey.PrivateKey,
		UserName:   req.UserName,
		ServerId:   req.ServerId,
	}, nil
}

func GenerateKeys() (string, string, error) {
	// generate key
	privatekey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Errorf("Cannot generate RSA key\n")
		return "", "", err
	}
	publicKey := &privatekey.PublicKey

	// Encode Private Key into a string variable
	var privateKeyBytes []byte = x509.MarshalPKCS1PrivateKey(privatekey)
	privateKeyBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	}
	privateKeyByteSlice := pem.EncodeToMemory(privateKeyBlock)
	if err != nil {
		log.Errorf("error when encode private key: %s \n", err)
		return "", "", err
	}

	// Encode Public Key into a string variable
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		log.Errorf("error when dumping publickey: %s \n", err)
		return "", "", err
	}
	publicKeyBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	}

	publicKeyByteSlice := pem.EncodeToMemory(publicKeyBlock)
	if err != nil {
		log.Errorf("error when encode public key: %s \n", err)
		return "", "", err
	}
	return string(publicKeyByteSlice), string(privateKeyByteSlice), nil
}
