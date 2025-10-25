package config

import "time"

type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
	PoolSize int
	TTL      time.Duration
}

func DefaultRedisConfig() RedisConfig {
	return RedisConfig{
		Host:     "localhost",
		Port:     6379,
		Password: "",
		DB:       0,
		PoolSize: 10,
		TTL:      24 * time.Hour,
	}
}
