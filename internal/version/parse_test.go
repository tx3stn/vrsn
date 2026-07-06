package version_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tx3stn/vrsn/internal/version"
)

func TestParse(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		input         string
		expectedError error
		expected      version.SemVer
	}{
		"ReturnsVersionStructForValidInput": {
			input:         "34.9.154",
			expectedError: nil,
			expected: version.SemVer{
				Major: 34,
				Minor: 9,
				Patch: 154,
			},
		},
		"ReturnsErrorIfVersionDoesNotContainSeparator": {
			input:         "100",
			expectedError: version.ErrNoVersionParts,
			expected:      version.SemVer{},
		},
		"ReturnsErrorIfInputDoesNotHaveThreeParts": {
			input:         "2.2",
			expectedError: version.ErrNumVersionParts,
			expected:      version.SemVer{},
		},
		"ReturnsErrorIfMajorVersionCannotBeConvertedToInt": {
			input:         "x.1.1",
			expectedError: version.ErrConvertingToInt,
			expected:      version.SemVer{},
		},
		"ReturnsErrorIfMinorVersionCannotBeConvertedToInt": {
			input:         "1.x.1",
			expectedError: version.ErrConvertingToInt,
			expected:      version.SemVer{},
		},
		"ReturnsErrorIfPatchVersionCannotBeConvertedToInt": {
			input:         "1.5.x",
			expectedError: version.ErrConvertingToInt,
			expected:      version.SemVer{},
		},
		"ReturnsVersionStructForPrefixedInput": {
			input:         "v1.2.3",
			expectedError: nil,
			expected: version.SemVer{
				Major:  1,
				Minor:  2,
				Patch:  3,
				Prefix: "v",
			},
		},
		"ReturnsErrorForEmptyInput": {
			input:         "",
			expectedError: version.ErrNoVersionParts,
			expected:      version.SemVer{},
		},
		"ReturnsErrorForOnlySeparators": {
			input:         "..",
			expectedError: version.ErrConvertingToInt,
			expected:      version.SemVer{},
		},
		"ReturnsErrorForEmptyMajorPart": {
			input:         ".1.2",
			expectedError: version.ErrConvertingToInt,
			expected:      version.SemVer{},
		},
		"ReturnsErrorForPrefixOnlyMajorPart": {
			input:         "v.1.2",
			expectedError: version.ErrConvertingToInt,
			expected:      version.SemVer{},
		},
		"ReturnsErrorForNegativeVersionPart": {
			input:         "1.-2.3",
			expectedError: version.ErrConvertingToInt,
			expected:      version.SemVer{},
		},
	}

	for name, testCase := range testCases {
		tc := testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			actual, err := version.Parse(tc.input)
			require.ErrorIs(t, err, tc.expectedError)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestAndroidVersionCode(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		input    version.SemVer
		expected int
	}{
		"CombinesEachPart": {
			input:    version.SemVer{Major: 1, Minor: 2, Patch: 3},
			expected: 10203,
		},
		"PadsSingleDigitParts": {
			input:    version.SemVer{Major: 0, Minor: 0, Patch: 1},
			expected: 1,
		},
		"HandlesMajorOnly": {
			input:    version.SemVer{Major: 12, Minor: 0, Patch: 0},
			expected: 120000,
		},
		"IgnoresPrefix": {
			input:    version.SemVer{Major: 2, Minor: 14, Patch: 74, Prefix: "v"},
			expected: 21474,
		},
	}

	for name, testCase := range testCases {
		tc := testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tc.expected, tc.input.AndroidVersionCode())
		})
	}
}
