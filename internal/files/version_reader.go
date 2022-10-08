package files

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// GetVersionFromFile reads the version file and returns the semantic
// version contained.
func GetVersionFromFile(dir string, inputFile string) (string, error) {
	matcher, err := getVersionMatcher(inputFile)
	if err != nil {
		return "", err
	}

	file, err := os.Open(filepath.Clean(filepath.Join(dir, inputFile)))
	if err != nil {
		return "", err
	}

	defer func() {
		if err := file.Close(); err != nil {
			log.Fatalf("error closing file: %s", inputFile)
		}
	}()

	scanner := bufio.NewScanner(file)

	version, err := matcher.getVersion(scanner)
	if err != nil {
		return "", err
	}
	return version, nil
}

// GetVersionFromString handles extracting the version from an file that has
// already been read and is passed as a string such as when getting the
// contents of a file from a git branch.
func GetVersionFromString(fileName string, input string) (string, error) {
	matcher, err := getVersionMatcher(fileName)
	if err != nil {
		return "", err
	}

	reader := strings.NewReader(input)
	scanner := bufio.NewScanner(reader)

	version, err := matcher.getVersion(scanner)
	if err != nil {
		return "", err
	}
	return version, nil
}
