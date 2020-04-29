package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/takama/grpc/client"
	"github.com/takama/grpc/pkg/config"
	"github.com/takama/grpc/pkg/info"
	"github.com/takama/grpc/pkg/logger"
	"github.com/takama/grpc/pkg/server"
	"github.com/takama/grpc/pkg/system"
	"github.com/takama/grpc/pkg/version"

	"go.uber.org/zap"
)

// Run the service.
func Run(cfg *config.Config) error {
	// Setup zap logger
	log := logger.New(&cfg.Logger)
	defer func(*zap.Logger) {
		if err := log.Sync(); err != nil {
			log.Error(err.Error())
		}
	}(log)

	log.Info(
		config.ServiceName,
		zap.String("version", version.RELEASE+"-"+version.COMMIT+"-"+version.BRANCH),
	)

	// Create new gRPC client connection.
	cl, err := client.New(&cfg.Client, log)
	if err != nil {
		return err
	}

	// Create new core server.
	srv, err := server.New(context.Background(), cl, &cfg.Server, log)
	if err != nil {
		return err
	}

	// Register info/health-check service.
	is := info.NewService(log)
	is.RegisterLivenessProbe(srv.LivenessProbe)
	is.RegisterReadinessProbe(srv.ReadinessProbe)

	// Run info/health-check service
	infoServer := is.Run(fmt.Sprintf(":%d", cfg.Info.Port))

	// Setup operator with info/health-check server, core server and data store.
	operator := system.NewOperator(
		&cfg.System,
		srv,
		cl,
		infoServer,
	)

	// Run core server.
	go func() {
		if err := srv.Run(context.Background()); err != nil {
			// Check for known errors
			if err != context.DeadlineExceeded &&
				err != context.Canceled &&
				err != http.ErrServerClosed {
				log.Fatal(err.Error())
			}

			log.Error(err.Error())
		}
	}()

	// Wait for signals.
	return system.NewSignals().Wait(log, operator)
}
