package config

import (
	"bytes"
	"encoding/json"
	"log"
	"strings"
	"testing"

	"github.com/integralist/go-findroot/find"
	"github.com/spf13/viper"
)

// Load system env config
func Load() (*Config, error) {
	/**
	|-------------------------------------------------------------------------
	| hacking to load reflect structure config into env
	|-----------------------------------------------------------------------*/
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	if testing.Testing() {
		root, err := find.Repo()
		if err != nil {
			return nil, err
		}
		viper.AddConfigPath(root.Path)
	}

	viper.AddConfigPath("./")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "__"))
	viper.AutomaticEnv()
	/**
	|-------------------------------------------------------------------------
	| You should set default config value here
	| 1. Populate the default value in (Source code)
	| 2. Then merge from config (YAML) and OS environment
	|-----------------------------------------------------------------------*/
	c := loadDefaultConfig()
	if configBuffer, err := json.Marshal(c); err != nil {
		log.Printf("[CONFIG] Failed to marshal default config: %v", err)
		return nil, err
	} else if err := viper.ReadConfig(bytes.NewBuffer(configBuffer)); err != nil {
		log.Printf("[CONFIG] Failed to read default config: %v", err)
		return nil, err
	}
	if err := viper.MergeInConfig(); err != nil {
		log.Printf("[CONFIG] Failed to merge config file (using defaults): %v", err)
	}
	// Populate all config again
	err := viper.Unmarshal(c)
	return c, err
}
