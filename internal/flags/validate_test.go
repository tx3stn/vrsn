package flags_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tx3stn/vrsn/internal/flags"
)

func TestValidate(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		was           string
		now           string
		expectedError error
	}{
		"ReturnsErrorWhenNoWasOrNowSupplied": {
			was:           "",
			now:           "",
			expectedError: flags.ErrNoValues,
		},
		"ReturnsErrorWhenNoWasSupplied": {
			was:           "",
			now:           "1.0.0",
			expectedError: flags.ErrNoWasValue,
		},
		"ReturnsErrorWhenNoNowSupplied": {
			was:           "2.4.6",
			now:           "",
			expectedError: flags.ErrNoNowValue,
		},
		"ReturnsNoErrorWhenBothSupplied": {
			was:           "6.6.6",
			now:           "9.9.9",
			expectedError: nil,
		},
	}

	for name, testCase := range testCases {
		tc := testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := flags.Validate(tc.was, tc.now)
			require.ErrorIs(t, err, tc.expectedError)
		})
	}
}
