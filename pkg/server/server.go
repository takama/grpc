package server

import (
	"context"
	"fmt"
	"net"

	"github.com/takama/grpc/client"
	"github.com/takama/grpc/contracts/echo"
	"github.com/takama/grpc/contracts/info"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

// Server contains core functionality of the service.
type Server struct {
	cfg *Config
	log *zap.Logger
	srv *grpc.Server
	cl  *client.Client
	is  *infoServer
	hs  *healthServer
	es  *echoServer
}

// New creates a new core server.
func New(ctx context.Context, cl *client.Client, cfg *Config, log *zap.Logger) (*Server, error) {
	return &Server{
		cfg: cfg,
		log: log,
		cl:  cl,
		is:  new(infoServer),
		hs:  new(healthServer),
		es:  &echoServer{cl: cl, log: log},
	}, nil
}

// LivenessProbe returns liveness probe of the server.
func (s Server) LivenessProbe() error {
	return nil
}

// ReadinessProbe returns readiness probe for the server.
func (s Server) ReadinessProbe() error {
	return nil
}

// Run starts the server.
func (s *Server) Run(ctx context.Context) error {
	// Register gRPC server
	s.srv = grpc.NewServer(Options(s.cfg)...)
	info.RegisterInfoServer(s.srv, s.is)
	grpc_health_v1.RegisterHealthServer(s.srv, s.hs)
	echo.RegisterEchoServer(s.srv, s.es)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", s.cfg.Port))
	if err != nil {
		return err
	}

	return s.srv.Serve(listener)
}

// Shutdown process graceful shutdown for the server.
func (s Server) Shutdown(ctx context.Context) error {
	if s.srv != nil {
		s.srv.GracefulStop()
	}

	if s.cl != nil {
		return s.cl.Shutdown(ctx)
	}

	return nil
}
