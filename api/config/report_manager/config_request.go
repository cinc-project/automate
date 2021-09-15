package report_manager

import (
	"fmt"
	"net"

	ac "github.com/chef/automate/api/config/shared"
	w "github.com/chef/automate/api/config/shared/wrappers"
	"github.com/chef/automate/lib/config"
	wrappers "github.com/golang/protobuf/ptypes/wrappers"
)

const (
	defaultResolverPort = 53
	defaultServerPort   = "5432"
)

// NewConfigRequest returns a new instance of ConfigRequest with zero values.
func NewConfigRequest() *ConfigRequest {
	return &ConfigRequest{
		V1: &ConfigRequest_V1{
			Sys: &ConfigRequest_V1_System{
				Mlsa:    &ac.Mlsa{},
				Service: &ConfigRequest_V1_System_Service{},
				Log:     &ac.Log{},
			},
		},
	}
}

// DefaultConfigRequest returns a new instance of ConfigRequest with default values.
func DefaultConfigRequest() *ConfigRequest {
	c := NewConfigRequest()

	c.V1.Sys.Service.Port = w.Int32(10152)
	return c
}

// Validate validates that the config is sufficient to start the Service. If
// validation succeeds it will return nil, if it fails it will return a new
// instance of config.InvalidConfigError that has the missing keys and invalid
// fields populated.
func (c *ConfigRequest) Validate() error {
	return nil
}

// SetGlobalConfig is a callback that allows us to populate the
// configuration of an individual service with global automate config.
func (c *ConfigRequest) SetGlobalConfig(g *ac.GlobalConfig) {
	c.V1.Sys.Mlsa = g.V1.Mlsa
}

// PrepareSystemConfig returns a system configuration that can be used
// to start the service.
func (c *ConfigRequest) PrepareSystemConfig(creds *ac.TLSCredentials) (ac.PreparedSystemConfig, error) {
	c.V1.Sys.Tls = creds
	return c.V1.Sys, nil
}

func getSystemResolvers() []*wrappers.StringValue {
	ns := config.GetSystemResolvers()
	resolvers := make([]*wrappers.StringValue, 0, len(ns))

	for _, n := range ns {
		_, _, err := net.SplitHostPort(n)
		if err == nil {
			resolvers = append(resolvers, w.String(n))
		} else {
			resolvers = append(resolvers, w.String(fmt.Sprintf("%s:%d", n, defaultResolverPort)))
		}
	}
	return resolvers
}

func isIPAddress(addr string) bool {
	return net.ParseIP(addr) != nil
}
