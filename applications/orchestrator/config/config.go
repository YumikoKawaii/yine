package config

import (
	"github.com/YumikoKawaii/shared/logger"
	"github.com/YumikoKawaii/shared/mysql"
	"github.com/YumikoKawaii/shared/redis"
	"github.com/YumikoKawaii/shared/tracer"
	"yumiko_kawaii.com/yine/applications/orchestrator/server"
)

type Config struct {
	Server       server.Config
	Logger       logger.Configuration
	MysqlCfg     mysql.Config
	RedisCfg     redis.Config
	TracerConfig tracer.Configuration
}

func loadDefaultConfig() *Config {
	c := &Config{
		Server: server.DefaultConfig(),
		Logger: logger.DefaultConfig(),
		MysqlCfg: mysql.Config{
			Username: "root",
			Password: "password",
			Host:     "localhost",
			Port:     3306,
			Database: "orchestrator",
		},
		RedisCfg: redis.Config{
			Address:       "localhost:6379",
			EnableTracing: true,
		},
		TracerConfig: *tracer.DefaultConfig(),
	}
	return c
}
