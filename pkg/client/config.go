package client

// Config contains params to setup client connection
type Config struct {
	Host         string
	Port         int
	Insecure     bool
	EnvoyProxy   bool
	WaitForReady bool
	BackOffDelay int
	Retry        Retry
}

// Retry config describes retry parameters
type Retry struct {
	Reason  Reason
	Count   int
	Timeout int
}

// Reason config describes reasons to retry
type Reason struct {
	Primary string
	GRPC    string
}
