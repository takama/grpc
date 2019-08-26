package boot

import (
	"crypto/tls"

	"github.com/takama/grpc/pkg/config"
	"github.com/takama/grpc/pkg/helper"
	"github.com/takama/grpc/pkg/logger"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Setup the configuration, logger
func Setup() (*config.Config, *zap.Logger) {
	cfg, err := config.New()
	helper.LogF("Load config error", err)
	// Setup zap logger
	log := logger.New(&cfg.Logger)

	return cfg, log
}

// TLSOption gives TLS secure/insecure option
func TLSOption(host string, insecure bool) grpc.DialOption {
	option := grpc.WithTransportCredentials(credentials.NewTLS(
		&tls.Config{
			ServerName: host,
		},
	))
	if insecure {
		option = grpc.WithInsecure()
	}
	return option
}
