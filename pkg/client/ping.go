package client

import (
	"context"
	"fmt"
	"time"

	"github.com/takama/grpc/contracts/echo"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// Ping command
func Ping(cfg *Config, log *zap.Logger, message string, count int, opts ...grpc.DialOption) error {
	// Set up a connection to the server.
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port), opts...)
	if err != nil {
		return err
	}
	defer func(log *zap.Logger) {
		if err := conn.Close(); err != nil {
			log.Error("Connection close error:", zap.Error(err))
		}
	}(log)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*300)
	defer cancel()

	metadata := new(metadata.MD)

	cl := echo.NewEchoClient(conn)
	for idx := 1; idx <= count; idx++ {
		response, err := cl.Ping(ctx, &echo.Request{
			Content: fmt.Sprintf("%s: %d", message, idx),
		}, grpc.Header(metadata))
		if err != nil {
			return err
		}
		log.Info(
			"ping",
			zap.String("message", response.Content),
			zap.Any("hostname", metadata.Get("hostname")),
		)
	}

	return nil
}
