package config

import (
	"fmt"

	"github.com/joeshaw/envdecode"
	log "github.com/sirupsen/logrus"
	"github.com/vivangkumar/skv/pkg/node"
)

// Config represents the configuration for SKV as a whole
type Config struct {
	// Node configuration
	Node node.Config
	Log  struct {
		Level log.Level `env:"LOG_LEVEL,default=info"`
	}
}

// New returns an instance of the SKV configuration
func New() (*Config, error) {
	var cfg Config

	err := envdecode.StrictDecode(&cfg)
	if err != nil {
		return nil, fmt.Errorf("envdecode: %w", err)
	}

	return &cfg, nil
}

func (c *Config) NodeConfig() node.Config {
	return c.Node
}
