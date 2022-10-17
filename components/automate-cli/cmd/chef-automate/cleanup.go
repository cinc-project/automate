package main

import (
	"bytes"
	"io/ioutil"
	"net"
	"strings"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
)

var cleanupFlags = struct {
	onprem        bool
}{}

// const (
// 	FSCRIPT = `systemctl stop chef-automate;
// rm -rf /var/automate-ha;
// rm -rf /hab`
// )

func init() {
	RootCmd.AddCommand(cleanupCmd)
	cleanupCmd.PersistentFlags().BoolVar(&cleanupFlags.onprem, "onprem", false, "Cleaning up all the instances related to onprem ")

}

var cleanupCmd = &cobra.Command {
		Use:  "cleanup",
	    Short: "cleanup the Automate HA instances",
	    Long:  "cleaning up the instance of all the Automate HA related Applications.",
		Annotations: map[string]string{
			NoCheckVersionAnnotation: NoCheckVersionAnnotation,
		},
		RunE: runCleanupCmd,
	}

func runCleanupCmd(cmd *cobra.Command, args []string) error {
	infra, err := getAutomateHAInfraDetails()
	if err != nil {
		return err
	}
	sshUser := infra.Outputs.SSHUser.Value
	sskKeyFile := infra.Outputs.SSHKeyFile.Value
	sshPort := infra.Outputs.SSHPort.Value
	writer.Printf(strings.Join(args, ""))
	if isA2HARBFileExist() {
		if cleanupFlags.onprem {
			frontendIps := append(infra.Outputs.ChefServerPrivateIps.Value, infra.Outputs.AutomatePrivateIps.Value...)
			backendIps := append(infra.Outputs.PostgresqlPrivateIps.Value, infra.Outputs.OpensearchPrivateIps.Value...)
			args = append(args, "--onprem")
			// writer.Printf("IPs: " + strings.Join(frontendIps, "") + args[0])
			for i := 0; i < len(frontendIps); i++ {
				executeCleanupOnRemote(sshUser, sshPort, sskKeyFile, frontendIps[i], args[0])
			}
			for i := 0; i < len(backendIps); i++ {
				executeCleanupOnBackendRemote(sshUser, sshPort, sskKeyFile, backendIps[i], args[0])
			}

		}
	}
	return nil
}

func executeCleanupOnRemote(sshUser string, sshPort string, sshKeyFile string, ip string, remoteType string) {

	pemBytes, err := ioutil.ReadFile(sshKeyFile)
	if err != nil {
		writer.Errorf("Unable to read private key: %v", err)
	}
	signer, err := ssh.ParsePrivateKey(pemBytes)
	if err != nil {
		writer.Errorf("Parsing key failed: %v", err)
	}
	config := &ssh.ClientConfig{
		User:            sshUser,
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	conn, err := ssh.Dial("tcp", net.JoinHostPort(ip, sshPort), config)
	if err != nil {
		writer.Errorf("dial failed:%v", err)
	}
	defer conn.Close()
	session, err := conn.NewSession()
	if err != nil {
		writer.Errorf("session failed:%v", err)
	}
	var stdoutBuf bytes.Buffer
	session.Stdout = &stdoutBuf
	writer.Printf("cleaning up the nodes on IP: " + ip)
	err = session.Run("sudo systemctl stop chef-automate;sudo rm -rf /hab")
	if err != nil {
		writer.Errorf("Run failed:%v", err)
	} else {
		writer.Success("Destroy successful...\n")
	}
	defer session.Close()
	writer.Printf(">%s", stdoutBuf)
	// script :=fmt.Sprintf(FSCRIPT)
	// command := exec.Command(script)
	// err = command.Run()
	// if err != nil {
	// 	writer.Fail(err.Error())
	// 	//return err
	// }
}
// err = session.Run("systemctl stop chef-automate;
// rm -rf /var/automate-ha;
// rm -rf /hab")

func executeCleanupOnBackendRemote(sshUser string, sshPort string, sshKeyFile string, ip string, remoteType string) {

	pemBytes, err := ioutil.ReadFile(sshKeyFile)
	if err != nil {
		writer.Errorf("Unable to read private key: %v", err)
	}
	signer, err := ssh.ParsePrivateKey(pemBytes)
	if err != nil {
		writer.Errorf("Parsing key failed: %v", err)
	}
	config := &ssh.ClientConfig{
		User:            sshUser,
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	conn, err := ssh.Dial("tcp", net.JoinHostPort(ip, sshPort), config)
	if err != nil {
		writer.Errorf("dial failed:%v", err)
	}
	defer conn.Close()
	session, err := conn.NewSession()
	if err != nil {
		writer.Errorf("session failed:%v", err)
	}
	var stdoutBuf bytes.Buffer
	session.Stdout = &stdoutBuf
	writer.Printf("cleaning up the nodes on IP: " + ip)
	err = session.Run("sudo systemctl stop hab-sup;sudo rm -rf /hab")
	if err != nil {
		writer.Errorf("Run failed:%v", err)
	} else {
		writer.Success("Destroy successful...\n")
	}
	defer session.Close()
	writer.Printf(">%s", stdoutBuf)
}