package client

import (
	"context"
	"fmt"
	"time"

	"github.com/takama/grpc/contracts/info"

	"github.com/golang/protobuf/ptypes/empty"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// Info command
func Info(cfg *Config, log *zap.Logger, opts ...grpc.DialOption) error {
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

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	info, err := info.NewInfoClient(conn).GetInfo(ctx, new(empty.Empty))
	if err != nil {
		return err
	}

	log.Info(
		"Info",
		zap.String("version", info.Version),
		zap.String("date", info.Date),
		zap.String("repo", info.Repo),
	)

	return nil
}
