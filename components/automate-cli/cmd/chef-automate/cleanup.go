package main

import (
	"io/ioutil"
	"strings"

	"github.com/spf13/cobra"
)

var (
	AWS        string = "aws"
	DEPLOYMENT string = "deployment"

	AWS_PROVISION_INCOMPLETE = `
	for i in 1;do i=$PWD;cd /hab/a2_deploy_workspace/terraform/;terraform destroy -state=/hab/a2_deploy_workspace/terraform/destroy/aws/terraform.tfstate -auto-approve;cd $i;done
	`

	AWS_PROVISION_COMPLETE = `
	for i in 1;do i=$PWD;cd /hab/a2_deploy_workspace/terraform/destroy/aws/;terraform init;cd $i;done
	for i in 1;do i=$PWD;cd /hab/a2_deploy_workspace/terraform/destroy/aws/;cp -r ../../.tf_arch .;cp -r ../../../a2ha.rb ..;terraform apply -var="destroy_bucket=true" -auto-approve;cd $i;done
	for i in 1;do i=$PWD;cd /hab/a2_deploy_workspace/terraform/destroy/aws/;terraform destroy -auto-approve;cd $i;done
`
)

var cleanupFlags = struct {
	onprem bool
	aws    bool
}{}

func init() {
	RootCmd.AddCommand(cleanupCmd)
	cleanupCmd.PersistentFlags().BoolVar(&cleanupFlags.onprem, "onprem", false, "Cleaning up all the instances related to onprem ")
	cleanupCmd.PersistentFlags().BoolVar(&cleanupFlags.aws, "aws", false, "Remove AWS resources created by provisioning and clean-up hab workspace")

}

var cleanupCmd = &cobra.Command{
	Use:   "cleanup",
	Short: "cleanup the Automate HA instances",
	Long:  "cleaning up the instance of all the Automate HA related Applications.",
	Annotations: map[string]string{
		NoCheckVersionAnnotation: NoCheckVersionAnnotation,
	},
	RunE:   runCleanupCmd,
	Hidden: true,
}

const (
	FRONTENDCLEANUP_COMMANDS = `
		sudo systemctl stop chef-automate;
		sudo rm -rf /hab;
		sudo rm -rf /var/automate-ha;
		`

	BACKENDCLEANUP_COMMANDS = `
		sudo systemctl stop hab-sup;
		sudo rm -rf /hab; 
		sudo rm -rf /var/automate-ha;
		`
)

func runCleanupCmd(cmd *cobra.Command, args []string) error {
	infra, err := getAutomateHAInfraDetails()
	if err != nil {
		return err
	}
	if infra != nil {
		writer.Printf(strings.Join(args, ""))
		if isA2HARBFileExist() {
			sshUser := infra.Outputs.SSHUser.Value
			sskKeyFile := infra.Outputs.SSHKeyFile.Value
			sshPort := infra.Outputs.SSHPort.Value
			if cleanupFlags.onprem {
				automateIps := infra.Outputs.AutomatePrivateIps.Value
				chefserverIps := infra.Outputs.ChefServerPrivateIps.Value
				postgresqlIps := infra.Outputs.PostgresqlPrivateIps.Value
				opensearchIps := infra.Outputs.OpensearchPrivateIps.Value
				for i := 0; i < len(automateIps); i++ {
					servername := "Automate"
					writer.Println("cleanup is starting on " + servername + " node : " + automateIps[i] + "\n")
					_, err := ConnectAndExecuteCommandOnRemote(sshUser, sshPort, sskKeyFile, automateIps[i], FRONTENDCLEANUP_COMMANDS)
					if err != nil {
						writer.Errorf("%s", err.Error())
						return err
					} else {
						writer.Success("cleanup is successfully completed on " + servername + " node : " + automateIps[i] + "\n")
					}
				}
				for i := 0; i < len(chefserverIps); i++ {
					servername := "chef server"
					writer.Println("cleanup is starting on " + servername + " node : " + chefserverIps[i] + "\n")
					_, err := ConnectAndExecuteCommandOnRemote(sshUser, sshPort, sskKeyFile, chefserverIps[i], FRONTENDCLEANUP_COMMANDS)
					if err != nil {
						writer.Error(err.Error())
						return err
					} else {
						writer.Success("cleanup is successful on " + servername + " node : " + chefserverIps[i] + "\n")
					}
				}
				for i := 0; i < len(postgresqlIps); i++ {
					servername := "postgresql"
					writer.Println("cleanup is starting on " + servername + " node : " + postgresqlIps[i] + "\n")
					_, err := ConnectAndExecuteCommandOnRemote(sshUser, sshPort, sskKeyFile, postgresqlIps[i], BACKENDCLEANUP_COMMANDS)
					if err != nil {
						writer.Error(err.Error())
						return err
					} else {
						writer.Success("cleanup is successful on " + servername + " node : " + postgresqlIps[i] + "\n")
					}
				}
				for i := 0; i < len(opensearchIps); i++ {
					servername := "opensearch"
					writer.Println("cleanup is starting on " + servername + " node : " + opensearchIps[i] + "\n")
					_, err := ConnectAndExecuteCommandOnRemote(sshUser, sshPort, sskKeyFile, opensearchIps[i], BACKENDCLEANUP_COMMANDS)
					if err != nil {
						writer.Error(err.Error())
						return err
					} else {
						writer.Success("cleanup is successful on " + servername + " node : " + opensearchIps[i] + "\n")
					}
				}
				cleanUpScript := "hab pkg uninstall chef/automate-ha-deployment"
				args := []string{
					"-c",
					cleanUpScript,
				}
				err := executeCommand("/bin/sh", args, "")
				if err != nil {
					return err
				}
			}
			if cleanupFlags.aws {
				archBytes, err := ioutil.ReadFile("/hab/a2_deploy_workspace/terraform/.tf_arch") // nosemgrep
				if err != nil {
					writer.Errorf("%s", err.Error())
					return err
				}
				var arch = strings.Trim(string(archBytes), "\n")

				if arch == DEPLOYMENT {
					writer.Println("DEPLOYMENT: " + DEPLOYMENT)
					args := []string{
						"-c",
						AWS_PROVISION_INCOMPLETE,
					}
					err = executeCommand("/bin/sh", args, "")
					if err != nil {
						return err
					}
				}
				if arch == "aws" {
					writer.Println("AWS: " + AWS)
					args := []string{
						"-c",
						AWS_PROVISION_COMPLETE,
					}
					err = executeCommand("/bin/sh", args, "")
					if err != nil {
						return err
					}
				}
				cleanUpScript := "hab pkg uninstall chef/automate-ha-deployment"
				args = []string{
					"-c",
					cleanUpScript,
				}
				err = executeCommand("/bin/sh", args, "")
				if err != nil {
					return err
				}
			}
		}
	} else {
		writer.Println("\nCleanup not executed.")
	}
	return nil
}
