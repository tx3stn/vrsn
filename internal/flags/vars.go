// Package flags holds logics for use of CLI flags.
package flags

var (
	// BaseBranch is the variable for the CLI flag `--base-branch` so you can set
	// your git base branch if it's anything other than `main`.
	BaseBranch string
	// Commit is the variable for the CLI flag `--commit` used to tell the `bump`
	// command to commit the version file after bumping.
	Commit bool
	// CommitMsg is the variable for the CLI flag `--commit-msg` used when
	// committing version file changes with the `bump` command.
	CommitMsg string
	// Now is the variable for the CLI flag `--now`.
	Now string
	// Was is the variable for the CLI flag `--was`.
	Was string
	// Verbose is the variable for the CLI flag `--verbose` to enable debug log output.
	Verbose bool
	// VersionFile is the variable for the CLI flag `--file` to provide a specific
	// version file path, rather than having vrsn try and work out what to use.
	VersionFile string
)
