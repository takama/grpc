package server

// Config contains params to setup server
type Config struct {
	Port    int
	Gateway Gateway
}

// Gateway contains params to setup gateway
type Gateway struct {
	Port int
}

// Certificates contains path to certificates and key
type Certificates struct {
	Crt string
	Key string
}
