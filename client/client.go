package client

import (
	"context"
	"fmt"
	"strconv"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/resolver/manual"
)

// Client provides access to the service using client connection
type Client struct {
	cfg  *Config
	log  *zap.Logger
	conn *grpc.ClientConn
}

// New gives a Client
func New(cfg *Config, log *zap.Logger, opts ...grpc.DialOption) (*Client, error) {
	if !cfg.EnvoyProxy {
		opts = append(opts, RetryOption(cfg))
	}

	var config string

	if len(cfg.Sockets) > 0 {
		if cfg.Balancer == roundrobin.Name {
			config += fmt.Sprintf(`{
			"loadBalancingConfig": [
				{"%v": {}}
			]
		}`, roundrobin.Name)
		}
	} else {
		cfg.Scheme = "dns"
	}

	// Initialize default service config
	opts = append(opts, grpc.WithDefaultServiceConfig(config))

	// Prepare resolver according to scheme
	mr := manual.NewBuilderWithScheme(cfg.Scheme)
	mr.InitialState(prepareResolverState(cfg.Scheme, cfg.Sockets))

	resolver.Register(mr)

	// Set up a connection to the server.
	conn, err := grpc.Dial(
		fmt.Sprintf("%s:///%s", cfg.Scheme, cfg.Host),
		DialOptions(cfg, opts...)...,
	)
	if err != nil {
		return nil, err
	}

	log.Info(
		"Connected with config:",
		zap.String("scheme", cfg.Scheme),
		zap.String("host", cfg.Host),
		zap.Any("sockets", cfg.Sockets),
		zap.String("balancer", cfg.Balancer),
		zap.Bool("insecure", cfg.Insecure),
		zap.Bool("wait for ready", cfg.WaitForReady),
		zap.Int("timeout (s)", cfg.Timeout),
		zap.Int("keepalive time (s)", cfg.Keepalive.Time),
		zap.Int("keepalive timeout (s)", cfg.Keepalive.Timeout),
		zap.Bool("keepalive force ping", cfg.Keepalive.Force),
		zap.String("retry reason", cfg.Retry.Reason.Primary),
		zap.String("retry reason for gRPC", cfg.Retry.Reason.GRPC),
		zap.String("retries count", strconv.Itoa(cfg.Retry.Count)),
		zap.String("retries timeout (ms)", strconv.Itoa(cfg.Retry.Timeout*1000)),
	)

	return &Client{
		cfg:  cfg,
		log:  log,
		conn: conn,
	}, nil
}

// Connection returns gRPC connection
func (c *Client) Connection() *grpc.ClientConn {
	return c.conn
}

// Context returns context
func (c *Client) Context(ctx context.Context) context.Context {
	if c.cfg.EnvoyProxy {
		return metadata.AppendToOutgoingContext(ctx,
			"x-envoy-retry-on", c.cfg.Retry.Reason.Primary,
			"x-envoy-retry-grpc-on", c.cfg.Retry.Reason.GRPC,
			"x-envoy-max-retries", strconv.Itoa(c.cfg.Retry.Count),
			"x-envoy-upstream-rq-timeout-ms", strconv.Itoa(c.cfg.Timeout*1000),
			"x-envoy-upstream-rq-per-try-timeout-ms", strconv.Itoa(c.cfg.Retry.Timeout*1000),
		)
	}

	return ctx
}

// Shutdown closes active Client connections
func (c *Client) Shutdown(ctx context.Context) error {
	if err := c.conn.Close(); err != nil {
		c.log.Error("Connection close error:", zap.Error(err))
	}

	return nil
}

func prepareResolverState(scheme string, sockets []string) resolver.State {
	addr := make([]resolver.Address, len(sockets))

	for ind, val := range sockets {
		addr[ind] = resolver.Address{
			Addr:       val,
			ServerName: fmt.Sprintf("%s-%d", scheme, ind),
		}
	}

	return resolver.State{Addresses: addr}
}
