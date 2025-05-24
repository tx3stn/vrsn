// Package flags holds logics for use of CLI flags.
package flags

var (
	// ConfigFile is the variable for the CLI flag `--config` used to specify a config
	// file not stored in the default location.
	ConfigFile string
	// Now is the variable for the CLI flag `--now`.
	Now string
	// VersionFile is the variable for the CLI flag `--file` to provide a specific
	// version file path, rather than having vrsn try and work out what to use.
	VersionFile string
	// Was is the variable for the CLI flag `--was`.
	Was string
)
