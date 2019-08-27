package boot

import (
	"crypto/tls"
	"time"

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

// PrepareDialOptions gives options that manage TLS options,
// retries and exponential backoff in calls to a service
func PrepareDialOptions(
	host string, insecure, waitForReady bool, maxDelay time.Duration,
	opts ...grpc.DialOption,
) []grpc.DialOption {
	tlsOption := grpc.WithTransportCredentials(credentials.NewTLS(
		&tls.Config{
			ServerName: host,
		},
	))
	if insecure {
		tlsOption = grpc.WithInsecure()
	}
	return append(opts,
		tlsOption,
		grpc.WithDefaultCallOptions(
			grpc.WaitForReady(waitForReady),
		),
		grpc.WithBackoffMaxDelay(maxDelay),
	)
}
