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
