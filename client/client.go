package client

import (
	"context"
	"fmt"
	"strconv"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
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

	// Set up a connection to the server.
	conn, err := grpc.Dial(
		fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		DialOptions(cfg, opts...)...,
	)
	if err != nil {
		return nil, err
	}

	log.Info(
		"Connected with config:",
		zap.String("host", cfg.Host),
		zap.Int("port", cfg.Port),
		zap.Bool("insecure", cfg.Insecure),
		zap.Bool("wait for ready", cfg.WaitForReady),
		zap.Int("timeout (s)", cfg.Timeout),
		zap.String("Retry reason", cfg.Retry.Reason.Primary),
		zap.String("Retry reason for gRPC", cfg.Retry.Reason.GRPC),
		zap.String("Retries count", strconv.Itoa(cfg.Retry.Count)),
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
