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

	for name, testCase := range testCases {
		tc := testCase

		originalFile, err := os.ReadFile(filepath.Join("testdata", tc.parentDir, tc.inputFile))
		require.NoError(t, err)

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			dir := filepath.Join("testdata", tc.parentDir)
			err := files.WriteVersionToFile(dir, tc.inputFile, tc.newVersion)
			require.ErrorIs(t, err, tc.expectedError)

			// Only assert the written contents if the writer func does not error.
			if err != nil {
				return
			}

			actual, err := files.GetVersionFromFile(dir, tc.inputFile)
			require.NoError(t, err)

			assert.Equal(t, tc.newVersion, actual)
		})

		// Revert the bumped file back to what we expect it to be.
		t.Cleanup(func() {
			if tc.newVersion == "" {
				return
			}

			err = os.WriteFile(
				filepath.Join("testdata", tc.parentDir, tc.inputFile),
				originalFile,
				0o600,
			)
			require.NoError(t, err)
		})
	}
}
