package git

import "fmt"

// Add adds the version file to the git staging area.
func Add(dir string, file string) (string, error) {
	// e.g.: git add package.json
	return gitCommand(
		dir,
		fmt.Sprintf("error staging %s, file will not be committed", file),
		"add", file,
	)
}

// Commit commits just the version file with the provided commit message.
func Commit(dir string, file string, msg string) (string, error) {
	// e.g.: git commit package.json -m "bump version"
	return gitCommand(
		dir,
		fmt.Sprintf("error committing %s", file),
		"commit", file, "-m", msg,
	)
}
