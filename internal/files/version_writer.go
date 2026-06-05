package files

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// WriteVersionToFile updates the version file with the provided new version
// value.
// The new contents are written to a temp file which then replaces the
// original, so a failure part way through never leaves a half written
// version file behind.
func WriteVersionToFile(dir string, inputFile string, newVersion string) error {
	path := filepath.Clean(versionFilePath(dir, inputFile))

	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}

	info, statErr := file.Stat()

	matcher := getVersionMatcher(inputFile)
	newContents, updateErr := matcher.updateVersionInPlace(newScanner(file), newVersion)

	// The whole file has been read so close it before any error handling,
	// it also has to be closed before the rename below works on Windows.
	if err := file.Close(); err != nil {
		return fmt.Errorf("error closing file %s: %w", inputFile, err)
	}

	if statErr != nil {
		return fmt.Errorf("error reading file info for %s: %w", inputFile, statErr)
	}

	if updateErr != nil {
		return updateErr
	}

	// The temp file is created next to the version file so the rename below
	// can't cross filesystems.
	tmpFile, err := os.CreateTemp(filepath.Dir(path), "vrsn-tmp-*")
	if err != nil {
		return fmt.Errorf("error creating temp file: %w", err)
	}

	if err := writeTempFile(tmpFile, newContents, info.Mode()); err != nil {
		// Best effort cleanup, the write error is more useful than any
		// remove error.
		_ = os.Remove(tmpFile.Name())

		return err
	}

	// #nosec G703 -- intentional: this CLI allows user-directed file paths.
	if err := os.Rename(tmpFile.Name(), path); err != nil {
		_ = os.Remove(tmpFile.Name())

		return fmt.Errorf("error renaming temp file: %w", err)
	}

	return nil
}

// writeTempFile writes the lines to the temp file and applies the original
// version file's permissions so they are preserved by the rename.
func writeTempFile(tmpFile *os.File, lines []string, mode fs.FileMode) error {
	for _, line := range lines {
		if _, err := tmpFile.WriteString(line + "\n"); err != nil {
			_ = tmpFile.Close()

			return fmt.Errorf("error writing string to file: %w", err)
		}
	}

	if err := tmpFile.Chmod(mode); err != nil {
		_ = tmpFile.Close()

		return fmt.Errorf("error setting temp file permissions: %w", err)
	}

	if err := tmpFile.Close(); err != nil {
		return fmt.Errorf("error closing temp file: %w", err)
	}

	return nil
}
