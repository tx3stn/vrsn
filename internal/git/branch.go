// Package git contains logic for git interactions.
package git

import (
	"fmt"
)

// CurrentBranch gets the name of the current branch.
func CurrentBranch(dir string) (string, error) {
	// e.g.: git rev-parse --abrev-ref HEAD
	return gitCommand(
		dir,
		"error trying to get current git branch name",
		"rev-parse", "--abbrev-ref", "HEAD",
	)
}

// VersionAtBranch returns the version file contents from the specific branch.
func VersionAtBranch(dir string, branchName string, versionFile string) (string, error) {
	// e.g.: git --no-pager show main:VERSION
	return gitCommand(
		dir,
		fmt.Sprintf("error trying to read %s from %s", versionFile, branchName),
		"--no-pager", "show", fmt.Sprintf("%s:%s", branchName, versionFile),
	)
}
