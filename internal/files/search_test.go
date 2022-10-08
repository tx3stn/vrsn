package files_test

import (
	"path/filepath"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tx3stn/vrsn/internal/files"
)

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
		"ReturnsSupportVersionFilesWhenFound": {
			directory:   "testdata/all",
			assertError: require.NoError,
			expectedFiles: []string{
				"Cargo.toml",
				"package.json",
				"pyproject.toml",
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
			reflect.DeepEqual(tc.expectedFiles, actual)
		})
	}
}
