// Package config contains logic for the handling of config files.
package config

import (
	"fmt"

	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/viper"
)

type (
	// Config represents the options available in the config file.
	Config struct {
		Commit    bool   `toml:"commit"`
		CommitMsg string `toml:"commit-msg"`
		Verbose   bool   `toml:"verbose"`
	}
)

// Get returns the config.
func Get() (Config, error) {
	conf := viper.AllSettings()

	tomlContent, err := toml.Marshal(conf)
	if err != nil {
		return Config{}, fmt.Errorf("error marshalling config file: %w", err)
	}

	parsedConfig := Config{}
	if err := toml.Unmarshal(tomlContent, &parsedConfig); err != nil {
		return Config{}, fmt.Errorf("error unmarshalling config file: %w", err)
	}

	return parsedConfig, nil
}
