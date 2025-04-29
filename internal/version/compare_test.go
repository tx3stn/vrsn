package version_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tx3stn/vrsn/internal/version"
)

func TestCompare(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		was           string
		now           string
		expectedError error
	}{
		"ReturnsVersionNotBumpedErrorWhenVersionsAreTheSame": {
			was:           "1.0.0",
			now:           "1.0.0",
			expectedError: version.ErrVersionNotBumped,
		},
		"ReturnsErrorWhenWasFailsValidation": {
			was:           "",
			now:           "1.1.1",
			expectedError: version.ErrNoVersionParts,
		},
		"ReturnsErrorWhenNowFailsValidation": {
			was:           "1.1.1",
			now:           "",
			expectedError: version.ErrNoVersionParts,
		},
		"ReturnsInvalidBumpErrorWhenNotValidSemVer": {
			was:           "1.0.0",
			now:           "1.0.3",
			expectedError: version.ErrInvalidBump,
		},
		"ReturnsNoErrorForValidPatch": {
			was:           "1.0.0",
			now:           "1.0.1",
			expectedError: nil,
		},
		"ReturnsNoErrorForValidMinor": {
			was:           "1.0.0",
			now:           "1.1.0",
			expectedError: nil,
		},
		"ReturnsNoErrorForValidMajor": {
			was:           "1.0.0",
			now:           "2.0.0",
			expectedError: nil,
		},
	}

	for name, testCase := range testCases {
		tc := testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := version.Compare(tc.was, tc.now)
			require.ErrorIs(t, err, tc.expectedError)
		})
	}
}
