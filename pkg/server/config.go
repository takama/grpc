package server

// Config contains params to setup server.
type Config struct {
	Port       int
	Gateway    Gateway
	Connection Connection
}

// Gateway contains params to setup gateway.
type Gateway struct {
	Port int
}

// Certificates contains path to certificates and key.
type Certificates struct {
	Crt string
	Key string
}

// Connection parameters.
type Connection struct {
	Idle      int
	Age       int
	Grace     int
	Keepalive Keepalive
}

// Keepalive connection parameters.
type Keepalive struct {
	Time    int
	Timeout int
}
