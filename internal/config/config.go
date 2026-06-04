// Package config contains logic for the handling of config files.
package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
	"github.com/tx3stn/vrsn/internal/flags"
)

type (
	// Config represents the options available in the config file.
	Config struct {
		Bump    BumpOpts  `toml:"bump"`
		Check   CheckOpts `toml:"check"`
		Files   []string  `toml:"files"`
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
func Get(fileFlag string) (Config, error) {
	var file string

	var err error

	if fileFlag == "" {
		file, err = FindConfigFile()
		if err != nil {
			return Config{}, err
		}
	} else {
		file = fileFlag
	}

	if file == "" {
		return Config{
			Bump: BumpOpts{
				Commit:    flags.Commit,
				CommitMsg: flags.CommitMsg,
				GitTag:    flags.GitTag,
				TagMsg:    flags.TagMsg,
			},
			Check: CheckOpts{
				BaseBranch: flags.BaseBranch,
			},
			Files:   filesFromFlag(flags.VersionFile),
			Verbose: flags.Verbose,
		}, nil
	}

	content, err := os.ReadFile(filepath.Clean(file))
	if err != nil {
		return Config{}, fmt.Errorf("error reading config file: %w", err)
	}

	var conf Config
	if err = toml.Unmarshal(content, &conf); err != nil {
		return Config{}, fmt.Errorf("error unmashalling config from file: %w", err)
	}

	// The files in the config file take precedence, but when none are
	// configured the --file flag still applies.
	if len(conf.Files) == 0 {
		conf.Files = filesFromFlag(flags.VersionFile)
	}

	return conf, nil
}

// filesFromFlag converts the --file flag value into the config files list.
func filesFromFlag(versionFile string) []string {
	if versionFile == "" {
		return nil
	}

	return []string{versionFile}
}

// FindConfigFile checks the expected paths for a vrsn config file and returns the
// path to it if found.
// The paths are checked in the order of precedence:
//   - current directory (project level config)
//   - XDG_CONFIG_DIR
//   - HOME/.config
func FindConfigFile() (string, error) {
	paths := []string{"."}

	if xdg, ok := os.LookupEnv("XDG_CONFIG_DIR"); ok {
		paths = append(paths, xdg)
	}

	if home, ok := os.LookupEnv("HOME"); ok {
		paths = append(paths, filepath.Join(home, ".config"))
	}

	configFileName := "vrsn.toml"

	for _, path := range paths {
		file := filepath.Join(path, configFileName)
		if _, err := os.Stat(file); os.IsNotExist(err) {
			// no config file at location, continue looking.
			continue
		}

		return file, nil
	}

	return "", nil
}
