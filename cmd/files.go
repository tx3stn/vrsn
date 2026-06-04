package cmd

import (
	"fmt"

	"github.com/tx3stn/vrsn/internal/files"
	"github.com/tx3stn/vrsn/internal/logger"
)

// resolveVersionFiles returns the list of version files to operate on.
// The configured files take precedence, when none are provided it falls back
// to searching the directory for a single supported version file.
func resolveVersionFiles(
	curDir string,
	configured []string,
	log logger.Basic,
	errorOnNoFilesFound bool,
) ([]string, error) {
	finder := files.VersionFileFinder{
		ErrorOnNoFilesFound: errorOnNoFilesFound,
		Logger:              log,
		SearchDir:           curDir,
	}

	if len(configured) == 0 {
		found, err := finder.Find()
		if err != nil {
			return nil, fmt.Errorf("error finding version file: %w", err)
		}

		if found == "" {
			return nil, nil
		}

		return []string{found}, nil
	}

	configured = dedupe(configured)
	versionFiles := make([]string, 0, len(configured))

	for _, file := range configured {
		finder.FileFlag = file

		versionFile, err := finder.Find()
		if err != nil {
			return nil, fmt.Errorf("error finding version file: %w", err)
		}

		versionFiles = append(versionFiles, versionFile)
	}

	return versionFiles, nil
}

// dedupe returns the provided slice with any duplicate entries removed,
// preserving the original order.
func dedupe(input []string) []string {
	seen := make(map[string]struct{}, len(input))
	out := make([]string, 0, len(input))

	for _, entry := range input {
		if _, ok := seen[entry]; ok {
			continue
		}

		seen[entry] = struct{}{}

		out = append(out, entry)
	}

	return out
}
