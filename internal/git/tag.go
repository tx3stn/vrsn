package git

import (
	"strings"

	"github.com/tx3stn/vrsn/internal/version"
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

// LatestTag returns the latest semantic version tag in the repository.
func LatestTag(dir string) (string, error) {
	allTags, err := VersionTags(dir)
	if err != nil {
		return "", err
	}

	if len(allTags) == 0 {
		return "", ErrNoGitTags
	}

	return allTags[len(allTags)-1], nil
}

// VersionTags returns all tags that match the semantic version syntax, sorted
// by version so the latest version is last rather than git's default
// lexicographic order (which sorts 0.0.9 after 0.0.10).
// Tags matching the glob but not parseable as a semantic version (e.g.
// 1.2.3-rc1) are filtered out so bumping is always based on a valid version.
func VersionTags(dir string) ([]string, error) {
	all, err := gitCommand(
		dir,
		"error getting version tags",
		"--no-pager", "tag", "--list", "--sort=v:refname", "*.*.*",
	)
	if err != nil {
		return []string{}, err
	}

	if all == "" {
		return []string{}, nil
	}

	versionTags := []string{}

	for tag := range strings.SplitSeq(all, "\n") {
		if _, err := version.Parse(tag); err == nil {
			versionTags = append(versionTags, tag)
		}
	}

	return versionTags, nil
}
