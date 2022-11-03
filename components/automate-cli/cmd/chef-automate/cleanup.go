package main

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/spf13/cobra"
)

var (
	AWS        string = "aws"
	DEPLOYMENT string = "deployment"

	AWS_PROVISION = `
	for i in 1;do i=$PWD;cd /hab/a2_deploy_workspace/terraform/destroy/aws/;terraform init;cd $i;done
	%s
	for i in 1;do i=$PWD;cd /hab/a2_deploy_workspace/terraform/destroy/aws/;terraform destroy  -state=/hab/a2_deploy_workspace/terraform/destroy/aws/terraform.tfstate -auto-approve;cd $i;done
`
)

var cleanupFlags = struct {
	onprem bool
	aws    bool
	force  bool
}{}

func init() {
	RootCmd.AddCommand(cleanupCmd)
	cleanupCmd.PersistentFlags().BoolVar(&cleanupFlags.onprem, "onprem-deployment", false, "Cleaning up all the instances related to onprem ")
	cleanupCmd.PersistentFlags().BoolVar(&cleanupFlags.aws, "aws-deployment", false, "Remove AWS resources created by provisioning and clean-up hab workspace")
	cleanupCmd.PersistentFlags().BoolVar(&cleanupFlags.force, "force", false, "Remove backup storage")

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
					writer.Println("Cleanup has started on " + servername + " node : " + automateIps[i] + "\n")
					_, err := ConnectAndExecuteCommandOnRemote(sshUser, sshPort, sskKeyFile, automateIps[i], FRONTENDCLEANUP_COMMANDS)
					if err != nil {
						writer.Errorf("%s", err.Error())
						return err
					} else {
						writer.Success("Cleanup is completed on " + servername + " node : " + automateIps[i] + "\n")
					}
				}
				for i := 0; i < len(chefserverIps); i++ {
					servername := "chef server"
					writer.Println("Cleanup has started on " + servername + " node : " + chefserverIps[i] + "\n")
					_, err := ConnectAndExecuteCommandOnRemote(sshUser, sshPort, sskKeyFile, chefserverIps[i], FRONTENDCLEANUP_COMMANDS)
					if err != nil {
						writer.Error(err.Error())
						return err
					} else {
						writer.Success("Cleanup is completed on " + servername + " node : " + chefserverIps[i] + "\n")
					}
				}
				for i := 0; i < len(postgresqlIps); i++ {
					servername := "postgresql"
					writer.Println("Cleanup has started on " + servername + " node : " + postgresqlIps[i] + "\n")
					_, err := ConnectAndExecuteCommandOnRemote(sshUser, sshPort, sskKeyFile, postgresqlIps[i], BACKENDCLEANUP_COMMANDS)
					if err != nil {
						writer.Error(err.Error())
						return err
					} else {
						writer.Success("Cleanup is completed on " + servername + " node : " + postgresqlIps[i] + "\n")
					}
				}
				for i := 0; i < len(opensearchIps); i++ {
					servername := "opensearch"
					writer.Println("Cleanup has started on " + servername + " node : " + opensearchIps[i] + "\n")
					_, err := ConnectAndExecuteCommandOnRemote(sshUser, sshPort, sskKeyFile, opensearchIps[i], BACKENDCLEANUP_COMMANDS)
					if err != nil {
						writer.Error(err.Error())
						return err
					} else {
						writer.Success("Cleanup is completed on " + servername + " node : " + opensearchIps[i] + "\n")
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

				appendString := ""
				if infra.Outputs.BackupConfigS3.Value == "true" && cleanupFlags.force {
					appendString = appendString + `for i in 1;do i=$PWD;cd /hab/a2_deploy_workspace/terraform/destroy/aws/;cp -r ../../.tf_arch .;cp -r ../../../a2ha.rb ..;terraform apply -var="destroy_bucket=true" -auto-approve;cd $i;done`
				} else if infra.Outputs.BackupConfigEFS.Value == "true" && !cleanupFlags.force {
					appendString = appendString + `for i in 1;do i=$PWD;cd /hab/a2_deploy_workspace/terraform/destroy/aws/;terraform state rm "module.efs[0].aws_efs_file_system.backups";cd $i;done`
				}

				writer.Println("Cleaning up all AWS provisioned resources.")
				if arch == DEPLOYMENT {
					args := []string{
						"-c",
						fmt.Sprintf(AWS_PROVISION, appendString),
					}
					err = executeCommand("/bin/sh", args, "")
					if err != nil {
						return err
					}
				}
				if arch == "aws" {
					args := []string{
						"-c",
						fmt.Sprintf(AWS_PROVISION, appendString),
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
				writer.Success("Cleaning up completed.")
			}
		}
	} else {
		writer.Println("\nCleanup not executed.")
	}
	return nil
}
