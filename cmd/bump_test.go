package cmd_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tx3stn/vrsn/cmd"
)

func TestValidateBumpOpts(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		gitTag        bool
		versionFiles  []string
		commit        bool
		expectedError error
	}{
		"ReturnsNoErrorWhenNoOptionsSupplied": {
			gitTag:        false,
			versionFiles:  nil,
			commit:        false,
			expectedError: nil,
		},
		"ReturnsNoErrorForGitTagOnly": {
			gitTag:        true,
			versionFiles:  nil,
			commit:        false,
			expectedError: nil,
		},
		"ReturnsNoErrorForGitTagWithCommit": {
			gitTag:        true,
			versionFiles:  nil,
			commit:        true,
			expectedError: nil,
		},
		"ReturnsNoErrorForFileOnly": {
			gitTag:        false,
			versionFiles:  []string{"VERSION"},
			commit:        false,
			expectedError: nil,
		},
		"ReturnsNoErrorForGitTagWithFileAndCommit": {
			gitTag:        true,
			versionFiles:  []string{"VERSION"},
			commit:        true,
			expectedError: nil,
		},
		"ReturnsNoErrorForGitTagWithMultipleFilesAndCommit": {
			gitTag:        true,
			versionFiles:  []string{"VERSION", "package.json"},
			commit:        true,
			expectedError: nil,
		},
		"ReturnsErrorForGitTagWithFileButNoCommit": {
			gitTag:        true,
			versionFiles:  []string{"VERSION"},
			commit:        false,
			expectedError: cmd.ErrGitTagFileNoCommit,
		},
		"ReturnsErrorForGitTagWithMultipleFilesButNoCommit": {
			gitTag:        true,
			versionFiles:  []string{"VERSION", "package.json"},
			commit:        false,
			expectedError: cmd.ErrGitTagFileNoCommit,
		},
	}

	for name, testCase := range testCases {
		tc := testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := cmd.ValidateBumpOpts(tc.gitTag, tc.versionFiles, tc.commit)
			require.ErrorIs(t, err, tc.expectedError)
		})
	}
}
