package main

import (
	"container/list"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/chef/automate/components/automate-cli/pkg/status"
	"github.com/chef/automate/components/automate-deployment/pkg/cli"
	"github.com/chef/automate/lib/io/fileutils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type DeleteNodeAWSImpl struct {
	config                  AwsConfigToml
	copyConfigForUserPrompt ConfigIp
	automateIpList          []string
	chefServerIpList        []string
	opensearchIpList        []string
	postgresqlIp            []string
	nodeUtils               NodeOpUtils
	flags                   AddDeleteNodeHACmdFlags
	configpath              string
	terraformPath           string
	writer                  *cli.Writer
	fileutils               fileutils.FileUtils
	sshUtil                 SSHUtil
	ConfigIp
}

type ConfigIp struct {
	configAutomateIpList   []string
	configChefServerIpList []string
	configOpensearchIpList []string
	configPostgresqlIpList []string
}

const (
	terrformStateFile = "-state=/hab/a2_deploy_workspace/terraform/terraform.tfstate"
)

func NewDeleteNodeAWS(writer *cli.Writer, flags AddDeleteNodeHACmdFlags, nodeUtils NodeOpUtils, haDirPath string, fileutils fileutils.FileUtils, sshUtil SSHUtil) (HAModifyAndDeploy, error) {
	outputDetails, err := getAutomateHAInfraDetails()
	if err != nil {
		return nil, err
	}

	var ConfigIp = ConfigIp{
		configAutomateIpList:   outputDetails.Outputs.AutomatePrivateIps.Value,
		configChefServerIpList: outputDetails.Outputs.ChefServerPrivateIps.Value,
		configOpensearchIpList: outputDetails.Outputs.OpensearchPrivateIps.Value,
		configPostgresqlIpList: outputDetails.Outputs.PostgresqlPrivateIps.Value,
	}
	return &DeleteNodeAWSImpl{
		config:                  AwsConfigToml{},
		copyConfigForUserPrompt: AwsConfigToml{},
		automateIpList:          []string{},
		chefServerIpList:        []string{},
		opensearchIpList:        []string{},
		postgresqlIp:            []string{},
		nodeUtils:               nodeUtils,
		flags:                   flags,
		configpath:              filepath.Join(haDirPath, "config.toml"),
		terraformPath:           filepath.Join(haDirPath, "terraform"),
		writer:                  writer,
		fileutils:               fileutils,
		sshUtil:                 sshUtil,
		ConfigIp:                ConfigIp,
	}, nil
}

func (dna *DeleteNodeAWSImpl) Execute(c *cobra.Command, args []string) error {
	if !dna.nodeUtils.isA2HARBFileExist() {
		return errors.New(AUTOMATE_HA_INVALID_BASTION)
	}
	if dna.flags.automateIp == "" &&
		dna.flags.chefServerIp == "" &&
		dna.flags.opensearchIp == "" &&
		dna.flags.postgresqlIp == "" {
		c.Help()
		return status.New(status.InvalidCommandArgsError, "Please provide service name and ip address of the node which you want to delete")
	}
	err := dna.validate()
	if err != nil {
		return err
	}
	err = dna.modifyConfig()
	if err != nil {
		return err
	}
	if !dna.flags.autoAccept {
		res, err := dna.promptUserConfirmation()
		if err != nil {
			return err
		}
		if !res {
			return nil
		}
	}
	dna.prepare()
	// run script
	// remove node from destory/aws/*.tfstate

	return dna.runDeploy()
	// return nil
}

func (dna *DeleteNodeAWSImpl) modifyConfig() error {
	err := modifyConfigForDeleteNode(
		&dna.config.Automate.Config.InstanceCount,
		&dna.configAutomateIpList,
		dna.automateIpList,
		&dna.config.Automate.Config.CertsByIP,
	)
	if err != nil {
		return status.Wrap(err, status.ConfigError, "Error modifying automate instance count")
	}
	err = modifyConfigForDeleteNode(
		&dna.config.ChefServer.Config.InstanceCount,
		&dna.configChefServerIpList,
		dna.chefServerIpList,
		&dna.config.ChefServer.Config.CertsByIP,
	)
	if err != nil {
		return status.Wrap(err, status.ConfigError, "Error modifying chef-server instance count")
	}
	err = modifyConfigForDeleteNode(
		&dna.config.Opensearch.Config.InstanceCount,
		&dna.configOpensearchIpList,
		dna.opensearchIpList,
		&dna.config.Opensearch.Config.CertsByIP,
	)
	if err != nil {
		return status.Wrap(err, status.ConfigError, "Error modifying opensearch instance count")
	}
	err = modifyConfigForDeleteNode(
		&dna.config.Postgresql.Config.InstanceCount,
		&dna.configPostgresqlIpList,
		dna.postgresqlIp,
		&dna.config.Postgresql.Config.CertsByIP,
	)
	if err != nil {
		return status.Wrap(err, status.ConfigError, "Error modifying postgresql instance count")
	}
	return nil
}

func (dna *DeleteNodeAWSImpl) prepare() error {
	return dna.nodeUtils.taintTerraform(dna.terraformPath)
}

func (dna *DeleteNodeAWSImpl) promptUserConfirmation() (bool, error) {
	dna.writer.Println("Existing nodes:")
	dna.writer.Println("================================================")
	dna.writer.Println("Automate => " + strings.Join(dna.copyConfigForUserPrompt.configAutomateIpList, ", "))
	dna.writer.Println("Chef-Server => " + strings.Join(dna.copyConfigForUserPrompt.configChefServerIpList, ", "))
	dna.writer.Println("OpenSearch => " + strings.Join(dna.copyConfigForUserPrompt.configOpensearchIpList, ", "))
	dna.writer.Println("Postgresql => " + strings.Join(dna.copyConfigForUserPrompt.configPostgresqlIpList, ", "))
	dna.writer.Println("")
	dna.writer.Println("Nodes to be deleted:")
	dna.writer.Println("================================================")
	if len(dna.automateIpList) > 0 {
		dna.writer.Println("Automate => " + strings.Join(dna.automateIpList, ", "))
	}
	if len(dna.chefServerIpList) > 0 {
		dna.writer.Println("Chef-Server => " + strings.Join(dna.chefServerIpList, ", "))
	}
	if len(dna.opensearchIpList) > 0 {
		dna.writer.Println("OpenSearch => " + strings.Join(dna.opensearchIpList, ", "))
	}
	if len(dna.postgresqlIp) > 0 {
		dna.writer.Println("Postgresql => " + strings.Join(dna.postgresqlIp, ", "))
	}
	return dna.writer.Confirm("This will delete the above nodes from your existing setup. It might take a while. Are you sure you want to continue?")
}

func (dna *DeleteNodeAWSImpl) runDeploy() error {
	tomlbytes, err := ptoml.Marshal(dna.config)
	if err != nil {
		return status.Wrap(err, status.ConfigError, "Error converting config to bytes")
	}
	err = dna.fileUtils.WriteToFile(dna.configpath, tomlbytes)
	if err != nil {
		return err
	}
	err = dna.nodeUtils.genConfig(dna.configpath)
	if err != nil {
		return err
	}
	argsdeploy := []string{"-y"}
	return dna.nodeUtils.executeAutomateClusterCtlCommandAsync("deploy", argsdeploy, upgradeHaHelpDoc)
}

func (dna *DeleteNodeAWSImpl) validate() error {
	dna.automateIpList, dna.chefServerIpList, dna.opensearchIpList, dna.postgresqlIp = splitIPCSV(
		dna.flags.automateIp,
		dna.flags.chefServerIp,
		dna.flags.opensearchIp,
		dna.flags.postgresqlIp,
	)
	var exceptionIps []string
	exceptionIps = append(exceptionIps, dna.automateIpList...)
	exceptionIps = append(exceptionIps, dna.chefServerIpList...)
	exceptionIps = append(exceptionIps, dna.opensearchIpList...)
	exceptionIps = append(exceptionIps, dna.postgresqlIp...)

	updatedConfig, err := dna.nodeUtils.pullAndUpdateConfigAws(&dna.sshUtil, exceptionIps)
	if err != nil {
		return err
	}
	dna.config = *updatedConfig
	dna.copyConfigForUserPrompt = dna.ConfigIp
	if dna.nodeUtils.isManagedServicesOn() {
		if len(dna.opensearchIpList) > 0 || len(dna.postgresqlIp) > 0 {
			return status.New(status.ConfigError, fmt.Sprintf(TYPE_ERROR, "remove"))
		}
	}
	errorList := dna.validateCmdArgs()
	if errorList != nil && errorList.Len() > 0 {
		return status.Wrap(getSingleErrorFromList(errorList), status.ConfigError, "IP address validation failed")
	}
	return nil
}

func (dna *DeleteNodeAWSImpl) validateCmdArgs() *list.List {
	errorList := list.New()
	if len(dna.automateIpList) == 1 {
		allowed, finalCount, err := isFinalInstanceCountAllowed(dna.config.Automate.Config.InstanceCount, -len(dna.automateIpList), AUTOMATE_MIN_INSTANCE_COUNT)
		if err != nil {
			errorList.PushBack("Error occurred in calculating automate final instance count")
		}
		if !allowed {
			errorList.PushBack(fmt.Sprintf("Unable to remove node. Automate instance count cannot be less than %d. Final count %d not allowed.", AUTOMATE_MIN_INSTANCE_COUNT, finalCount))
		}
		fmt.Println("automate", allowed)
		errorList.PushBackList(checkIfPresentInPrivateIPList(dna.configAutomateIpList, dna.automateIpList, "Automate"))
	} else if len(dna.automateIpList) != 0 {
		errorList.PushBack(fmt.Sprintf("only One automate ip adress is allowed"))
	}
	if len(dna.chefServerIpList) == 1 {
		allowed, finalCount, err := isFinalInstanceCountAllowed(dna.config.ChefServer.Config.InstanceCount, -len(dna.chefServerIpList), CHEF_SERVER_MIN_INSTANCE_COUNT)
		if err != nil {
			errorList.PushBack("Error occurred in calculating chef server final instance count")
		}
		if !allowed {
			errorList.PushBack(fmt.Sprintf("Unable to remove node. Chef Server instance count cannot be less than %d. Final count %d not allowed.", CHEF_SERVER_MIN_INSTANCE_COUNT, finalCount))
		}
		fmt.Println("chef-server", allowed)
		errorList.PushBackList(checkIfPresentInPrivateIPList(dna.configChefServerIpList, dna.chefServerIpList, "Chef-Server"))
	} else if len(dna.chefServerIpList) != 0 {
		errorList.PushFront(fmt.Sprintf("only One Chef server ip adress is allowed"))
	}
	if len(dna.opensearchIpList) == 1 {
		allowed, finalCount, err := isFinalInstanceCountAllowed(dna.config.Opensearch.Config.InstanceCount, -len(dna.opensearchIpList), OPENSEARCH_MIN_INSTANCE_COUNT)
		if err != nil {
			errorList.PushBack("Error occurred in calculating opensearch final instance count")
		}
		if !allowed {
			errorList.PushBack(fmt.Sprintf("Unable to remove node. OpenSearch instance count cannot be less than %d. Final count %d not allowed.", OPENSEARCH_MIN_INSTANCE_COUNT, finalCount))
		}
		fmt.Println("os", allowed)
		errorList.PushBackList(checkIfPresentInPrivateIPList(dna.configPostgresqlIpList, dna.opensearchIpList, "OpenSearch"))
	} else if len(dna.opensearchIpList) != 0 {
		errorList.PushBack(fmt.Sprintf("only One postgres ip adress is allowed"))
	}
	if len(dna.postgresqlIp) == 1 {
		allowed, finalCount, err := isFinalInstanceCountAllowed(dna.config.Postgresql.Config.InstanceCount, -len(dna.postgresqlIp), POSTGRESQL_MIN_INSTANCE_COUNT)
		if err != nil {
			errorList.PushBack("Error occurred in calculating postgresql final instance count")
		}
		if !allowed {
			errorList.PushBack(fmt.Sprintf("Unable to remove node. Postgresql instance count cannot be less than %d. Final count %d not allowed.", POSTGRESQL_MIN_INSTANCE_COUNT, finalCount))
		}
		fmt.Println("pg", allowed)
		errorList.PushBackList(checkIfPresentInPrivateIPList(dna.configPostgresqlIpList, dna.postgresqlIp, "Postgresql"))
	} else if len(dna.postgresqlIp) != 0 {
		errorList.PushBack(fmt.Sprintf("only One opensearch ip adress is allowed"))
	}
	return errorList
}
