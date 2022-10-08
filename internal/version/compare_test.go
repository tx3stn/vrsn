package version_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tx3stn/vrsn/internal/test"
	"github.com/tx3stn/vrsn/internal/version"
)

func TestCompare(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		was         string
		now         string
		assertError require.ErrorAssertionFunc
	}{
		"ReturnsVersionNotBumpedErrorWhenVersionsAreTheSame": {
			was:         "1.0.0",
			now:         "1.0.0",
			assertError: test.IsSentinelError(version.ErrVersionNotBumped),
		},
		"ReturnsErrorWhenWasFailsValidation": {
			was:         "",
			now:         "1.1.1",
			assertError: require.Error,
		},
		"ReturnsErrorWhenNowFailsValidation": {
			was:         "1.1.1",
			now:         "",
			assertError: require.Error,
		},
		"ReturnsInvalidBumpErrorWhenNotValidSemVer": {
			was:         "1.0.0",
			now:         "1.0.3",
			assertError: test.IsSentinelError(version.ErrInvalidBump),
		},
		"ReturnsNoErrorForValidPatch": {
			was:         "1.0.0",
			now:         "1.0.1",
			assertError: require.NoError,
		},
		"ReturnsNoErrorForValidMinor": {
			was:         "1.0.0",
			now:         "1.1.0",
			assertError: require.NoError,
		},
		"ReturnsNoErrorForValidMajor": {
			was:         "1.0.0",
			now:         "2.0.0",
			assertError: require.NoError,
		},
	}

	for name, testCase := range testCases {
		tc := testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := version.Compare(tc.was, tc.now)
			tc.assertError(t, err)
		})
	}
}
