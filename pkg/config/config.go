package config

import (
	"github.com/takama/grpc/pkg/client"
	"github.com/takama/grpc/pkg/info"
	"github.com/takama/grpc/pkg/logger"
	"github.com/takama/grpc/pkg/server"

	"github.com/spf13/viper"
)

// Default values: host, port, etc
const (
	// ServiceName - default service name
	ServiceName       = "grpc"
	ClientServiceName = "grpc"

	APIVersion = "v1"

	DefaultConfigPath = "config/default.conf"

	DefaultClientPort     = 8000
	DefaultServerPort     = 8000
	DefaultInfoPort       = 8080
	DefaultClientInsecure = false
	DefaultInfoStatistics = true
	DefaultLoggerLevel    = logger.LevelInfo
)

// Config -- Base config structure
type Config struct {
	Client client.Config
	Server server.Config
	Info   info.Config
	Logger logger.Config
}

// New - returns new config record initialized with default values
func New() (*Config, error) {
	cfg := &Config{
		Client: client.Config{
			Host:     ClientServiceName,
			Port:     DefaultClientPort,
			Insecure: DefaultClientInsecure,
		},
		Server: server.Config{
			Port: DefaultServerPort,
		},
		Info: info.Config{
			Port:       DefaultInfoPort,
			Statistics: DefaultInfoStatistics,
		},
		Logger: logger.Config{
			Level: DefaultLoggerLevel,
		},
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
