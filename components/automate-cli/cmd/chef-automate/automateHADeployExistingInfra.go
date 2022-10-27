package main

import (
	"container/list"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/chef/automate/components/automate-cli/pkg/status"
	ptoml "github.com/pelletier/go-toml"
)

type existingInfra struct {
	config     ExistingInfraConfigToml
	configPath string
}

func newExistingInfa(configPath string) *existingInfra {
	return &existingInfra{
		configPath: configPath,
	}
}

func (e *existingInfra) doDeployWork(args []string) error {
	var err = bootstrapEnv(e, deployCmdFlags.airgap, false)
	if err != nil {
		return err
	}
	err = executeSecretsInitCommand(e.config.Architecture.ConfigInitials.SecretsKeyFile)
	if err != nil {
		return err
	}
	return executeDeployment(args)
}

func (e *existingInfra) doProvisionJob(args []string) error {
	return nil
}

func (e *existingInfra) generateConfig() error {
	templateBytes, err := ioutil.ReadFile(e.getConfigPath())
	if err != nil {
		return status.Wrap(err, status.FileAccessError, "error in reading config toml file")
	}
	e.config = ExistingInfraConfigToml{}
	err = ptoml.Unmarshal(templateBytes, &e.config)
	if err != nil {
		return status.Wrap(err, status.ConfigError, "error in unmarshalling config toml file")
	}
	errList := e.validateConfigFields()
	if errList != nil && errList.Len() > 0 {
		return status.Wrap(getSingleErrorFromList(errList), status.ConfigError, "config is invalid.")
	}

	//If CustomCertsEnabled for OpenSearch is enabled, then get admin_dn and nodes_dn from the certs
	if e.config.Opensearch.Config.CustomCertsEnabled {
		//If AdminCert is given then get the admin_dn from the cert
		if len(strings.TrimSpace(e.config.Opensearch.Config.AdminCert)) > 0 {
			admin_dn, err := e.getDistinguishedNameFromKey(e.config.Opensearch.Config.AdminCert)
			if err != nil {
				return err
			}
			e.config.Opensearch.Config.AdminDn = admin_dn
		}
		//If PublicKey is given then get the nodes_dn from the cert
		if len(strings.TrimSpace(e.config.Opensearch.Config.PublicKey)) > 0 {
			nodes_dn, err := e.getDistinguishedNameFromKey(e.config.Opensearch.Config.PublicKey)
			if err != nil {
				return err
			}
			e.config.Opensearch.Config.NodesDn = nodes_dn
		}

		//Set the admin_dn and nodes_dn in the config for all IP addresses
		for i := 0; i < len(e.config.Opensearch.Config.CertsByIP); i++ {
			//If AdminCert is given then get the admin_dn from the cert
			adminCert := e.config.Opensearch.Config.CertsByIP[i].AdminCert
			if len(strings.TrimSpace(adminCert)) > 0 {
				admin_dn, err := e.getDistinguishedNameFromKey(adminCert)
				if err != nil {
					return err
				}
				e.config.Opensearch.Config.CertsByIP[i].AdminDn = admin_dn
			}

			//If PublicKey is given then get the nodes_dn from the cert
			publicKey := e.config.Opensearch.Config.CertsByIP[i].PublicKey
			if len(strings.TrimSpace(publicKey)) > 0 {
				nodes_dn, err := e.getDistinguishedNameFromKey(publicKey)
				if err != nil {
					return err
				}
				e.config.Opensearch.Config.CertsByIP[i].NodesDn = nodes_dn
			}
		}
	}

	finalTemplate := renderSettingsToA2HARBFile(existingNodesA2harbTemplate, e.config)
	writeToA2HARBFile(finalTemplate, initConfigHabA2HAPathFlag.a2haDirPath+"a2ha.rb")
	return nil
}

func (e *existingInfra) getDistinguishedNameFromKey(publicKey string) (string, error) {
	block, _ := pem.Decode([]byte(publicKey))
	if block == nil {
		return "", status.New(status.ConfigError, "failed to decode certificate PEM")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return "", status.Wrap(err, status.ConfigError, "failed to parse certificate PEM")
	}
	return fmt.Sprintf("%v", cert.Subject), nil
}

func (e *existingInfra) getConfigPath() string {
	return e.configPath
}

func (e *existingInfra) validateConfigFields() *list.List {
	errorList := list.New()
	if len(e.config.Architecture.ConfigInitials.SecretsKeyFile) < 1 {
		errorList.PushBack("Invalid or empty secrets_key_file")
	}
	if len(e.config.Architecture.ConfigInitials.SecretsStoreFile) < 1 {
		errorList.PushBack("Invalid or empty secrets_store_file")
	}
	if len(e.config.Architecture.ConfigInitials.Architecture) < 1 {
		errorList.PushBack("Invalid or empty Architecture")
	}
	if len(e.config.Architecture.ConfigInitials.WorkspacePath) < 1 {
		errorList.PushBack("Invalid or empty workspace_path")
	}
	if len(e.config.Architecture.ConfigInitials.SSHUser) < 1 {
		errorList.PushBack("Invalid or empty ssh_user")
	}
	if len(e.config.Architecture.ConfigInitials.SSHKeyFile) < 1 {
		errorList.PushBack("Invalid or empty ssh_key_file")
	}
	if len(e.config.Architecture.ConfigInitials.BackupMount) < 1 {
		errorList.PushBack("Invalid or empty backup_mount")
	}
	if len(e.config.Automate.Config.InstanceCount) < 1 {
		errorList.PushBack("Invalid or empty automate instance_count")
	}
	if len(e.config.ChefServer.Config.InstanceCount) < 1 {
		errorList.PushBack("Invalid or empty chef-server instance_count")
	}
	if len(e.config.Opensearch.Config.InstanceCount) < 1 {
		errorList.PushBack("Invalid or empty open-search instance_count")
	}
	if len(e.config.Postgresql.Config.InstanceCount) < 1 {
		errorList.PushBack("Invalid or empty postgres-sql instance_count")
	}

	if len(e.config.ExistingInfra.Config.AutomatePrivateIps) < 1 {
		errorList.PushBack("Invalid or empty automate_private_ips")
	}

	if len(e.config.ExistingInfra.Config.ChefServerPrivateIps) < 1 {
		errorList.PushBack("Invalid or empty chef_server_private_ips")
	}

	if len(e.config.ExistingInfra.Config.OpensearchPrivateIps) < 1 {
		errorList.PushBack("Invalid or empty opensearch_private_ips")
	}

	if len(e.config.ExistingInfra.Config.PostgresqlPrivateIps) < 1 {
		errorList.PushBack("Invalid or empty postgresql_private_ips")
	}

	if len(e.config.Architecture.ConfigInitials.BackupConfig) > 0 {
		if e.config.Architecture.ConfigInitials.BackupConfig == "object_storage" {
			if len(e.config.ObjectStorage.Config.AccessKey) < 1 {
				errorList.PushBack("Invalid or empty access_key")
			}
			if len(e.config.ObjectStorage.Config.SecretKey) < 1 {
				errorList.PushBack("Invalid or empty secret_key")
			}
			if len(e.config.ObjectStorage.Config.BucketName) < 1 {
				errorList.PushBack("Invalid or empty bucket_name")
			}
			if len(e.config.ObjectStorage.Config.Endpoint) < 1 {
				errorList.PushBack("Invalid or empty endpoint")
			}
		} else if e.config.Architecture.ConfigInitials.BackupConfig == "file_system" {
			// if len(e.config.ObjectStorage.Config.AccessKey) < 1 {
			// 	errorList.PushBack("Invalid or empty access_key")
			// }
		} else {
			errorList.PushBack("Invalid or empty backup_config")
		}
	}

	errorList.PushBackList(e.validateIPs())
	errorList.PushBackList(e.validateCerts())
	return errorList
}

func (e *existingInfra) validateCerts() *list.List {

	errorList := list.New()

	// if CustomCertsEnabled is disabled, then skip validation for custom certs and use self signed certs
	if e.config.Automate.Config.CustomCertsEnabled {
		// check if all the default certs are given
		if len(strings.TrimSpace(e.config.Automate.Config.RootCA)) > 0 ||
			len(strings.TrimSpace(e.config.Automate.Config.PrivateKey)) > 0 ||
			len(strings.TrimSpace(e.config.Automate.Config.PublicKey)) > 0 {
			if len(strings.TrimSpace(e.config.Automate.Config.RootCA)) < 1 ||
				len(strings.TrimSpace(e.config.Automate.Config.PrivateKey)) < 1 ||
				len(strings.TrimSpace(e.config.Automate.Config.PublicKey)) < 1 {
				errorList.PushBack("Automate RootCA and/or Public Key and/or Private Key are missing. Set custom_certs_enabled to false to continue without custom certificates.")
			}
		}
		// check if all the certs are valid for given IPs
		for _, node := range e.config.Automate.Config.CertsByIP {
			if len(strings.TrimSpace(node.RootCA)) < 1 ||
				len(strings.TrimSpace(node.PrivateKey)) < 1 ||
				len(strings.TrimSpace(node.PublicKey)) < 1 {
				errorList.PushBack("Automate RootCA and/or Public Key and/or Private Key are missing. Set custom_certs_enabled to false to continue without custom certificates.")
			}
		}
	}

	// if CustomCertsEnabled is disabled, then skip validation for custom certs and use self signed certs
	if e.config.ChefServer.Config.CustomCertsEnabled {
		// check if all the default certs are given
		if len(strings.TrimSpace(e.config.ChefServer.Config.PrivateKey)) > 0 ||
			len(strings.TrimSpace(e.config.ChefServer.Config.PublicKey)) > 0 {
			if len(strings.TrimSpace(e.config.ChefServer.Config.PrivateKey)) < 1 ||
				len(strings.TrimSpace(e.config.ChefServer.Config.PublicKey)) < 1 {
				errorList.PushBack("ChefServer Public Key and/or Private Key are missing. Set custom_certs_enabled to false to continue without custom certificates.")
			}
		}
		// check if all the certs are valid for given IPs
		for _, node := range e.config.ChefServer.Config.CertsByIP {
			if len(strings.TrimSpace(node.PrivateKey)) < 1 ||
				len(strings.TrimSpace(node.PublicKey)) < 1 {
				errorList.PushBack("ChefServer Public Key and/or Private Key are missing. Set custom_certs_enabled to false to continue without custom certificates.")
			}
		}
	}

	// if CustomCertsEnabled is disabled, then skip validation for custom certs and use self signed certs
	if e.config.Postgresql.Config.CustomCertsEnabled {
		// check if all the default certs are given
		if len(strings.TrimSpace(e.config.Postgresql.Config.RootCA)) > 0 ||
			len(strings.TrimSpace(e.config.Postgresql.Config.PrivateKey)) > 0 ||
			len(strings.TrimSpace(e.config.Postgresql.Config.PublicKey)) > 0 {
			if len(strings.TrimSpace(e.config.Postgresql.Config.RootCA)) < 1 ||
				len(strings.TrimSpace(e.config.Postgresql.Config.PrivateKey)) < 1 ||
				len(strings.TrimSpace(e.config.Postgresql.Config.PublicKey)) < 1 {
				errorList.PushBack("Postgresql RootCA and/or Public Key and/or Private Key are missing. Set custom_certs_enabled to false to continue without custom certificates.")
			}
		}
		// check if all the certs are valid for given IPs
		for _, node := range e.config.Postgresql.Config.CertsByIP {
			if len(strings.TrimSpace(node.RootCA)) < 1 ||
				len(strings.TrimSpace(node.PrivateKey)) < 1 ||
				len(strings.TrimSpace(node.PublicKey)) < 1 {
				errorList.PushBack("Postgresql RootCA and/or Public Key and/or Private Key are missing. Set custom_certs_enabled to false to continue without custom certificates.")
			}
		}
	}

	// if CustomCertsEnabled is disabled, then skip validation for custom certs and use self signed certs
	if e.config.Opensearch.Config.CustomCertsEnabled {
		// check if all the default certs are given
		if len(strings.TrimSpace(e.config.Opensearch.Config.RootCA)) > 0 ||
			len(strings.TrimSpace(e.config.Opensearch.Config.AdminKey)) > 0 ||
			len(strings.TrimSpace(e.config.Opensearch.Config.AdminCert)) > 0 ||
			len(strings.TrimSpace(e.config.Opensearch.Config.PrivateKey)) > 0 ||
			len(strings.TrimSpace(e.config.Opensearch.Config.PublicKey)) > 0 {
			if len(strings.TrimSpace(e.config.Opensearch.Config.RootCA)) < 1 ||
				len(strings.TrimSpace(e.config.Opensearch.Config.AdminKey)) < 1 ||
				len(strings.TrimSpace(e.config.Opensearch.Config.AdminCert)) < 1 ||
				len(strings.TrimSpace(e.config.Opensearch.Config.PrivateKey)) < 1 ||
				len(strings.TrimSpace(e.config.Opensearch.Config.PublicKey)) < 1 {
				errorList.PushBack("Opensearch RootCA and/or Admin Key and/or Admin Cert and/or Public Key and/or Private Key are missing. Set custom_certs_enabled to false to continue without custom certificates.")
			}
		}
		// check if all the certs are valid for given IPs
		for _, node := range e.config.Opensearch.Config.CertsByIP {
			if len(strings.TrimSpace(node.RootCA)) < 1 ||
				len(strings.TrimSpace(node.AdminKey)) < 1 ||
				len(strings.TrimSpace(node.AdminCert)) < 1 ||
				len(strings.TrimSpace(node.PrivateKey)) < 1 ||
				len(strings.TrimSpace(node.PublicKey)) < 1 {
				errorList.PushBack("Opensearch RootCA and/or Admin Key and/or Admin Cert and/or Public Key and/or Private Key are missing. Set custom_certs_enabled to false to continue without custom certificates.")
			}
		}
	}
	return errorList
}

func (e *existingInfra) validateIPs() *list.List {
	const notValidErrorString = "is not valid"
	errorList := list.New()

	for _, element := range e.config.ExistingInfra.Config.AutomatePrivateIps {
		if checkIPAddress(element) != nil {
			errorList.PushBack("Automate private Ip " + element + notValidErrorString)
		}
	}

	for _, element := range e.config.ExistingInfra.Config.ChefServerPrivateIps {
		if checkIPAddress(element) != nil {
			errorList.PushBack("chef server private Ip " + element + notValidErrorString)
		}
	}

	for _, element := range e.config.ExistingInfra.Config.OpensearchPrivateIps {
		if checkIPAddress(element) != nil {
			errorList.PushBack("open search private Ip " + element + notValidErrorString)
		}
	}

	for _, element := range e.config.ExistingInfra.Config.PostgresqlPrivateIps {
		if checkIPAddress(element) != nil {
			errorList.PushBack("Postgresql private Ip " + element + notValidErrorString)
		}
	}

	//
	for _, element := range e.config.Automate.Config.CertsByIP {
		if checkIPAddress(element.IP) != nil {
			errorList.PushBack("Automate IP " + element.IP + " for certs " + notValidErrorString)
		}
	}

	for _, element := range e.config.ChefServer.Config.CertsByIP {
		if checkIPAddress(element.IP) != nil {
			errorList.PushBack("ChefServer IP " + element.IP + " for certs " + notValidErrorString)
		}
	}

	for _, element := range e.config.Opensearch.Config.CertsByIP {
		if checkIPAddress(element.IP) != nil {
			errorList.PushBack("Opensearch IP " + element.IP + " for certs " + notValidErrorString)
		}
	}

	for _, element := range e.config.Postgresql.Config.CertsByIP {
		if checkIPAddress(element.IP) != nil {
			errorList.PushBack("Postgresql IP " + element.IP + " for certs " + notValidErrorString)
		}
	}

	return errorList
}
