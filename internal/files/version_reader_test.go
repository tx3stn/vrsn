package files_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tx3stn/vrsn/internal/files"
	"github.com/tx3stn/vrsn/internal/test"
)

func TestGetVersionFromFile(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		parentDir   string
		inputFile   string
		assertError require.ErrorAssertionFunc
		expected    string
	}{
		"ReturnsErrorForUnsupportedVersionFile": {
			parentDir:   "all",
			inputFile:   "foo.txt",
			assertError: require.Error,
			expected:    "",
		},
		"ReturnsVersionFromBuildGradle": {
			parentDir:   "all",
			inputFile:   "build.gradle",
			assertError: require.NoError,
			expected:    "1.3.0",
		},
		"ReturnsErrorFromInvalidBuildGradle": {
			parentDir:   "no-version",
			inputFile:   "build.gradle",
			assertError: test.IsSentinelError(files.ErrGettingVersionFromBuildGradle),
			expected:    "",
		},
		"ReturnsVersionFromBuildGradleKTS": {
			parentDir:   "all",
			inputFile:   "build.gradle.kts",
			assertError: require.NoError,
			expected:    "0.9.12",
		},
		"ReturnsErrorFromInvalidBuildGradleKTS": {
			parentDir:   "no-version",
			inputFile:   "build.gradle.kts",
			assertError: test.IsSentinelError(files.ErrGettingVersionFromBuildGradle),
			expected:    "",
		},
		"ReturnsVersionFromCargoTOML": {
			parentDir:   "all",
			inputFile:   "Cargo.toml",
			assertError: require.NoError,
			expected:    "2.14.741",
		},
		"ReturnsErrorFromInvalidCargoTOML": {
			parentDir:   "no-version",
			inputFile:   "Cargo.toml",
			assertError: test.IsSentinelError(files.ErrGettingVersionFromTOML),
			expected:    "",
		},
		"ReturnsVersionFromCMakeLists": {
			parentDir:   "all",
			inputFile:   "CMakeLists.txt",
			assertError: require.NoError,
			expected:    "1.3.0",
		},
		"ReturnsErrorFromInvalidCMakeLists": {
			parentDir:   "no-version",
			inputFile:   "CMakeLists.txt",
			assertError: test.IsSentinelError(files.ErrGettingVersionFromCMakeLists),
			expected:    "",
		},
		"ReturnsVersionFromPackageJSON": {
			parentDir:   "all",
			inputFile:   "package.json",
			assertError: require.NoError,
			expected:    "1.0.4",
		},
		"ReturnsErrorFromInvalidPackageJSON": {
			parentDir:   "no-version",
			inputFile:   "package.json",
			assertError: test.IsSentinelError(files.ErrGettingVersionFromPackageJSON),
			expected:    "",
		},
		"ReturnsVersionFromPyprojectTOML": {
			parentDir:   "all",
			inputFile:   "pyproject.toml",
			assertError: require.NoError,
			expected:    "9.8.123456",
		},
		"ReturnsErrorFromInvalidPyprojectTOML": {
			parentDir:   "no-version",
			inputFile:   "pyproject.toml",
			assertError: test.IsSentinelError(files.ErrGettingVersionFromTOML),
			expected:    "",
		},
		"ReturnsVersionFromSetupPy": {
			parentDir:   "all",
			inputFile:   "setup.py",
			assertError: require.NoError,
			expected:    "0.2.0",
		},
		"ReturnsErrorFromInvalidSetupPy": {
			parentDir:   "no-version",
			inputFile:   "setup.py",
			assertError: test.IsSentinelError(files.ErrGettingVersionFromSetupPy),
			expected:    "",
		},
		"ReturnsVersionFromVERSIONFile": {
			parentDir:   "all",
			inputFile:   "VERSION",
			assertError: require.NoError,
			expected:    "6.6.6",
		},
		"ReturnsErrorFromInvalidVERSIONFile": {
			parentDir:   "no-version",
			inputFile:   "VERSION",
			assertError: test.IsSentinelError(files.ErrGettingVersionFromVERSION),
			expected:    "",
		},
	}

	for name, testCase := range testCases {
		tc := testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			dir := filepath.Join("testdata", tc.parentDir)
			actual, err := files.GetVersionFromFile(dir, tc.inputFile)
			tc.assertError(t, err)

			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestGetVersionFromString(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		parentDir   string
		inputFile   string
		assertError require.ErrorAssertionFunc
		expected    string
	}{
		"ReturnsVersionFromBuildGradle": {
			parentDir:   "all",
			inputFile:   "build.gradle",
			assertError: require.NoError,
			expected:    "1.3.0",
		},
		"ReturnsErrorFromInvalidBuildGradle": {
			parentDir:   "no-version",
			inputFile:   "build.gradle",
			assertError: test.IsSentinelError(files.ErrGettingVersionFromBuildGradle),
			expected:    "",
		},
		"ReturnsVersionFromBuildGradleKTS": {
			parentDir:   "all",
			inputFile:   "build.gradle.kts",
			assertError: require.NoError,
			expected:    "0.9.12",
		},
		"ReturnsErrorFromInvalidBuildGradleKTS": {
			parentDir:   "no-version",
			inputFile:   "build.gradle.kts",
			assertError: test.IsSentinelError(files.ErrGettingVersionFromBuildGradle),
			expected:    "",
		},
		"ReturnsVersionFromCargoTOML": {
			parentDir:   "all",
			inputFile:   "Cargo.toml",
			assertError: require.NoError,
			expected:    "2.14.741",
		},
		"ReturnsErrorFromInvalidCargoTOML": {
			parentDir:   "no-version",
			inputFile:   "Cargo.toml",
			assertError: test.IsSentinelError(files.ErrGettingVersionFromTOML),
			expected:    "",
		},
		"ReturnsVersionFromPackageJSON": {
			parentDir:   "all",
			inputFile:   "package.json",
			assertError: require.NoError,
			expected:    "1.0.4",
		},
		"ReturnsErrorFromInvalidPackageJSON": {
			parentDir:   "no-version",
			inputFile:   "package.json",
			assertError: test.IsSentinelError(files.ErrGettingVersionFromPackageJSON),
			expected:    "",
		},
		"ReturnsVersionFromPyprojectTOML": {
			parentDir:   "all",
			inputFile:   "pyproject.toml",
			assertError: require.NoError,
			expected:    "9.8.123456",
		},
		"ReturnsErrorFromInvalidPyprojectTOML": {
			parentDir:   "no-version",
			inputFile:   "pyproject.toml",
			assertError: test.IsSentinelError(files.ErrGettingVersionFromTOML),
			expected:    "",
		},
		"ReturnsVersionFromSetupPy": {
			parentDir:   "all",
			inputFile:   "setup.py",
			assertError: require.NoError,
			expected:    "0.2.0",
		},
		"ReturnsErrorFromInvalidSetupPy": {
			parentDir:   "no-version",
			inputFile:   "setup.py",
			assertError: test.IsSentinelError(files.ErrGettingVersionFromSetupPy),
			expected:    "",
		},
		"ReturnsVersionFromVERSIONFile": {
			parentDir:   "all",
			inputFile:   "VERSION",
			assertError: require.NoError,
			expected:    "6.6.6",
		},
		"ReturnsErrorFromInvalidVERSIONFile": {
			parentDir:   "no-version",
			inputFile:   "VERSION",
			assertError: test.IsSentinelError(files.ErrGettingVersionFromVERSION),
			expected:    "",
		},
	}

	for name, testCase := range testCases {
		tc := testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			content, err := os.ReadFile(filepath.Join("testdata", tc.parentDir, tc.inputFile))
			require.NoError(t, err)

			actual, err := files.GetVersionFromString(tc.inputFile, string(content))
			tc.assertError(t, err)

			assert.Equal(t, tc.expected, actual)
		})
	}
}
