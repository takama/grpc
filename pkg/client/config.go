package client

// Config contains params to setup client connection
type Config struct {
	Host     string
	Port     int
	Insecure bool
}
