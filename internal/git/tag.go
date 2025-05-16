package git

import (
	"errors"
	"strings"
)

// AddTag adds the specified tag.
func AddTag(dir string, tag string, message string) error {
	_, err := gitCommand(
		dir,
		"error adding tag",
		"tag", "-a", tag, "-m", message,
	)

	return err
}

// LatestTag returns the latest git tag on the current branch.
func LatestTag(dir string) (string, error) {
	allTags, err := VersionTags(dir)
	if err != nil {
		return "", err
	}

	if len(allTags) == 0 {
		//nolint:err113
		return "", errors.New("no git tags found")
	}

	return allTags[len(allTags)-1], nil
}

// VersionTags returns all tags that match the semantic version syntax.
func VersionTags(dir string) ([]string, error) {
	all, err := gitCommand(
		dir,
		"error getting version tags",
		"--no-pager", "tag", "--list", "*.*.*",
	)
	if err != nil {
		return []string{}, err
	}

	return strings.Split(all, "\n"), nil
}
