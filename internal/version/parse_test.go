package version_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tx3stn/vrsn/internal/test"
	"github.com/tx3stn/vrsn/internal/version"
)

func TestParse(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		input         string
		errorExpected require.ErrorAssertionFunc
		expected      version.SemVer
	}{
		"ReturnsVersionStructForValidInput": {
			input:         "34.9.154",
			errorExpected: require.NoError,
			expected: version.SemVer{
				Major: 34,
				Minor: 9,
				Patch: 154,
			},
		},
		"ReturnsErrorIfVersionDoesNotContainSeparator": {
			input:         "100",
			errorExpected: test.IsSentinelError(version.ErrNoVersionParts),
			expected:      version.SemVer{},
		},
		"ReturnsErrorIfInputDoesNotHaveThreeParts": {
			input:         "2.2",
			errorExpected: test.IsSentinelError(version.ErrNumVersionParts),
			expected:      version.SemVer{},
		},
		"ReturnsErrorIfMajorVersionCannotBeConvertedToInt": {
			input:         "x.1.1",
			errorExpected: test.IsSentinelError(version.ErrConvertingToInt),
			expected:      version.SemVer{},
		},
		"ReturnsErrorIfMinorVersionCannotBeConvertedToInt": {
			input:         "1.x.1",
			errorExpected: test.IsSentinelError(version.ErrConvertingToInt),
			expected:      version.SemVer{},
		},
		"ReturnsErrorIfPatchVersionCannotBeConvertedToInt": {
			input:         "1.5.x",
			errorExpected: test.IsSentinelError(version.ErrConvertingToInt),
			expected:      version.SemVer{},
		},
	}

	for name, testCase := range testCases {
		tc := testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			actual, err := version.Parse(tc.input)
			tc.errorExpected(t, err)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
