// Package flags holds logics for use of CLI flags.
package flags

var (
	// BaseBranch is the variable for the CLI flag `--base-branch` so you can set
	// your git base branch if it's anything other than `main`.
	BaseBranch string
	// ConfigFile is the variable for the CLI flag `--config` used to specify a config
	// file not stored in the default location.
	ConfigFile string
	// GitTag is the variable for the CLI flag `--git-tag` used to read the version from
	// the git tags.
	GitTag bool
	// Now is the variable for the CLI flag `--now`.
	Now string
	// TagMsg is the variable for the CLI flag `--tag-msg` to add a custom git tag
	// message. Only used with the `--git-tag` flag.
	TagMsg string
	// Verbose is the variable for the CLI flag `--verbose` to enable debug log output.
	Verbose bool
	// VersionFile is the variable for the CLI flag `--file` to provide a specific
	// version file path, rather than having vrsn try and work out what to use.
	VersionFile string
	// Was is the variable for the CLI flag `--was`.
	Was string
)
