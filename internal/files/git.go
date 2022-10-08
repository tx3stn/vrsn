package files

import "os"

// IsGitDir returns if the specified directory is a git dir.
func IsGitDir(dir string) (bool, error) {
	allFiles, err := os.ReadDir(dir)
	if err != nil {
		return false, err
	}

	for _, file := range allFiles {
		if !file.IsDir() {
			continue
		}

		if file.Name() == ".git" {
			return true, nil
		}
	}

	return false, nil
}
