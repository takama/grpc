package client

import (
	"crypto/tls"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
)

// DialOptions gives options that manage TLS options,
// retries and exponential backoff in calls to a service
func DialOptions(cfg *Config, opts ...grpc.DialOption) []grpc.DialOption {
	tlsOption := grpc.WithTransportCredentials(credentials.NewTLS(
		&tls.Config{
			ServerName: cfg.Host,
		},
	))
	if cfg.Insecure {
		tlsOption = grpc.WithInsecure()
	}

	if !cfg.EnvoyProxy {
		opts = append(opts, grpc.WithBackoffMaxDelay(time.Duration(cfg.Timeout)*time.Second))
	}

	return append(opts,
		tlsOption,
		grpc.WithDefaultCallOptions(
			grpc.WaitForReady(cfg.WaitForReady),
		),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:    time.Duration(cfg.Keepalive.Time) * time.Second,
			Timeout: time.Duration(cfg.Keepalive.Timeout) * time.Second,
		}),
	)
}
