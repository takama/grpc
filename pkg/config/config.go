package config

import (
	"github.com/takama/grpc/pkg/info"
	"github.com/takama/grpc/pkg/logger"
	"github.com/takama/grpc/pkg/server"

	"github.com/spf13/viper"
)

// Default values: host, port, etc
const (
	// ServiceName - default service name
	ServiceName = "grpc"

	APIVersion = "v1"

	DefaultConfigPath = "config/default.conf"

	DefaultServerPort     = 8000
	DefaultInfoPort       = 8080
	DefaultInfoStatistics = true
	DefaultLoggerLevel    = logger.LevelInfo
)

// Config -- Base config structure
type Config struct {
	Server server.Config
	Info   info.Config
	Logger logger.Config
}

// New - returns new config record initialized with default values
func New() (*Config, error) {
	cfg := &Config{
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
