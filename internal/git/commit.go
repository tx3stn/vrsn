package git

import (
	"fmt"
	"strings"
)

// Add adds the version files to the git staging area.
func Add(dir string, files ...string) (string, error) {
	// e.g.: git add package.json
	return gitCommand(
		dir,
		fmt.Sprintf("error staging %s, files will not be committed", strings.Join(files, ", ")),
		append([]string{"add"}, files...)...,
	)
}

// Commit commits just the version files with the provided commit message.
func Commit(dir string, msg string, files ...string) (string, error) {
	// e.g.: git commit package.json -m "bump version"
	args := append([]string{"commit"}, files...)
	args = append(args, "-m", msg)

	return gitCommand(
		dir,
		"error committing "+strings.Join(files, ", "),
		args...,
	)
}
