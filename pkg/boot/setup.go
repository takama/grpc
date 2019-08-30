package boot

import (
	"github.com/takama/grpc/pkg/config"
	"github.com/takama/grpc/pkg/helper"
	"github.com/takama/grpc/pkg/logger"

	"go.uber.org/zap"
)

// Setup the configuration, logger
func Setup() (*config.Config, *zap.Logger) {
	cfg, err := config.New()
	helper.LogF("Load config error", err)
	// Setup zap logger
	log := logger.New(&cfg.Logger)

	return cfg, log
}
