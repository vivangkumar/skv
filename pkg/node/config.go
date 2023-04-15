package node

// Config stores the config for a single `skv` node
type Config struct {
	// Addr is the TCP address for the node
	Addr string `env:"NODE_ADDR,default=:2303"`
}
