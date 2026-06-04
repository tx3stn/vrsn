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
		versionFile   string
		commit        bool
		expectedError error
	}{
		"ReturnsNoErrorWhenNoOptionsSupplied": {
			gitTag:        false,
			versionFile:   "",
			commit:        false,
			expectedError: nil,
		},
		"ReturnsNoErrorForGitTagOnly": {
			gitTag:        true,
			versionFile:   "",
			commit:        false,
			expectedError: nil,
		},
		"ReturnsNoErrorForGitTagWithCommit": {
			gitTag:        true,
			versionFile:   "",
			commit:        true,
			expectedError: nil,
		},
		"ReturnsNoErrorForFileOnly": {
			gitTag:        false,
			versionFile:   "VERSION",
			commit:        false,
			expectedError: nil,
		},
		"ReturnsNoErrorForGitTagWithFileAndCommit": {
			gitTag:        true,
			versionFile:   "VERSION",
			commit:        true,
			expectedError: nil,
		},
		"ReturnsErrorForGitTagWithFileButNoCommit": {
			gitTag:        true,
			versionFile:   "VERSION",
			commit:        false,
			expectedError: cmd.ErrGitTagFileNoCommit,
		},
	}

	for name, testCase := range testCases {
		tc := testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := cmd.ValidateBumpOpts(tc.gitTag, tc.versionFile, tc.commit)
			require.ErrorIs(t, err, tc.expectedError)
		})
	}
}
