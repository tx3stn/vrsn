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
			parentDir:     "bump",
			inputFile:     "foo.txt",
			newVersion:    "",
			expectedError: files.ErrUnsuportedFile,
		},
		"WritesVersionToBuildGradle": {
			parentDir:     "bump",
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
			parentDir:     "bump",
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
			parentDir:     "bump",
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
			parentDir:     "bump",
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
			parentDir:     "bump",
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
			parentDir:     "bump",
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
			parentDir:     "bump",
			inputFile:     "VERSION",
			newVersion:    "6.6.6",
			expectedError: nil,
		},
	}

	for name, testCase := range testCases {
		tc := testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			dir := filepath.Join("testdata", tc.parentDir)
			err := files.WriteVersionToFile(dir, tc.inputFile, tc.newVersion)
			require.ErrorIs(t, err, tc.expectedError)

			// Only assert the written contents if the writer func does not error.
			if err != nil {
				return
			}

			expected, err := os.ReadFile(filepath.Join("testdata", "all", tc.inputFile))
			require.NoError(t, err)

			actual, err := os.ReadFile(filepath.Clean(filepath.Join(dir, tc.inputFile)))
			require.NoError(t, err)

			assert.Equal(t, string(expected), string(actual))
		})

		// Revert the bumped file back to what we expect it to be.
		t.Cleanup(func() {
			if tc.newVersion == "" {
				return
			}

			originalValues, err := os.ReadFile(filepath.Join("testdata", "bump-og", tc.inputFile))
			require.NoError(t, err)

			err = os.WriteFile(
				filepath.Join("testdata", "bump", tc.inputFile),
				originalValues,
				0o600,
			)
			require.NoError(t, err)
		})
	}
}
