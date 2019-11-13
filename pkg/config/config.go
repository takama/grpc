package config

import (
	"github.com/takama/grpc/client"
	"github.com/takama/grpc/pkg/info"
	"github.com/takama/grpc/pkg/logger"
	"github.com/takama/grpc/pkg/server"
	"github.com/takama/grpc/pkg/system"

	"github.com/spf13/viper"
)

// Default values: host, port, etc
const (
	// ServiceName - default service name
	ServiceName       = "grpc"
	ClientServiceName = "grpc"

	APIVersion = "v1"

	DefaultConfigPath = "config/default.conf"

	DefaultServerPort                 = 8000
	DefaultInfoPort                   = 8080
	DefaultClientBalancer             = "round_robin"
	DefaultClientInsecure             = false
	DefaultClientEnvoyProxy           = false
	DefaultClientWaitForReady         = false
	DefaultClientTimeout              = 15
	DefaultClientKeepAliveTime        = 10
	DefaultClientKeepAliveTimeout     = 5
	DefaultClientKeepAliveForce       = false
	DefaultClientRetryActive          = false
	DefaultClientEnvoyRetryReason     = "5xx"
	DefaultClientEnvoyRetryGRPCReason = "unavailable"
	DefaultClientEnvoyRetryCount      = 30
	DefaultClientEnvoyRetryTimeout    = 5
	DefaultClientBackoffMultiplier    = 1.6
	DefaultClientBackoffJitter        = 0.2
	DefaultClientBackoffDelayMin      = 1
	DefaultClientBackoffDelayMax      = 120
	DefaultServerConnectionIdle       = 0
	DefaultServerConnectionAge        = 0
	DefaultServerConnectionAgeGrace   = 0
	DefaultServertKeepAliveTime       = 300
	DefaultServerKeepAliveTimeout     = 10
	DefaultGracePeriod                = 30
	DefaultInfoStatistics             = true
	DefaultLoggerLevel                = logger.LevelInfo
)

// Config -- Base config structure
type Config struct {
	Client client.Config
	Server server.Config
	Info   info.Config
	Logger logger.Config
	System system.Config
}

// New - returns new config record initialized with default values
// nolint: funlen
func New() (*Config, error) {
	cfg := &Config{
		Client: client.Config{
			Scheme:       ServiceName,
			Host:         ClientServiceName,
			Balancer:     DefaultClientBalancer,
			Insecure:     DefaultClientInsecure,
			EnvoyProxy:   DefaultClientEnvoyProxy,
			WaitForReady: DefaultClientWaitForReady,
			Timeout:      DefaultClientTimeout,
			Keepalive: client.Keepalive{
				Time:    DefaultClientKeepAliveTime,
				Timeout: DefaultClientKeepAliveTimeout,
				Force:   DefaultClientKeepAliveForce,
			},
			Retry: client.Retry{
				Active: DefaultClientRetryActive,
				Envoy: client.Envoy{
					Reason: client.Reason{
						Primary: DefaultClientEnvoyRetryReason,
						GRPC:    DefaultClientEnvoyRetryGRPCReason,
					},
					Count:   DefaultClientEnvoyRetryCount,
					Timeout: DefaultClientEnvoyRetryTimeout,
				},
				Backoff: client.Backoff{
					Multiplier: DefaultClientBackoffMultiplier,
					Jitter:     DefaultClientBackoffJitter,
					Delay: client.Delay{
						Min: DefaultClientBackoffDelayMin,
						Max: DefaultClientBackoffDelayMax,
					},
				},
			},
		},
		Server: server.Config{
			Port: DefaultServerPort,
			Connection: server.Connection{
				Idle:  DefaultServerConnectionIdle,
				Age:   DefaultServerConnectionAge,
				Grace: DefaultServerConnectionAgeGrace,
				Keepalive: server.Keepalive{
					Time:    DefaultServertKeepAliveTime,
					Timeout: DefaultServerKeepAliveTimeout,
				},
			},
		},
		Info: info.Config{
			Port:       DefaultInfoPort,
			Statistics: DefaultInfoStatistics,
		},
		System: system.Config{
			Grace: system.Grace{
				Period: DefaultGracePeriod,
			},
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
