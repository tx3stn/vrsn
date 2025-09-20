package files_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tx3stn/vrsn/internal/files"
)

func TestWriteVersionToFile(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		parentDir     string
		inputFile     string
		newVersion    string
		expectedError error
	}{
		"ReturnsErrorForUnsupportedVersionFile": {
			parentDir:     "all",
			inputFile:     "foo.txt",
			newVersion:    "",
			expectedError: files.ErrUnsuportedFile,
		},
		"WritesVersionToBuildGradle": {
			parentDir:     "all",
			inputFile:     "build.gradle",
			newVersion:    "1.3.0",
			expectedError: nil,
		},
		"ReturnsErrorForInvalidBuildGradle": {
			parentDir:     "no-version",
			inputFile:     "build.gradle",
			newVersion:    "",
			expectedError: files.ErrGettingVersionFromBuildGradle,
		},
		"WritesVersionToBuildGradleKTS": {
			parentDir:     "all",
			inputFile:     "build.gradle.kts",
			newVersion:    "0.9.12",
			expectedError: nil,
		},
		"ReturnsErrorForInvalidBuildGradleKTS": {
			parentDir:     "no-version",
			inputFile:     "build.gradle.kts",
			newVersion:    "",
			expectedError: files.ErrGettingVersionFromBuildGradle,
		},
		"WritesVersionToCargoTOML": {
			parentDir:     "all",
			inputFile:     "Cargo.toml",
			newVersion:    "2.14.741",
			expectedError: nil,
		},
		"ReturnsErrorForInvalidCargoTOML": {
			parentDir:     "no-version",
			inputFile:     "Cargo.toml",
			newVersion:    "",
			expectedError: files.ErrGettingVersionFromTOML,
		},
		"WritesVersionToCMakeLists": {
			parentDir:     "all",
			inputFile:     "CMakeLists.txt",
			newVersion:    "1.3.0",
			expectedError: nil,
		},
		"ReturnsErrorForInvalidCMakeLists": {
			parentDir:     "no-version",
			inputFile:     "CMakeLists.txt",
			newVersion:    "",
			expectedError: files.ErrGettingVersionFromCMakeLists,
		},
		"WritesVersionToPackageJSON": {
			parentDir:     "all",
			inputFile:     "package.json",
			newVersion:    "1.0.4",
			expectedError: nil,
		},
		"ReturnsErrorForInvalidPackageJSON": {
			parentDir:     "no-version",
			inputFile:     "package.json",
			newVersion:    "",
			expectedError: files.ErrGettingVersionFromPackageJSON,
		},
		"WritesVersionToPyProjectTOML": {
			parentDir:     "all",
			inputFile:     "pyproject.toml",
			newVersion:    "9.8.123456",
			expectedError: nil,
		},
		"ReturnsErrorForInvalidPyProjectTOML": {
			parentDir:     "no-version",
			inputFile:     "pyproject.toml",
			newVersion:    "",
			expectedError: files.ErrGettingVersionFromTOML,
		},
		"WritesVersionToVERSIONFile": {
			parentDir:     "all",
			inputFile:     "VERSION",
			newVersion:    "6.6.6",
			expectedError: nil,
		},
		"WritesPrefixedVersionToBuildGradle": {
			parentDir:     "prefixed",
			inputFile:     "build.gradle",
			newVersion:    "v1.3.0",
			expectedError: nil,
		},
		"WritesPrefixedVersionToBuildGradleKTS": {
			parentDir:     "prefixed",
			inputFile:     "build.gradle.kts",
			newVersion:    "v0.9.12",
			expectedError: nil,
		},
		"WritesPrefixedVersionToCargoTOML": {
			parentDir:     "prefixed",
			inputFile:     "Cargo.toml",
			newVersion:    "v2.14.741",
			expectedError: nil,
		},
		"WritesPrefixedVersionToCMakeLists": {
			parentDir:     "prefixed",
			inputFile:     "CMakeLists.txt",
			newVersion:    "v1.3.0",
			expectedError: nil,
		},
		"WritesPrefixedVersionToPackageJSON": {
			parentDir:     "prefixed",
			inputFile:     "package.json",
			newVersion:    "v1.0.4",
			expectedError: nil,
		},
		"WritesPrefixedVersionToPyProjectTOML": {
			parentDir:     "prefixed",
			inputFile:     "pyproject.toml",
			newVersion:    "v9.8.123456",
			expectedError: nil,
		},
		"WritesPrefixedVersionToVERSIONFile": {
			parentDir:     "prefixed",
			inputFile:     "VERSION",
			newVersion:    "v6.6.6",
			expectedError: nil,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			tmpDir := copyTestFile(t, tc.parentDir, tc.inputFile)
			err := files.WriteVersionToFile(tmpDir, tc.inputFile, tc.newVersion)
			require.ErrorIs(t, err, tc.expectedError)

			if err != nil {
				return
			}

			actual, err := files.GetVersionFromFile(tmpDir, tc.inputFile)
			require.NoError(t, err)

			assert.Equal(t, tc.newVersion, actual)
		})
	}
}

func copyTestFile(t *testing.T, parentDir, filename string) string {
	t.Helper()

	tmpDir := t.TempDir()
	originalPath := filepath.Join("testdata", parentDir, filename)

	data, err := os.ReadFile(filepath.Clean(originalPath))
	require.NoError(t, err)

	testPath := filepath.Join(tmpDir, filename)
	err = os.WriteFile(testPath, data, 0o600)
	require.NoError(t, err)

	return tmpDir
}
