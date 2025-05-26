package files_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tx3stn/vrsn/internal/files"
)

func TestGetVersionFromFile(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		parentDir     string
		inputFile     string
		expectedError error
		expected      string
	}{
		"ReturnsErrorForUnsupportedVersionFile": {
			parentDir:     "all",
			inputFile:     "foo.txt",
			expectedError: files.ErrUnsuportedFile,
			expected:      "",
		},
		"ReturnsVersionFromBuildGradle": {
			parentDir:     "all",
			inputFile:     "build.gradle",
			expectedError: nil,
			expected:      "1.3.0",
		},
		"ReturnsErrorFromInvalidBuildGradle": {
			parentDir:     "no-version",
			inputFile:     "build.gradle",
			expectedError: files.ErrGettingVersionFromBuildGradle,
			expected:      "",
		},
		"ReturnsVersionFromBuildGradleKTS": {
			parentDir:     "all",
			inputFile:     "build.gradle.kts",
			expectedError: nil,
			expected:      "0.9.12",
		},
		"ReturnsErrorFromInvalidBuildGradleKTS": {
			parentDir:     "no-version",
			inputFile:     "build.gradle.kts",
			expectedError: files.ErrGettingVersionFromBuildGradle,
			expected:      "",
		},
		"ReturnsVersionFromCargoTOML": {
			parentDir:     "all",
			inputFile:     "Cargo.toml",
			expectedError: nil,
			expected:      "2.14.741",
		},
		"ReturnsErrorFromInvalidCargoTOML": {
			parentDir:     "no-version",
			inputFile:     "Cargo.toml",
			expectedError: files.ErrGettingVersionFromTOML,
			expected:      "",
		},
		"ReturnsVersionFromCMakeLists": {
			parentDir:     "all",
			inputFile:     "CMakeLists.txt",
			expectedError: nil,
			expected:      "1.3.0",
		},
		"ReturnsErrorFromInvalidCMakeLists": {
			parentDir:     "no-version",
			inputFile:     "CMakeLists.txt",
			expectedError: files.ErrGettingVersionFromCMakeLists,
			expected:      "",
		},
		"ReturnsVersionFromPackageJSON": {
			parentDir:     "all",
			inputFile:     "package.json",
			expectedError: nil,
			expected:      "1.0.4",
		},
		"ReturnsErrorFromInvalidPackageJSON": {
			parentDir:     "no-version",
			inputFile:     "package.json",
			expectedError: files.ErrGettingVersionFromPackageJSON,
			expected:      "",
		},
		"ReturnsVersionFromPyprojectTOML": {
			parentDir:     "all",
			inputFile:     "pyproject.toml",
			expectedError: nil,
			expected:      "9.8.123456",
		},
		"ReturnsErrorFromInvalidPyprojectTOML": {
			parentDir:     "no-version",
			inputFile:     "pyproject.toml",
			expectedError: files.ErrGettingVersionFromTOML,
			expected:      "",
		},
		"ReturnsVersionFromSetupPy": {
			parentDir:     "all",
			inputFile:     "setup.py",
			expectedError: nil,
			expected:      "0.2.0",
		},
		"ReturnsErrorFromInvalidSetupPy": {
			parentDir:     "no-version",
			inputFile:     "setup.py",
			expectedError: files.ErrGettingVersionFromSetupPy,
			expected:      "",
		},
		"ReturnsVersionFromVERSIONFile": {
			parentDir:     "all",
			inputFile:     "VERSION",
			expectedError: nil,
			expected:      "6.6.6",
		},
		"ReturnsErrorFromInvalidVERSIONFile": {
			parentDir:     "no-version",
			inputFile:     "VERSION",
			expectedError: files.ErrGettingVersionFromVERSION,
			expected:      "",
		},
	}

	for name, testCase := range testCases {
		tc := testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			dir := filepath.Join("testdata", tc.parentDir)
			actual, err := files.GetVersionFromFile(dir, tc.inputFile)

			require.ErrorIs(t, err, tc.expectedError)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestGetVersionFromString(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		parentDir     string
		inputFile     string
		expectedError error
		expected      string
	}{
		"ReturnsVersionFromBuildGradle": {
			parentDir:     "all",
			inputFile:     "build.gradle",
			expectedError: nil,
			expected:      "1.3.0",
		},
		"ReturnsErrorFromInvalidBuildGradle": {
			parentDir:     "no-version",
			inputFile:     "build.gradle",
			expectedError: files.ErrGettingVersionFromBuildGradle,
			expected:      "",
		},
		"ReturnsVersionFromBuildGradleKTS": {
			parentDir:     "all",
			inputFile:     "build.gradle.kts",
			expectedError: nil,
			expected:      "0.9.12",
		},
		"ReturnsErrorFromInvalidBuildGradleKTS": {
			parentDir:     "no-version",
			inputFile:     "build.gradle.kts",
			expectedError: files.ErrGettingVersionFromBuildGradle,
			expected:      "",
		},
		"ReturnsVersionFromCargoTOML": {
			parentDir:     "all",
			inputFile:     "Cargo.toml",
			expectedError: nil,
			expected:      "2.14.741",
		},
		"ReturnsErrorFromInvalidCargoTOML": {
			parentDir:     "no-version",
			inputFile:     "Cargo.toml",
			expectedError: files.ErrGettingVersionFromTOML,
			expected:      "",
		},
		"ReturnsVersionFromPackageJSON": {
			parentDir:     "all",
			inputFile:     "package.json",
			expectedError: nil,
			expected:      "1.0.4",
		},
		"ReturnsErrorFromInvalidPackageJSON": {
			parentDir:     "no-version",
			inputFile:     "package.json",
			expectedError: files.ErrGettingVersionFromPackageJSON,
			expected:      "",
		},
		"ReturnsVersionFromPyprojectTOML": {
			parentDir:     "all",
			inputFile:     "pyproject.toml",
			expectedError: nil,
			expected:      "9.8.123456",
		},
		"ReturnsErrorFromInvalidPyprojectTOML": {
			parentDir:     "no-version",
			inputFile:     "pyproject.toml",
			expectedError: files.ErrGettingVersionFromTOML,
			expected:      "",
		},
		"ReturnsVersionFromSetupPy": {
			parentDir:     "all",
			inputFile:     "setup.py",
			expectedError: nil,
			expected:      "0.2.0",
		},
		"ReturnsErrorFromInvalidSetupPy": {
			parentDir:     "no-version",
			inputFile:     "setup.py",
			expectedError: files.ErrGettingVersionFromSetupPy,
			expected:      "",
		},
		"ReturnsVersionFromVERSIONFile": {
			parentDir:     "all",
			inputFile:     "VERSION",
			expectedError: nil,
			expected:      "6.6.6",
		},
		"ReturnsErrorFromInvalidVERSIONFile": {
			parentDir:     "no-version",
			inputFile:     "VERSION",
			expectedError: files.ErrGettingVersionFromVERSION,
			expected:      "",
		},
		"ReturnsVersionWithPrefixFromBuildGradle": {
			parentDir:     "prefixed",
			inputFile:     "build.gradle",
			expectedError: nil,
			expected:      "v1.3.0",
		},
		"ReturnsVersionWithPrefixFromBuildGradleKTS": {
			parentDir:     "prefixed",
			inputFile:     "build.gradle.kts",
			expectedError: nil,
			expected:      "v0.9.12",
		},
		"ReturnsVersionWithPrefixFromCargoTOML": {
			parentDir:     "prefixed",
			inputFile:     "Cargo.toml",
			expectedError: nil,
			expected:      "v2.14.741",
		},
		"ReturnsVersionWithPrefixFromCMakeLists": {
			parentDir:     "prefixed",
			inputFile:     "CMakeLists.txt",
			expectedError: nil,
			expected:      "v1.3.0",
		},
		"ReturnsVersionWithPrefixFromPackageJSON": {
			parentDir:     "prefixed",
			inputFile:     "package.json",
			expectedError: nil,
			expected:      "v1.0.4",
		},
		"ReturnsVersionWithPrefixFromPyprojectTOML": {
			parentDir:     "prefixed",
			inputFile:     "pyproject.toml",
			expectedError: nil,
			expected:      "v9.8.123456",
		},
		"ReturnsVersionWithPrefixFromSetupPy": {
			parentDir:     "prefixed",
			inputFile:     "setup.py",
			expectedError: nil,
			expected:      "v0.2.0",
		},
		"ReturnsVersionWithPrefixFromVERSIONFile": {
			parentDir:     "prefixed",
			inputFile:     "VERSION",
			expectedError: nil,
			expected:      "v6.6.5",
		},
	}

	for name, testCase := range testCases {
		tc := testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			content, err := os.ReadFile(filepath.Join("testdata", tc.parentDir, tc.inputFile))
			require.NoError(t, err)

			actual, err := files.GetVersionFromString(tc.inputFile, string(content))

			require.ErrorIs(t, err, tc.expectedError)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
