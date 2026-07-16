// Package config contains logic for the handling of config files.
package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
	"github.com/tx3stn/vrsn/internal/flags"
)

// FlagChecker reports whether a flag was explicitly set on the command line.
// *pflag.FlagSet (returned by cobra's cmd.Flags()) satisfies it without the
// config package needing to depend on pflag directly.
type FlagChecker interface {
	Changed(name string) bool
}

type (
	// Config represents the options available in the config file.
	Config struct {
		Bump    BumpOpts  `toml:"bump"`
		Check   CheckOpts `toml:"check"`
		Set     SetOpts   `toml:"set"`
		Files   []string  `toml:"files"`
		Verbose bool      `toml:"verbose"`
	}

	// BumpOpts are the vrsn bump specific options in the config file.
	BumpOpts struct {
		AndroidVersionCode bool   `toml:"android-version-code"`
		Commit             bool   `toml:"commit"`
		CommitMsg          string `toml:"commit-msg"`
		GitTag             bool   `toml:"git-tag"`
		TagMsg             string `toml:"tag-msg"`
	}

	// CheckOpts are the vrsn check specific options in the config file.
	CheckOpts struct {
		BaseBranch string `toml:"base-branch"`
	}

	// SetOpts are the vrsn set specific options in the config file.
	SetOpts struct {
		AndroidVersionCode bool `toml:"android-version-code"`
	}
)

// Get returns the effective config: values from a config file (if one is
// found) layered over the CLI flag defaults, with any flags explicitly passed
// on the command line taking precedence over the config file.
// The one exception is the documented behaviour that `files` in the config
// file takes precedence over the --file flag.
func Get(fileFlag string, flagSet FlagChecker) (Config, error) {
	conf := Config{
		Bump: BumpOpts{
			AndroidVersionCode: flags.AndroidVersionCode,
			Commit:             flags.Commit,
			CommitMsg:          flags.CommitMsg,
			GitTag:             flags.GitTag,
			TagMsg:             flags.TagMsg,
		},
		Check: CheckOpts{
			BaseBranch: flags.BaseBranch,
		},
		Set: SetOpts{
			AndroidVersionCode: flags.AndroidVersionCode,
		},
		Files:   filesFromFlag(flags.VersionFile),
		Verbose: flags.Verbose,
	}

	file := fileFlag
	if file == "" {
		var err error

		file, err = FindConfigFile()
		if err != nil {
			return Config{}, err
		}
	}

	if file == "" {
		return conf, nil
	}

	content, err := os.ReadFile(filepath.Clean(file))
	if err != nil {
		return Config{}, fmt.Errorf("error reading config file: %w", err)
	}

	// Unmarshalling over the flag based config means options missing from the
	// config file keep the flag default, while `files` from the config file
	// still replaces the --file flag value.
	if err = toml.Unmarshal(content, &conf); err != nil {
		return Config{}, fmt.Errorf("error unmarshalling config from file: %w", err)
	}

	applyChangedFlags(&conf, flagSet)

	return conf, nil
}

// applyChangedFlags overrides config file values with any flags that were
// explicitly set on the command line, so passing a flag always works
// regardless of which config file is found.
// Flags not registered on the current command report as unchanged.
func applyChangedFlags(conf *Config, flagSet FlagChecker) {
	if flagSet == nil {
		return
	}

	if flagSet.Changed("android-version-code") {
		conf.Bump.AndroidVersionCode = flags.AndroidVersionCode
		conf.Set.AndroidVersionCode = flags.AndroidVersionCode
	}

	if flagSet.Changed("commit") {
		conf.Bump.Commit = flags.Commit
	}

	if flagSet.Changed("commit-msg") {
		conf.Bump.CommitMsg = flags.CommitMsg
	}

	if flagSet.Changed("git-tag") {
		conf.Bump.GitTag = flags.GitTag
	}

	if flagSet.Changed("tag-msg") {
		conf.Bump.TagMsg = flags.TagMsg
	}

	if flagSet.Changed("base-branch") {
		conf.Check.BaseBranch = flags.BaseBranch
	}

	if flagSet.Changed("verbose") {
		conf.Verbose = flags.Verbose
	}
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
//   - XDG_CONFIG_HOME (with XDG_CONFIG_DIR still supported for backwards
//     compatibility)
//   - HOME/.config
func FindConfigFile() (string, error) {
	paths := []string{"."}

	if xdg, ok := os.LookupEnv("XDG_CONFIG_HOME"); ok {
		paths = append(paths, xdg)
	}

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
