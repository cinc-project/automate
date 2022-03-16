package server

import (
	"context"
	"fmt"

	"github.com/chef/automate/api/interservice/infra_proxy/request"
	"github.com/chef/automate/api/interservice/infra_proxy/response"
	"github.com/chef/automate/components/infra-proxy-service/service"
	"github.com/chef/automate/components/infra-proxy-service/storage"
	"github.com/chef/automate/components/infra-proxy-service/validation"
	chef "github.com/go-chef/chef"
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

func (s *Server) ResetInfraServerUserKey(ctx context.Context, req *request.ResetInfraServerUserKeyReq) (*response.ResetInfraServerUserKeyRes, error) {
	err := validation.New(validation.Options{
		Target:  "client",
		Request: *req,
		Rules: validation.Rules{
			"OrgId":    []string{"required"},
			"ServerId": []string{"required"},
			"Name":     []string{"required"},
		},
	}).Validate()

	if err != nil {
		return nil, err
	}

	c, err := s.createClient(ctx, req.OrgId, req.ServerId)
	if err != nil {
		return nil, err
	}

	key := req.Key

	if key == "" {
		key = "default"
	}

	// Deletes the existing key
	_, err = c.client.Clients.DeleteKey(req.Name, key)
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
	addReq, err := c.client.NewRequest("POST", fmt.Sprintf("clients/%s/keys", req.Name), body)

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

	return &response.ResetClient{
		Name: req.Name,
		ClientKey: &response.ClientKey{
			Name:           key,
			PublicKey:      chefKey.PublicKey,
			ExpirationDate: chefKey.ExpirationDate,
			PrivateKey:     chefKey.PrivateKey,
		},
	}, nil
	// Either Cll Reset Key and get the private key or grnerate the key and and send public key to chef and privste to automate
	return &response.ResetInfraServerUserKeyRes{}, nil
}
