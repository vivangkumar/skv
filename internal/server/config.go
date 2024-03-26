package server

// Config stores the config for the skv server.
type Config struct {
	// Addr is the TCP address for the server.
	Addr string `env:"SERVER_ADDR,default=:2303"`
}
