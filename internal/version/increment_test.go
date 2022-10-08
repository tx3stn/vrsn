package version_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tx3stn/vrsn/internal/version"
)

func TestGetBumpOptions(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		version     string
		assertError require.ErrorAssertionFunc
		expected    version.BumpOptions
	}{
		"ReturnsErrorForInvalidVersionString": {
			version:     "foo",
			assertError: require.Error,
			expected:    version.BumpOptions{},
		},
		"ReturnsIncrementedVersionsForValidInput": {
			version:     "1.0.0",
			assertError: require.NoError,
			expected: version.BumpOptions{
				Major: "2.0.0",
				Minor: "1.1.0",
				Patch: "1.0.1",
			},
		},
	}

	for name, testCase := range testCases {
		tc := testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			actual, err := version.GetBumpOptions(tc.version)
			tc.assertError(t, err)

			assert.Equal(t, tc.expected, actual)
		})
	}
}
