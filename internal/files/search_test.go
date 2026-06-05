package files_test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tx3stn/vrsn/internal/files"
	"github.com/tx3stn/vrsn/internal/logger"
)

func TestFind(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		fileFlag            string
		searchDir           string
		errorOnNoFilesFound bool
		expected            string
		expectedError       error
	}{
		"ReturnsFileFlagWhenFileExists": {
			fileFlag:      "testdata/single/VERSION",
			expected:      "testdata/single/VERSION",
			expectedError: nil,
		},
		"ReturnsFileFlagForUnsupportedFileType": {
			fileFlag:      "testdata/all/foo.txt",
			expected:      "testdata/all/foo.txt",
			expectedError: nil,
		},
		"ReturnsErrorWhenFileFlagDoesNotExist": {
			fileFlag:      "testdata/single/nope",
			expected:      "",
			expectedError: files.ErrFileNotFound,
		},
		"ReturnsErrorWhenFileFlagIsDirectory": {
			fileFlag:      "testdata/single",
			expected:      "",
			expectedError: files.ErrFileIsDirectory,
		},
		"ReturnsSingleVersionFileInSearchDir": {
			searchDir:     "testdata/single",
			expected:      "VERSION",
			expectedError: nil,
		},
		"ReturnsErrorWhenNoVersionFilesFoundAndErrorOnNoFilesFound": {
			searchDir:           "testdata/empty",
			errorOnNoFilesFound: true,
			expected:            "",
			expectedError:       files.ErrNoVersionFilesInDir,
		},
		"ReturnsNoFileAndNoErrorWhenNoVersionFilesFound": {
			searchDir:           "testdata/empty",
			errorOnNoFilesFound: false,
			expected:            "",
			expectedError:       nil,
		},
		"ReturnsErrorWhenMultipleVersionFilesFound": {
			searchDir:     "testdata/all",
			expected:      "",
			expectedError: files.ErrMultipleVersionFiles,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			finder := files.VersionFileFinder{
				ErrorOnNoFilesFound: tc.errorOnNoFilesFound,
				FileFlag:            filepath.FromSlash(tc.fileFlag),
				Logger:              logger.NewBasic(false, false),
				SearchDir:           filepath.FromSlash(tc.searchDir),
			}

			actual, err := finder.Find()
			require.ErrorIs(t, err, tc.expectedError)
			assert.Equal(t, filepath.FromSlash(tc.expected), actual)
		})
	}
}

func TestGetVersionFilesInDirectory(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		directory     string
		assertError   require.ErrorAssertionFunc
		expectedFiles []string
	}{
		"ReturnsEmptySliceWhenNoVersionFilesFound": {
			directory:     "testdata/empty",
			assertError:   require.NoError,
			expectedFiles: []string{},
		},
		"ReturnsSupportedVersionFilesWhenFound": {
			directory:   "testdata/all",
			assertError: require.NoError,
			expectedFiles: []string{
				"BUILD.bazel",
				"build.gradle",
				"build.gradle.kts",
				"Cargo.toml",
				"CMakeLists.txt",
				"package.json",
				"pyproject.toml",
				"setup.py",
				"VERSION",
			},
		},
		"ReturnsErrorWhenDirectoryDoesNotExist": {
			directory:     "testdata/foo",
			assertError:   require.Error,
			expectedFiles: []string{},
		},
	}

	for name, testCase := range testCases {
		tc := testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			path := filepath.FromSlash(tc.directory)
			actual, err := files.GetVersionFilesInDirectory(path)
			tc.assertError(t, err)
			assert.ElementsMatch(t, tc.expectedFiles, actual)
		})
	}
}
