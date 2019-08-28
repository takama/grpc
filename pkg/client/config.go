package client

import (
	"time"
)

// Config contains params to setup client connection
type Config struct {
	Host         string
	Port         int
	Insecure     bool
	WaitForReady bool
	BackOffDelay time.Duration
}
