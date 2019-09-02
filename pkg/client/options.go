package client

import (
	"crypto/tls"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// DialOptions gives options that manage TLS options,
// retries and exponential backoff in calls to a service
func DialOptions(
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
