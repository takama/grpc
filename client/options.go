package client

import (
	"context"
	"crypto/tls"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
)

// DialOptions gives options that manage TLS options,
// retries and exponential backoff in calls to a service.
func DialOptions(cfg *Config, opts ...grpc.DialOption) []grpc.DialOption {
	tlsOption := grpc.WithTransportCredentials(credentials.NewTLS(
		&tls.Config{
			ServerName: cfg.Host,
		},
	))
	if cfg.Insecure {
		tlsOption = grpc.WithInsecure()
	}

	if cfg.Retry.Active {
		opts = append(opts, grpc.WithConnectParams(grpc.ConnectParams{
			MinConnectTimeout: time.Duration(cfg.Timeout) * time.Second,
			Backoff: backoff.Config{
				BaseDelay:  time.Duration(cfg.Retry.Backoff.Delay.Min) * time.Second,
				Multiplier: cfg.Retry.Backoff.Multiplier,
				Jitter:     cfg.Retry.Backoff.Jitter,
				MaxDelay:   time.Duration(cfg.Retry.Backoff.Delay.Max) * time.Second,
			},
		}))
	}

	return append(opts,
		tlsOption,
		grpc.WithDefaultCallOptions(
			grpc.WaitForReady(cfg.WaitForReady),
		),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                time.Duration(cfg.Keepalive.Time) * time.Second,
			Timeout:             time.Duration(cfg.Keepalive.Timeout) * time.Second,
			PermitWithoutStream: cfg.Keepalive.Force,
		}),
	)
}

type tokenAuth struct {
	token string
}

// Return value is mapped to request headers.
func (t tokenAuth) GetRequestMetadata(ctx context.Context, in ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": "Bearer " + t.token,
	}, nil
}

func (tokenAuth) RequireTransportSecurity() bool {
	return true
}

// TokenOption gives gRPC token option.
func TokenOption(active bool, path string) grpc.DialOption {
	if !active {
		return new(grpc.EmptyDialOption)
	}

	token, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("Read token error: %s", err)
	}

	return grpc.WithPerRPCCredentials(
		tokenAuth{
			token: strings.TrimSpace(string(token)),
		},
	)
}
