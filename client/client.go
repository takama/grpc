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

const (
	// It represents timeout in 1000 ms = 1s
	timeoutSecond int = 1000
)

// Client provides access to the service using client connection
type Client struct {
	cfg  *Config
	log  *zap.Logger
	conn *grpc.ClientConn
}

// New gives a Client
func New(cfg *Config, log *zap.Logger, opts ...grpc.DialOption) (*Client, error) {
	var config string

	if len(cfg.Sockets) > 0 {
		if cfg.Balancer == roundrobin.Name {
			config += fmt.Sprintf(`{
			"loadBalancingConfig": [{"%v": {}}]
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
		zap.Bool("used Envoy proxy", cfg.EnvoyProxy),
		zap.Bool("wait for ready", cfg.WaitForReady),
		zap.Int("timeout (s)", cfg.Timeout),
		zap.Int("keepalive time (s)", cfg.Keepalive.Time),
		zap.Int("keepalive timeout (s)", cfg.Keepalive.Timeout),
		zap.Bool("keepalive force ping", cfg.Keepalive.Force),
		zap.Bool("Retry is active", cfg.Retry.Active),
		zap.String("envoy proxy retry reason", cfg.Retry.Envoy.Reason.Primary),
		zap.String("envoy proxy retry reason for gRPC", cfg.Retry.Envoy.Reason.GRPC),
		zap.Int("envoy proxy retries count", cfg.Retry.Envoy.Count),
		zap.Int("envoy proxy retry timeout (ms)", cfg.Retry.Envoy.Timeout*timeoutSecond),
		zap.Float64("backoff multiplier", cfg.Retry.Backoff.Multiplier),
		zap.Float64("backoff jitter", cfg.Retry.Backoff.Jitter),
		zap.Int("backoff min delay", cfg.Retry.Backoff.Delay.Min),
		zap.Int("backoff max delay", cfg.Retry.Backoff.Delay.Max),
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
	if c.cfg.EnvoyProxy && c.cfg.Retry.Active {
		return metadata.AppendToOutgoingContext(ctx,
			"x-envoy-retry-on", c.cfg.Retry.Envoy.Reason.Primary,
			"x-envoy-retry-grpc-on", c.cfg.Retry.Envoy.Reason.GRPC,
			"x-envoy-max-retries", strconv.Itoa(c.cfg.Retry.Envoy.Count),
			"x-envoy-upstream-rq-timeout-ms", strconv.Itoa(c.cfg.Timeout*timeoutSecond),
			"x-envoy-upstream-rq-per-try-timeout-ms", strconv.Itoa(c.cfg.Retry.Envoy.Timeout*timeoutSecond),
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
