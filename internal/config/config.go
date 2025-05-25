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
		Bump    BumpOpts  `toml:"bump"`
		Check   CheckOpts `toml:"check"`
		Verbose bool      `toml:"verbose"`
	}

	// BumpOpts are the vrsn bump specific options in the config file.
	BumpOpts struct {
		Commit    bool   `toml:"commit"`
		CommitMsg string `toml:"commit-msg"`
		GitTag    bool   `toml:"git-tag"`
		TagMsg    string `toml:"tag-msg"`
	}

	// CheckOpts are the vrsn check specific options in the config file.
	CheckOpts struct {
		BaseBranch string `toml:"base-branch"`
	}
)

// Get returns the config.
func Get() (Config, error) {
	usingConfigFile := true

	if err := viper.ReadInConfig(); err != nil {
		//nolint:errorlint
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			usingConfigFile = false

			if viper.GetBool("verbose") {
				fmt.Println("no config file found")
			}
		} else {
			// Config file was found but another error was produced
			return Config{}, fmt.Errorf("error reading config file: %w", err)
		}
	}

	conf := viper.AllSettings()

	tomlContent, err := toml.Marshal(conf)
	if err != nil {
		return Config{}, fmt.Errorf("error marshalling config file: %w", err)
	}

	parsedConfig := Config{}
	if err := toml.Unmarshal(tomlContent, &parsedConfig); err != nil {
		return Config{}, fmt.Errorf("error unmarshalling config file: %w", err)
	}

	parsedConfig.setDefaults(usingConfigFile)

	return parsedConfig, nil
}

func (c *Config) setDefaults(useConfigFile bool) {
	if !useConfigFile {
		c.Bump.Commit = viper.GetBool("commit")
		c.Bump.CommitMsg = viper.GetString("commit-msg")
	}

	if c.Check.BaseBranch == "" {
		c.Check.BaseBranch = "main"
	}

	if c.Bump.CommitMsg == "" {
		c.Bump.CommitMsg = "bump version"
	}
}
