package config

import (
	"fmt"

	"github.com/joeshaw/envdecode"
	log "github.com/sirupsen/logrus"
	"github.com/vivangkumar/skv/internal/server"
)

// Config represents the configuration for skv as a whole.
type Config struct {
	// Server configuration.
	Server server.Config

	// Log represents Logging configuration.
	Log struct {
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

// ServerConfig returns server specific config
func (c *Config) ServerConfig() server.Config {
	return c.Server
}
