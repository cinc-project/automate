package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/chef/automate/api/external/sso"
	deployment "github.com/chef/automate/api/interservice/deployment"
	license_control "github.com/chef/automate/api/interservice/license_control"
	"github.com/golang/protobuf/ptypes/empty"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// SsoConfig - the ssoconfig service data structure
type SsoConfig struct {
	license_client license_control.LicenseControlServiceClient
	client         deployment.DeploymentClient
}

type PostConfig struct {
    CaContents         string  	`json:"CaContents"`
	SsoUrl             string	`json:"SsoUrl"`
	EmailAttr          string	`json:"EmailAttr"`
	UsernameAttr       string	`json:"UsernameAttr"`
	GroupsAttr         string	`json:"GroupAttr"`
	AllowedGroups      []string	`json:"AllowedGroup"`
	EntityIssuer       string	`json:"EntityIssuer"`
	NameIdPolicyFormat string	`json:"NameIdPolicyFormat"`
}
// NewSsoConfigHandler - create a new ssoconfig service handler
func NewSsoConfigHandler(license_client license_control.LicenseControlServiceClient, client deployment.DeploymentClient) *SsoConfig {
	return &SsoConfig{
		license_client: license_client,
		client:         client,
	}
}

func (a *SsoConfig) GetSsoConfig(ctx context.Context, in *empty.Empty) (*sso.GetSsoConfigResponse, error) {

	err := a.validateDeploymentType(ctx)
	if err != nil {
		return nil, err
	}


	req := &deployment.GetAutomateConfigRequest{}

	res, err := a.client.GetAutomateConfig(ctx, req)
	if err != nil {
		return nil, err
	}

	if res.Config.Dex != nil {
		ssoConfig := res.Config.Dex.V1.Sys.Connectors.Saml

		return &sso.GetSsoConfigResponse{
			CaContents:         ssoConfig.CaContents.GetValue(),
			SsoUrl:             ssoConfig.SsoUrl.GetValue(),
			EmailAttr:          ssoConfig.EmailAttr.GetValue(),
			UsernameAttr:       ssoConfig.UsernameAttr.GetValue(),
			GroupsAttr:         ssoConfig.GroupsAttr.GetValue(),
			AllowedGroups:      ssoConfig.AllowedGroups,
			EntityIssuer:       ssoConfig.EntityIssuer.GetValue(),
			NameIdPolicyFormat: ssoConfig.NameIdPolicyFormat.GetValue(),
		}, nil
	}

	return &sso.GetSsoConfigResponse{
		CaContents:         "",
		SsoUrl:             "",
		EmailAttr:          "",
		UsernameAttr:       "",
		GroupsAttr:         "",
		AllowedGroups:      []string{},
		EntityIssuer:       "",
		NameIdPolicyFormat: "",
	}, nil
}

func (a *SsoConfig) getDeploymentDetails(ctx context.Context) (string, error) {
	deployIDResponse, err := a.license_client.GetDeploymentID(ctx, &license_control.GetDeploymentIDRequest{})
	if err != nil {
		if isServiceDownError(err) {
			log.WithFields(log.Fields{
				"err":  err,
				"func": "getDeploymentDetails",
			}).Error("connecting to the license client")
			return "", nil
		}
		return "", err
	}

	log.Debugf("deployIDResponse.DeploymentType: %s ", deployIDResponse.DeploymentType)

	return deployIDResponse.DeploymentType, nil
}

func(a *SsoConfig) validateDeploymentType(ctx context.Context) error {
	deploymentType, err := a.getDeploymentDetails(ctx)
	if err != nil {
		return err
	}

	if deploymentType != "SAAS" {
		msg := "Unauthorized: Deployment type is not SAAS"
		return status.Error(codes.PermissionDenied, msg)
	}
	return nil
}

func(a *SsoConfig) SetSsoConfig(ctx context.Context, in *sso.SetSsoConfigRequest) (*sso.SetSsoConfigResponse , error) {
	err := a.validateDeploymentType(ctx)
	if err != nil {
		return nil , err
	}

	req := &sso.SetSsoConfigRequest{
		CaContents: in.CaContents,
		SsoUrl: in.SsoUrl,
		EmailAttr: in.EmailAttr,
		UsernameAttr: in.UsernameAttr,
		GroupsAttr: in.GroupsAttr,
		AllowedGroups: in.AllowedGroups,
		EntityIssuer: in.EntityIssuer,
		NameIdPolicyFormat: in.NameIdPolicyFormat,
	}
	body_params:= &PostConfig{
		CaContents:         req.CaContents,
		SsoUrl:             req.SsoUrl,
		EmailAttr:          req.EmailAttr,
		UsernameAttr:       req.UsernameAttr,
		GroupsAttr:         req.GroupsAttr,
		AllowedGroups:      req.AllowedGroups,
		EntityIssuer:       req.EntityIssuer,
		NameIdPolicyFormat: req.NameIdPolicyFormat,
	}
	// jsonValue, _ :=  json.Marshal(body_params)
	buf := new(bytes.Buffer)
	ip := getBastionIp()
	url := "http://" + string(ip)
	json.NewEncoder(buf).Encode(body_params)
	
	request, _ := http.NewRequest("POST",url,buf)
	client := &http.Client{}
	res, err := client.Do(request)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}
	fmt.Println("response Body:", string(body))
	return &sso.SetSsoConfigResponse{
		Response: "Config patch was successfull",
	}, nil
	

}

func getBastionIp() string {
	content, err := ioutil.ReadFile("/var/automate-ha/bastion_info.txt")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(content))
	return strings.TrimSpace(string(content))
}