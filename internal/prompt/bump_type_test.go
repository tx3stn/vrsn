package prompt_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tx3stn/vrsn/internal/prompt"
	"github.com/tx3stn/vrsn/internal/version"
)

func TestSelect(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		currentVersion string
		selectorFunc   prompt.BumpTypeSelectorFunc
		expected       string
		expectedError  error
	}{
		"returns the selected version string": {
			currentVersion: "6.6.6",
			selectorFunc: func(opts version.BumpOptions) (string, error) {
				return "minor", nil
			},
			expected:      "6.7.0",
			expectedError: nil,
		},
		"returns error for invalid version string": {
			currentVersion: "",
			selectorFunc: func(opts version.BumpOptions) (string, error) {
				return "", nil
			},
			expected:      "",
			expectedError: prompt.ErrGettingBumpOptions,
		},
	}

	for name, testCase := range testCases {
		tc := testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			bumpSelector := prompt.BumpSelector{
				SelectorFunc: tc.selectorFunc,
			}

			actual, err := bumpSelector.Select(tc.currentVersion)
			require.ErrorIs(t, err, tc.expectedError)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
