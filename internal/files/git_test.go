package files_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tx3stn/vrsn/internal/files"
)

func TestIsGitDir(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		inputDir      string
		errorExpected error
		expected      bool
		needsGitDir   bool
	}{
		"ReturnsTrueIfIsGitDir": {
			inputDir:      "testdata/all",
			errorExpected: nil,
			expected:      true,
			needsGitDir:   true,
		},
		"ReturnsFalseIfNotGitDir": {
			inputDir:      "testdata/no-version",
			errorExpected: nil,
			expected:      false,
			needsGitDir:   false,
		},
		"ReturnsErrorIfDirectoryDoesNotExist": {
			inputDir:      "testdata/foo",
			errorExpected: files.ErrGettingFilesInDirectory,
			expected:      false,
			needsGitDir:   false,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			testDir := tc.inputDir

			// Copy testdata to temp dir and create .git directory
			if tc.needsGitDir {
				tmpDir := t.TempDir()
				err := os.MkdirAll(filepath.Join(tmpDir, ".git"), 0o750)
				require.NoError(t, err)

				testDir = tmpDir
			}

			actual, err := files.IsGitDir(testDir)
			require.ErrorIs(t, err, tc.errorExpected)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
