package files

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// WriteVersionToFile updates the version file with the provided new version
// value.
func WriteVersionToFile(dir string, inputFile string, newVersion string) error {
	matcher, err := getVersionMatcher(inputFile)
	if err != nil {
		return err
	}

	file, err := os.Open(filepath.Clean(filepath.Join(dir, inputFile)))
	if err != nil {
		return err
	}

	defer func() {
		if err := file.Close(); err != nil {
			log.Fatalf("error closing file: %s\n%s", inputFile, err)
		}
	}()

	scanner := bufio.NewScanner(file)

	newContents, err := matcher.updateVersionInPlace(scanner, newVersion)
	if err != nil {
		return err
	}

	tmpFile, err := os.CreateTemp(dir, "vrsn-tmp-*")
	if err != nil {
		return err
	}

	defer func() {
		if err := tmpFile.Close(); err != nil {
			log.Fatalf("error closing temp file while bumping version: %s\n", err)
		}
	}()

	for _, line := range newContents {
		if _, err := tmpFile.WriteString(fmt.Sprintf("%s\n", line)); err != nil {
			return err
		}
	}

	if err := os.Rename(tmpFile.Name(), file.Name()); err != nil {
		return err
	}

	return nil
}
