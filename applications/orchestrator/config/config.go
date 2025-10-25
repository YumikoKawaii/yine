package config

import (
	"yumiko_kawaii.com/yine/applications/orchestrator/pkg/logger"
	"yumiko_kawaii.com/yine/applications/orchestrator/server"
)

type Config struct {
	Server   server.Config
	Logger   logger.Configuration
	Redis    RedisConfig
	Database DatabaseConfig
}

func loadDefaultConfig() *Config {
	c := &Config{
		Server:   server.DefaultConfig(),
		Logger:   logger.DefaultConfig(),
		Redis:    DefaultRedisConfig(),
		Database: DefaultDatabaseConfig(),
	}
	return c
}
