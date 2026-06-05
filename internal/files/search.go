// Package files handles logic for interacting with files.
package files

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/tx3stn/vrsn/internal/logger"
)

// VersionFileFinder looks for the relevant version file based on the options
// specified.
type VersionFileFinder struct {
	ErrorOnNoFilesFound bool
	FileFlag            string
	Logger              logger.Basic
	SearchDir           string
}

// Find returns the version file based on the config provided.
func (v VersionFileFinder) Find() (string, error) {
	if v.FileFlag != "" {
		return v.findFromFlag()
	}

	v.Logger.Debugf("looking for version files in %s", v.SearchDir)

	allVersionFiles, err := GetVersionFilesInDirectory(v.SearchDir)
	if err != nil {
		return "", err
	}

	v.Logger.Debugf("found version files: %v", allVersionFiles)

	switch len(allVersionFiles) {
	case 1:
		return allVersionFiles[0], nil

	case 0:
		if v.ErrorOnNoFilesFound {
			return "", ErrNoVersionFilesInDir
		}

		return "", nil

	default:
		return "", ErrMultipleVersionFiles
	}
}

// findFromFlag validates the version file explicitly provided with the --file
// flag exists and returns it.
func (v VersionFileFinder) findFromFlag() (string, error) {
	v.Logger.Debugf("using specified version file %s", v.FileFlag)

	info, err := os.Stat(v.FileFlag)
	// Handle not exists error first for better error output.
	if errors.Is(err, fs.ErrNotExist) {
		return "", fmt.Errorf("%w: file:%s", ErrFileNotFound, v.FileFlag)
	}

	if err != nil {
		return "", fmt.Errorf("error checking for file %s: %w", v.FileFlag, err)
	}

	if info.IsDir() {
		return "", fmt.Errorf("%w: file:%s", ErrFileIsDirectory, v.FileFlag)
	}

	if _, supported := versionFileMatchers[filepath.Base(v.FileFlag)]; !supported {
		v.Logger.Debugf(
			"%s is not a natively supported version file, will attempt best effort matching",
			v.FileFlag,
		)
	}

	return v.FileFlag, nil
}

// GetVersionFilesInDirectory checks the provided directory for supported
// version files and returns a list of ones found.
func GetVersionFilesInDirectory(dir string) ([]string, error) {
	allFiles, err := os.ReadDir(dir)
	if err != nil {
		return []string{}, fmt.Errorf("error getting version files in directory: %w", err)
	}

	versionFiles := []string{}

	for _, file := range allFiles {
		if file.IsDir() {
			continue
		}

		name := file.Name()
		if _, supported := versionFileMatchers[name]; supported {
			versionFiles = append(versionFiles, name)
		}
	}

	return versionFiles, nil
}
