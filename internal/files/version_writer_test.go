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
		"ReturnsErrorWhenBestEffortMatchingFindsNoVersion": {
			parentDir:     "all",
			inputFile:     "foo.txt",
			newVersion:    "",
			expectedError: files.ErrGettingVersionBestEffort,
		},
		"WritesVersionToSingleQuotedFileWithBestEffort": {
			parentDir:     "all",
			inputFile:     "version.ts",
			newVersion:    "0.0.11",
			expectedError: nil,
		},
		"ReturnsErrorForInvalidBestEffortFile": {
			parentDir:     "no-version",
			inputFile:     "version.ts",
			newVersion:    "",
			expectedError: files.ErrGettingVersionBestEffort,
		},
		"WritesVersionToDoubleQuotedFileWithBestEffort": {
			parentDir:     "all",
			inputFile:     "version.js",
			newVersion:    "1.5.1",
			expectedError: nil,
		},
		"WritesVersionToUnquotedFileWithBestEffort": {
			parentDir:     "all",
			inputFile:     "app.conf",
			newVersion:    "0.3.3",
			expectedError: nil,
		},
		"WritesVersionToBUILDBazel": {
			parentDir:     "all",
			inputFile:     "BUILD.bazel",
			newVersion:    "0.16.4",
			expectedError: nil,
		},
		"ReturnsErrorForInvalidBUILDBazel": {
			parentDir:     "no-version",
			inputFile:     "BUILD.bazel",
			newVersion:    "",
			expectedError: files.ErrGettingVersionFromBuildBazel,
		},
		"WritesVersionToBazelModule": {
			parentDir:     "all",
			inputFile:     "MODULE.bazel",
			newVersion:    "0.0.24",
			expectedError: nil,
		},
		"ReturnsErrorForInvalidBazelModule": {
			parentDir:     "no-version",
			inputFile:     "MODULE.bazel",
			newVersion:    "",
			expectedError: files.ErrGettingVersionFromBazelModule,
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
		"WritesPrefixedVersionToBUILDBazel": {
			parentDir:     "prefixed",
			inputFile:     "BUILD.bazel",
			newVersion:    "v0.16.4",
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
		"WritesPrefixedVersionToBestEffortFile": {
			parentDir:     "prefixed",
			inputFile:     "version.ts",
			newVersion:    "v0.0.10",
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

// TestWriteVersionToFileOnlyUpdatesFirstMatch is the regression test for
// files where other version strings appear after the package version, such as
// dependency constraints, which must not be rewritten by a bump.
func TestWriteVersionToFileOnlyUpdatesFirstMatch(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		inputFile        string
		expectedContents string
	}{
		"OnlyUpdatesPackageVersionInCargoTOML": {
			inputFile: "Cargo.toml",
			expectedContents: `[package]
name = "with-deps"
version = "2.0.0"
authors = ["me"]
license = "GPL-3.0"

[dependencies]
tokio = { version = "1.38.2", features = ["full"] }
serde = { version = "1.0.219" }
`,
		},
		"OnlyUpdatesPackageVersionInPyProjectTOML": {
			inputFile: "pyproject.toml",
			expectedContents: `[tool.poetry]
name = "with-deps"
version = "2.0.0"
description = "testing dependency versions are not clobbered"

[tool.poetry.dependencies]
python = "3.11.4"
requests = { version = "2.32.3" }
`,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			tmpDir := copyTestFile(t, "with-deps", tc.inputFile)
			err := files.WriteVersionToFile(tmpDir, tc.inputFile, "2.0.0")
			require.NoError(t, err)

			// #nosec G304 -- reading a test fixture from the temp dir.
			actual, err := os.ReadFile(filepath.Join(tmpDir, tc.inputFile))
			require.NoError(t, err)

			assert.Equal(t, tc.expectedContents, string(actual))
		})
	}
}

// TestWriteVersionToFilePreservesPermissions checks the original file mode
// survives the temp file replacing the version file.
func TestWriteVersionToFilePreservesPermissions(t *testing.T) {
	t.Parallel()

	tmpDir := copyTestFile(t, "all", "VERSION")
	path := filepath.Join(tmpDir, "VERSION")
	// #nosec G302 -- a non default mode is the point of this test.
	require.NoError(t, os.Chmod(path, 0o644))

	err := files.WriteVersionToFile(tmpDir, "VERSION", "1.2.3")
	require.NoError(t, err)

	info, err := os.Stat(path)
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0o644), info.Mode().Perm())
}

// TestWriteVersionToFileWithAbsolutePath checks absolute file paths are used
// as is rather than being joined to the current directory.
func TestWriteVersionToFileWithAbsolutePath(t *testing.T) {
	t.Parallel()

	tmpDir := copyTestFile(t, "all", "VERSION")
	absPath := filepath.Join(tmpDir, "VERSION")

	err := files.WriteVersionToFile("/some/other/dir", absPath, "4.5.6")
	require.NoError(t, err)

	actual, err := files.GetVersionFromFile("/another/dir", absPath)
	require.NoError(t, err)
	assert.Equal(t, "4.5.6", actual)
}

func copyTestFile(t *testing.T, parentDir, filename string) string {
	t.Helper()

	tmpDir := t.TempDir()
	originalPath := filepath.Join("testdata", parentDir, filename)

	data, err := os.ReadFile(filepath.Clean(originalPath))
	require.NoError(t, err)

	testPath := filepath.Join(tmpDir, filename)
	//#nosec: G703
	err = os.WriteFile(testPath, data, 0o600)
	require.NoError(t, err)

	return tmpDir
}
