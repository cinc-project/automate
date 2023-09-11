package config

import (
	"github.com/chef/automate/components/user-settings-service/pkg/storage"
	"github.com/chef/automate/lib/logger"
	"github.com/chef/automate/lib/tls/certs"
)

var ServiceName = "user-settings-service"

type UserSettings struct {
	Service          `mapstructure:"service"`
	Postgres         `mapstructure:"postgres"`
	*certs.TLSConfig `mapstructure:"tls"`
	storageClient    storage.Client
}

// Service is a base config options struct for all services
type Service struct {
	Host        string `mapstructure:"host"`
	Port        int    `mapstructure:"port"`
	MetricsPort int    `mapstructure:"metrics_port"`
	LogLevel    string `mapstructure:"log_level"`
	LogFormat   string `mapstructure:"log-format"`
	Logger      logger.Logger
}

type Postgres struct {
	URI          string `mapstructure:"uri"`
	Database     string `mapstructure:"database"`
	SchemaPath   string `mapstructure:"schema_path"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
}

// SetStorage sets the storage client for the service
func (s *UserSettings) SetStorage(c storage.Client) {
	s.storageClient = c
}

// GetStorage returns the storage client for the service
func (s *UserSettings) GetStorage() storage.Client {
	return s.storageClient
}

// SetLogLevel sets the log level and format for the service
func (s *Service) SetLogConfig() error {
	l, err := logger.NewLogger(s.LogFormat, s.LogLevel)
	if err != nil {
		return err
	}
	s.Logger = l
	return nil
}
