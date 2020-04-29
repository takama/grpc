package client

import (
	"context"
	"fmt"

	"github.com/takama/grpc/client"
	"github.com/takama/grpc/contracts/echo"
	"github.com/takama/grpc/contracts/info"

	"github.com/golang/protobuf/ptypes/empty"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// Client provides access to the service using client connection.
type Client struct {
	cl  *client.Client
	log *zap.Logger
}

// New gives a Client.
func New(cfg *client.Config, log *zap.Logger) (*Client, error) {
	// Set up a connection to the server.
	// Create new gRPC client
	cl, err := client.New(cfg, log)
	if err != nil {
		return nil, err
	}

	return &Client{
		cl:  cl,
		log: log,
	}, nil
}

// Info command.
func (c *Client) Info(ctx context.Context) error {
	// Set up a connection to the server.
	info, err := info.NewInfoClient(c.cl.Connection()).
		GetInfo(c.cl.Context(ctx), new(empty.Empty))
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

// Ping command.
func (c *Client) Ping(ctx context.Context, message string, count int) error {
	md := new(metadata.MD)

	cl := echo.NewEchoClient(c.cl.Connection())

	for idx := 1; idx <= count; idx++ {
		response, err := cl.Ping(c.cl.Context(ctx), &echo.Request{
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

// Reverse command.
func (c *Client) Reverse(ctx context.Context, message string, count int) error {
	md := new(metadata.MD)

	cl := echo.NewEchoClient(c.cl.Connection())

	for idx := 1; idx <= count; idx++ {
		response, err := cl.Reverse(c.cl.Context(ctx), &echo.Request{
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

// Shutdown closes active Client connections.
func (c *Client) Shutdown(ctx context.Context) error {
	return c.cl.Shutdown(ctx)
}
