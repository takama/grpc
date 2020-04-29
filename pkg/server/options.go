package server

import (
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

// Options gives options that manage server connection parameters.
func Options(cfg *Config, opts ...grpc.ServerOption) []grpc.ServerOption {
	return append(opts,
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle:     time.Duration(cfg.Connection.Idle) * time.Second,
			MaxConnectionAge:      time.Duration(cfg.Connection.Age) * time.Second,
			MaxConnectionAgeGrace: time.Duration(cfg.Connection.Grace) * time.Second,
			Time:                  time.Duration(cfg.Connection.Keepalive.Time) * time.Second,
			Timeout:               time.Duration(cfg.Connection.Keepalive.Timeout) * time.Second,
		}),
	)
}
