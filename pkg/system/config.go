package system

// Config contains system parameters.
type Config struct {
	Grace Grace
}

// Grace contains attributes of graceful shutdown process.
type Grace struct {
	Period int
}
