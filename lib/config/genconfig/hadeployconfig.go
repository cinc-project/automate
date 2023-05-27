package genconfig

import (
	"github.com/chef/automate/lib/pmt"
	"github.com/chef/automate/lib/toml"
)

type HaDeployConfig struct {
	Prompt        pmt.Prompt             `toml:"-"`
	Architecture  *Architecture          `toml:"architecture,omitempty"`
	ObjectStorage *ObjectStorage         `toml:"object_storage,omitempty"`
	Automate      *AutomateSettings      `toml:"automate,omitempty"`
	ChefServer    *ServerConfigSettings  `toml:"chef_server,omitempty"`
	Opensearch    *ServerConfigSettings  `toml:"opensearch,omitempty"`
	Postgresql    *ServerConfigSettings  `toml:"postgresql,omitempty"`
	ExistingInfra *ExistingInfraSettings `toml:"existing_infra,omitempty"`
	External      *ExternalSettings      `toml:"external,omitempty"`
}

type Architecture struct {
	ExistingInfra *ExistingInfraArch `toml:"existing_infra,omitempty"`
}

type ExistingInfraArch struct {
	SSHUser          string `toml:"ssh_user,omitempty"`
	SSHGroupName     string `toml:"ssh_group_name,omitempty"`
	SSHKeyFile       string `toml:"ssh_key_file,omitempty"`
	SSHPort          string `toml:"ssh_port,omitempty"`
	SecretsKeyFile   string `toml:"secrets_key_file,omitempty"`
	SecretsStoreFile string `toml:"secrets_store_file,omitempty"`
	Architecture     string `toml:"architecture,omitempty"`
	WorkspacePath    string `toml:"workspace_path,omitempty"`
	BackupMount      string `toml:"backup_mount,omitempty"`
	BackupConfig     string `toml:"backup_config,omitempty"`
}

type ObjectStorage struct {
	Config *ConfigObjectStorage `toml:"config,omitempty"`
}

type ConfigObjectStorage struct {
	BucketName string `toml:"bucket_name,omitempty"`
	AccessKey  string `toml:"access_key,omitempty"`
	SecretKey  string `toml:"secret_key,omitempty"`
	Endpoint   string `toml:"endpoint,omitempty"`
	Region     string `toml:"region,omitempty"`
}

type AutomateSettings struct {
	Config *ConfigAutomateSettings `toml:"config,omitempty"`
}

type ConfigAutomateSettings struct {
	AdminPassword string `toml:"admin_password,omitempty"`
	Fqdn          string `toml:"fqdn,omitempty"`
	ConfigFile    string `toml:"config_file,omitempty"`
	ConfigSettings
}

type ServerConfigSettings struct {
	Config *ConfigSettings `toml:"config,omitempty"`
}

type ConfigSettings struct {
	InstanceCount     string `toml:"instance_count,omitempty"`
	EnableCustomCerts bool   `toml:"enable_custom_certs,omitempty"`
}

type ExistingInfraSettings struct {
	Config *ConfigExistingInfraSettings `toml:"config,omitempty"`
}

type ConfigExistingInfraSettings struct {
	AutomatePrivateIps   []string `toml:"automate_private_ips,omitempty"`
	ChefServerPrivateIps []string `toml:"chef_server_private_ips,omitempty"`
	OpensearchPrivateIps []string `toml:"opensearch_private_ips,omitempty"`
	PostgresqlPrivateIps []string `toml:"postgresql_private_ips,omitempty"`
}

type ExternalSettings struct {
	Database *ExternalDBSettings `toml:"database,omitempty"`
}

type ExternalDBSettings struct {
	Type       string              `toml:"type,omitempty"`
	PostgreSQL *ExternalPgSettings `toml:"postgre_sql,omitempty"`
	OpenSearch *ExternalOsSettings `toml:"open_search,omitempty"`
}

type ExternalPgSettings struct {
	InstanceURL        string `toml:"instance_url,omitempty"`
	SuperuserUsername  string `toml:"superuser_username,omitempty"`
	SuperuserPassword  string `toml:"superuser_password,omitempty"`
	DbuserUsername     string `toml:"dbuser_username,omitempty"`
	DbuserPassword     string `toml:"dbuser_password,omitempty"`
	PostgresqlRootCert string `toml:"postgresql_root_cert,omitempty"`
}

type ExternalOsSettings struct {
	OpensearchDomainName   string                 `toml:"opensearch_domain_name,omitempty"`
	OpensearchDomainURL    string                 `toml:"opensearch_domain_url,omitempty"`
	OpensearchUsername     string                 `toml:"opensearch_username,omitempty"`
	OpensearchUserPassword string                 `toml:"opensearch_user_password,omitempty"`
	OpensearchRootCert     string                 `toml:"opensearch_root_cert,omitempty"`
	Aws                    *AwsExternalOsSettings `toml:"aws,omitempty"`
}

type AwsExternalOsSettings struct {
	AwsOsSnapshotRoleArn          string `toml:"aws_os_snapshot_role_arn,omitempty"`
	OsSnapshotUserAccessKeyID     string `toml:"os_snapshot_user_access_key_id,omitempty"`
	OsSnapshotUserAccessKeySecret string `toml:"os_snapshot_user_access_key_secret,omitempty"`
}

func HaDeployConfigFactory(p pmt.Prompt) *HaDeployConfig {
	return &HaDeployConfig{
		Prompt: p,
	}
}

func (c *HaDeployConfig) Toml() (tomlBytes []byte, err error) {
	return toml.Marshal(c)
}

func (c *HaDeployConfig) Prompts() (err error) {
	err = c.PromptSsh()
	if err != nil {
		return
	}
	return
}

func (c *HaDeployConfig) PromptSsh() (err error) {
	sshUser, err := c.Prompt.InputWord("SSH User Name")
	if err != nil {
		return
	}
	if c.Architecture == nil {
		c.Architecture = &Architecture{}
	}
	if c.Architecture.ExistingInfra == nil {
		c.Architecture.ExistingInfra = &ExistingInfraArch{}
	}
	c.Architecture.ExistingInfra.SSHUser = sshUser

	sshGroup, err := c.Prompt.InputWordDefault("SSH Group", sshUser)
	if err != nil {
		return
	}
	c.Architecture.ExistingInfra.SSHGroupName = sshGroup

	sshPort, err := c.Prompt.InputWordDefault("SSH Port", "22")
	if err != nil {
		return
	}
	c.Architecture.ExistingInfra.SSHPort = sshPort

	sshKeyFile, err := c.Prompt.InputStringRegxDefault("SSH Key File [eg. \"~/.ssh/A2HA.pem\"]", "^[~]{0,1}(\\/[\\w^\\. ]+)+\\/?$", "")
	if err != nil {
		return
	}
	c.Architecture.ExistingInfra.SSHKeyFile = sshKeyFile

	c.Architecture.ExistingInfra.SecretsKeyFile = "/hab/a2_deploy_workspace/secrets.key"
	c.Architecture.ExistingInfra.SecretsStoreFile = "/hab/a2_deploy_workspace/secrets.json"
	c.Architecture.ExistingInfra.Architecture = "existing_nodes"
	c.Architecture.ExistingInfra.WorkspacePath = "workspace_path"

	err = c.PromptBackup()
	if err != nil {
		return
	}
	return
}

func (c *HaDeployConfig) PromptBackup() (err error) {
	isBackupNeeded, err := c.Prompt.Confirm("Backup need to be configured during deployment", "yes", "no")
	if isBackupNeeded {
		if c.Architecture == nil {
			c.Architecture = &Architecture{}
		}
		if c.Architecture.ExistingInfra == nil {
			c.Architecture.ExistingInfra = &ExistingInfraArch{}
		}

		c.Architecture.ExistingInfra.BackupMount = "/mnt/automate_backups"

		backupOption, err1 := c.Prompt.Select("Which backup option will you use", "AWS S3", "Minio", "Object Store", "File System", "NFS", "EFS")
		if err1 != nil {
			return err1
		}
		backupConfig := ""
		switch backupOption {
		case "AWS S3", "Minio", "Object Store":
			backupConfig = "object_storage"
		case "File System", "NFS", "EFS":
			backupConfig = "file_system"
		}
		c.Architecture.ExistingInfra.BackupConfig = backupConfig

		if backupConfig == "object_storage" {
			err1 := c.PromptObjectStorageSettings(backupOption)
			if err1 != nil {
				return err1
			}
		} else if backupConfig == "file_system" {
			backupMountLoc, err1 := c.Prompt.InputStringRegxDefault("Backup Mount Location", "^[~]{0,1}(\\/[\\w^\\. ]+)+\\/?$", "/mnt/automate_backups")
			if err1 != nil {
				return err1
			}
			c.Architecture.ExistingInfra.BackupMount = backupMountLoc
		}
	}
	return
}

func (c *HaDeployConfig) PromptObjectStorageSettings(backupOption string) (err error) {

	bucketName, err := c.Prompt.InputStringRegx("Backup bucket name", "^[a-zA-Z0-9_-]+$")
	if err != nil {
		return
	}
	if c.Architecture == nil {
		c.Architecture = &Architecture{}
	}
	if c.Architecture.ExistingInfra == nil {
		c.Architecture.ExistingInfra = &ExistingInfraArch{}
	}
	if c.ObjectStorage == nil {
		c.ObjectStorage = &ObjectStorage{}
	}
	if c.ObjectStorage.Config == nil {
		c.ObjectStorage.Config = &ConfigObjectStorage{}
	}
	c.ObjectStorage.Config.BucketName = bucketName

	if backupOption == "AWS S3" {
		accessKey, err1 := c.Prompt.InputStringRegx("Access key for bucket", "^[A-Z0-9]{20}$")
		if err1 != nil {
			return err1
		}
		c.ObjectStorage.Config.AccessKey = accessKey

		secretKey, err1 := c.Prompt.InputStringRegx("SecretKey for bucket", "^[A-Za-z0-9/+=]{40}$")
		if err1 != nil {
			return err1
		}
		c.ObjectStorage.Config.SecretKey = secretKey
		c.ObjectStorage.Config.Endpoint = "https://s3.amazonaws.com"

		bucketRegion, err1 := GetAwsRegion(c.Prompt)
		if err1 != nil {
			return err1
		}
		c.ObjectStorage.Config.Endpoint = bucketRegion
	} else {
		accessKey, err1 := c.Prompt.InputStringRegx("Access key for bucket", "^[A-Za-z0-9/+=]+$")
		if err1 != nil {
			return err1
		}
		c.ObjectStorage.Config.AccessKey = accessKey

		secretKey, err1 := c.Prompt.InputStringRegx("SecretKey for bucket", "^[A-Za-z0-9/+=]+$")
		if err1 != nil {
			return err1
		}
		c.ObjectStorage.Config.SecretKey = secretKey
		bucketEndpoint, err1 := c.Prompt.InputStringRegx("Endpoint for bucket", "^((http|https)://)[-a-zA-Z0-9@:%._\\+~#?&//=]{2,256}\\.[a-z]{2,6}\\b([-a-zA-Z0-9@:%._\\+~#?&//=]*)$")
		if err1 != nil {
			return err1
		}
		c.ObjectStorage.Config.Endpoint = bucketEndpoint
	}
	return
}
