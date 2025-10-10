package config

import (
	"yumiko_kawaii.com/yine/applications/orchestrator/pkg/logger"
	"yumiko_kawaii.com/yine/applications/orchestrator/server"
)

type Config struct {
	Server server.Config
	Logger logger.Configuration
}

func loadDefaultConfig() *Config {
	c := &Config{
		Server: server.DefaultConfig(),
		Logger: logger.DefaultConfig(),
	}
	return c
}
