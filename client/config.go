package client

// Config contains params to setup client connection
type Config struct {
	Scheme       string
	Host         string
	Sockets      []string
	Balancer     string
	Insecure     bool
	EnvoyProxy   bool
	WaitForReady bool
	Timeout      int
	Keepalive    Keepalive
	Retry        Retry
}

// Retry config describes retry parameters
type Retry struct {
	Active  bool
	Envoy   Envoy
	Backoff Backoff
}

// Envoy contains envoy proxy retry parameters
type Envoy struct {
	Reason  Reason
	Count   int
	Timeout int
}

// Reason config describes reasons to retry
type Reason struct {
	Primary string
	GRPC    string
}

// Keepalive connection parameters
type Keepalive struct {
	Time    int
	Timeout int
	Force   bool
}

// Backoff is exponential backoff attributes
type Backoff struct {
	Multiplier float64
	Jitter     float64
	Delay      Delay
}

// Delay are the amounts of time to backoff after the first failure
// and the upper bound of backoff delay
type Delay struct {
	Min int
	Max int
}
