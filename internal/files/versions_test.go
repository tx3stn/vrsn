package files_test

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tx3stn/vrsn/internal/files"
	"github.com/tx3stn/vrsn/internal/logger"
)

func TestGetVersionsFromFiles(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		dirFiles      map[string]string
		versionFiles  []string
		expected      string
		expectedError error
	}{
		"ReturnsVersionFromSingleFile": {
			dirFiles: map[string]string{
				"VERSION": "0.1.0\n",
			},
			versionFiles:  []string{"VERSION"},
			expected:      "0.1.0",
			expectedError: nil,
		},
		"ReturnsCommonVersionWhenAllFilesMatch": {
			dirFiles: map[string]string{
				"VERSION":      "0.1.0\n",
				"package.json": `{"version": "0.1.0"}`,
				"Cargo.toml":   "[package]\nversion = \"0.1.0\"\n",
			},
			versionFiles:  []string{"VERSION", "package.json", "Cargo.toml"},
			expected:      "0.1.0",
			expectedError: nil,
		},
		"ReturnsErrorWhenVersionsDoNotMatch": {
			dirFiles: map[string]string{
				"VERSION":      "0.1.0\n",
				"package.json": `{"version": "0.2.0"}`,
			},
			versionFiles:  []string{"VERSION", "package.json"},
			expected:      "",
			expectedError: files.ErrVersionsDoNotMatch,
		},
		"ReturnsErrorWhenFileDoesNotExist": {
			dirFiles: map[string]string{
				"VERSION": "0.1.0\n",
			},
			versionFiles:  []string{"VERSION", "package.json"},
			expected:      "",
			expectedError: fs.ErrNotExist,
		},
		"ReturnsErrorWhenNoFilesProvided": {
			dirFiles:      map[string]string{},
			versionFiles:  []string{},
			expected:      "",
			expectedError: files.ErrNoVersionFilesInDir,
		},
	}

	for name, testCase := range testCases {
		tc := testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			dir := t.TempDir()
			for file, content := range tc.dirFiles {
				err := os.WriteFile(filepath.Join(dir, file), []byte(content), 0o600)
				require.NoError(t, err)
			}

			log := logger.NewBasic(false, false)

			version, err := files.GetVersionsFromFiles(dir, tc.versionFiles, log)
			require.ErrorIs(t, err, tc.expectedError)
			assert.Equal(t, tc.expected, version)
		})
	}
}
