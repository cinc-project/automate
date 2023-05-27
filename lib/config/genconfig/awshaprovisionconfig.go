package genconfig

import (
	"github.com/chef/automate/lib/toml"
)

type AwsHaProvisionConfig struct {
	Fqdn string `json:"fqdn,omitempty" toml:"fqdn,omitempty" mapstructure:"fqdn,omitempty"`
}

func AwsHaProvisionConfigFactory() *AwsHaProvisionConfig {
	return &AwsHaProvisionConfig{}
}

func (c *AwsHaProvisionConfig) Toml() (tomlBytes []byte, err error) {
	return toml.Marshal(c)
}

func (c *AwsHaProvisionConfig) Prompts() (err error) {
	return
}
