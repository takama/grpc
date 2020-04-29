package config_test

import (
	"testing"

	"github.com/takama/grpc/pkg/config"
)

func TestConfig(t *testing.T) {
	cfg, err := config.New()
	if err != nil {
		t.Error("Expected loading of config, got", err)
	}

	if cfg.Server.Port != config.DefaultServerPort {
		t.Errorf("Expected %d, got %d", config.DefaultServerPort, cfg.Server.Port)
	}

	if cfg.Info.Port != config.DefaultInfoPort {
		t.Errorf("Expected %d, got %d", config.DefaultInfoPort, cfg.Info.Port)
	}

	if cfg.Logger.Level != config.DefaultLoggerLevel {
		t.Errorf("Expected %d, got %d", config.DefaultLoggerLevel, cfg.Logger.Level)
	}
}
