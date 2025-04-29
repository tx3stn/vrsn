package files_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tx3stn/vrsn/internal/files"
)

func TestIsGitDir(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		inputDir      string
		errorExpected require.ErrorAssertionFunc
		expected      bool
	}{
		"ReturnsTrueIfIsGitDir": {
			inputDir:      "testdata/all",
			errorExpected: require.NoError,
			expected:      true,
		},
		"ReturnsFalseIfNotGitDir": {
			inputDir:      "testdata/no-version",
			errorExpected: require.NoError,
			expected:      false,
		},
		"ReturnsErrorIfDirectoryDoesNotExist": {
			inputDir:      "testdata/foo",
			errorExpected: require.Error,
			expected:      false,
		},
	}

	for name, testCase := range testCases {
		tc := testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if tc.inputDir == "testdata/all" {
				renameDir("gitdir", ".git")
			}

			actual, err := files.IsGitDir(tc.inputDir)
			tc.errorExpected(t, err)
			assert.Equal(t, tc.expected, actual)
		})
	}

	t.Cleanup(func() {
		renameDir(".git", "gitdir")
	})
}

func renameDir(from, to string) {
	// Git won't let you commit the `.git` directory but that's needed for this
	// test, so just rename the directory before the test runs.
	_ = os.Rename("testdata/all/"+from, "testdata/all/"+to)
}
