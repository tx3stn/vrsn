package files

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// GetVersionFromFile reads the version file and returns the semantic
// version contained.
func GetVersionFromFile(dir string, inputFile string) (string, error) {
	file, err := os.Open(filepath.Clean(versionFilePath(dir, inputFile)))
	if err != nil {
		return "", fmt.Errorf("error opening version file: %w", err)
	}

	// The file is only read so a close error can't affect the result.
	defer func() {
		_ = file.Close()
	}()

	return getVersionFromReader(inputFile, file)
}

// GetVersionFromString handles extracting the version from a file that has
// already been read and is passed as a string such as when getting the
// contents of a file from a git branch.
func GetVersionFromString(fileName string, input string) (string, error) {
	return getVersionFromReader(fileName, strings.NewReader(input))
}

// getVersionFromReader extracts the version from the reader using the
// matcher config for the provided file name.
func getVersionFromReader(fileName string, reader io.Reader) (string, error) {
	matcher := getVersionMatcher(fileName)

	return matcher.getVersion(newScanner(reader))
}

// versionFilePath resolves the path to the version file, supporting absolute
// paths provided with the --file flag.
func versionFilePath(dir string, inputFile string) string {
	if filepath.IsAbs(inputFile) {
		return inputFile
	}

	return filepath.Join(dir, inputFile)
}
