package client

import (
	"context"
	"fmt"
	"strconv"

	"github.com/takama/grpc/contracts/echo"
	"github.com/takama/grpc/contracts/info"

	"github.com/golang/protobuf/ptypes/empty"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// Client provides access to the service using client connection
type Client struct {
	ctx  context.Context
	cfg  *Config
	log  *zap.Logger
	conn *grpc.ClientConn
}

// New gives a Client
func New(ctx context.Context, cfg *Config, log *zap.Logger) (*Client, error) {
	// Set up a connection to the server.
	conn, err := grpc.Dial(
		fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		DialOptions(cfg)...,
	)
	if err != nil {
		return nil, err
	}

	if cfg.EnvoyProxy {
		ctx = metadata.AppendToOutgoingContext(ctx,
			"x-envoy-retry-on", cfg.Retry.Reason.Primary,
			"x-envoy-retry-grpc-on", cfg.Retry.Reason.GRPC,
			"x-envoy-max-retries", strconv.Itoa(cfg.Retry.Count),
			"x-envoy-upstream-rq-timeout-ms", strconv.Itoa(cfg.Timeout*1000),
			"x-envoy-upstream-rq-per-try-timeout-ms", strconv.Itoa(cfg.Retry.Timeout*1000),
		)
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
		ctx:  ctx,
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
func (c *Client) Context() context.Context {
	return c.ctx
}

// Info command
func (c *Client) Info(ctx context.Context) error {
	// Set up a connection to the server.
	info, err := info.NewInfoClient(c.conn).GetInfo(ctx, new(empty.Empty))
	if err != nil {
		return err
	}

	c.log.Info(
		"Info",
		zap.String("version", info.Version),
		zap.String("date", info.Date),
		zap.String("repo", info.Repo),
	)

	return nil
}

// Ping command
func (c *Client) Ping(ctx context.Context, message string, count int) error {
	md := new(metadata.MD)

	cl := echo.NewEchoClient(c.conn)

	for idx := 1; idx <= count; idx++ {
		response, err := cl.Ping(ctx, &echo.Request{
			Content: fmt.Sprintf("%s: %d", message, idx),
		}, grpc.Header(md))
		if err != nil {
			return err
		}

		c.log.Info(
			"ping",
			zap.String("message", response.Content),
			zap.Any("hostname", md.Get("hostname")),
		)
	}

	return nil
}

// Reverse command
func (c *Client) Reverse(ctx context.Context, message string, count int) error {
	md := new(metadata.MD)

	cl := echo.NewEchoClient(c.conn)

	for idx := 1; idx <= count; idx++ {
		response, err := cl.Reverse(ctx, &echo.Request{
			Content: fmt.Sprintf("%s: %d", message, idx),
		}, grpc.Header(md))
		if err != nil {
			return err
		}

		c.log.Info(
			"reverse",
			zap.String("message", response.Content),
			zap.Any("hostname", md.Get("hostname")),
		)
	}

	return nil
}

// Shutdown closes active Client connections
func (c *Client) Shutdown(ctx context.Context) error {
	if err := c.conn.Close(); err != nil {
		c.log.Error("Connection close error:", zap.Error(err))
	}

	return nil
}
