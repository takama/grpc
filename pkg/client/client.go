package client

import (
	"context"
	"fmt"
	"time"

	"github.com/takama/grpc/contracts/echo"
	"github.com/takama/grpc/contracts/info"

	"github.com/golang/protobuf/ptypes/empty"
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
	// Set up a connection to the server.
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port), opts...)
	if err != nil {
		return nil, err
	}
	log.Info(
		"Connected with config:",
		zap.String("host", cfg.Host),
		zap.Int("port", cfg.Port),
		zap.Bool("insecure", cfg.Insecure),
		zap.Bool("wait for ready", cfg.WaitForReady),
		zap.Any("back off delay", cfg.BackOffDelay),
	)

	return &Client{
		cfg:  cfg,
		log:  log,
		conn: conn,
	}, nil
}

// Info command
func (c *Client) Info() error {
	// Set up a connection to the server.

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

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
func (c *Client) Ping(message string, count int) error {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*300)
	defer cancel()

	metadata := new(metadata.MD)

	cl := echo.NewEchoClient(c.conn)
	for idx := 1; idx <= count; idx++ {
		response, err := cl.Ping(ctx, &echo.Request{
			Content: fmt.Sprintf("%s: %d", message, idx),
		}, grpc.Header(metadata))
		if err != nil {
			return err
		}
		c.log.Info(
			"ping",
			zap.String("message", response.Content),
			zap.Any("hostname", metadata.Get("hostname")),
		)
	}

	return nil
}

// Reverse command
func (c *Client) Reverse(message string, count int) error {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*300)
	defer cancel()

	metadata := new(metadata.MD)

	cl := echo.NewEchoClient(c.conn)
	for idx := 1; idx <= count; idx++ {
		response, err := cl.Reverse(ctx, &echo.Request{
			Content: fmt.Sprintf("%s: %d", message, idx),
		}, grpc.Header(metadata))
		if err != nil {
			return err
		}
		c.log.Info(
			"reverse",
			zap.String("message", response.Content),
			zap.Any("hostname", metadata.Get("hostname")),
		)
	}

	return nil
}

// Shutdown closes active Client connections
func (c *Client) Shutdown() {
	if err := c.conn.Close(); err != nil {
		c.log.Error("Connection close error:", zap.Error(err))
	}
}
