package files

import (
	"fmt"

	"github.com/tx3stn/vrsn/internal/logger"
)

// GetVersionsFromFiles reads the version from each of the provided files and
// returns the common version they all contain.
// The version found in each file is debug logged, and if the versions do not
// all match an ErrVersionsDoNotMatch error is returned.
func GetVersionsFromFiles(dir string, versionFiles []string, log logger.Basic) (string, error) {
	if len(versionFiles) == 0 {
		return "", ErrNoVersionFilesInDir
	}

	versions := make([]string, 0, len(versionFiles))

	for _, file := range versionFiles {
		version, err := GetVersionFromFile(dir, file)
		if err != nil {
			return "", fmt.Errorf("error getting version from file %s: %w", file, err)
		}

		log.Debugf("file %s has version %s", file, version)

		versions = append(versions, version)
	}

	return CommonVersion(versions)
}

// CommonVersion returns the version shared by all of the provided versions,
// or an ErrVersionsDoNotMatch error when they aren't all the same.
func CommonVersion(versions []string) (string, error) {
	for _, version := range versions {
		if version != versions[0] {
			return "", ErrVersionsDoNotMatch
		}
	}

	return versions[0], nil
}
