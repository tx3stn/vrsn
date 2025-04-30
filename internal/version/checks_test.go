package version_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tx3stn/vrsn/internal/version"
)

//nolint:dupl
func TestIsValidMajor(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		was      version.SemVer
		now      version.SemVer
		expected bool
	}{
		"ReturnsTrueForValidMajor": {
			was:      version.SemVer{Major: 0, Minor: 0, Patch: 420},
			now:      version.SemVer{Major: 1, Minor: 0, Patch: 0},
			expected: true,
		},
		"ReturnsFalseWhenMajorTooHigh": {
			was:      version.SemVer{Major: 8, Minor: 3, Patch: 19},
			now:      version.SemVer{Major: 10, Minor: 0, Patch: 0},
			expected: false,
		},
		"ReturnsFalseWhenMinorIsNotReset": {
			was:      version.SemVer{Major: 2, Minor: 0, Patch: 4},
			now:      version.SemVer{Major: 3, Minor: 1, Patch: 0},
			expected: false,
		},
		"ReturnsFalseWhenPatchIsNotReset": {
			was:      version.SemVer{Major: 30, Minor: 812, Patch: 1},
			now:      version.SemVer{Major: 31, Minor: 0, Patch: 1},
			expected: false,
		},
	}

	for name, testCase := range testCases {
		tc := testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			actual := version.IsValidMajor(tc.was, tc.now)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

//nolint:dupl
func TestIsValidMinor(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		was      version.SemVer
		now      version.SemVer
		expected bool
	}{
		"ReturnsTrueForValidMinor": {
			was:      version.SemVer{Major: 19, Minor: 4, Patch: 23},
			now:      version.SemVer{Major: 19, Minor: 5, Patch: 0},
			expected: true,
		},
		"ReturnsFalseWhenMinorTooHigh": {
			was:      version.SemVer{Major: 1, Minor: 4, Patch: 8},
			now:      version.SemVer{Major: 1, Minor: 6, Patch: 0},
			expected: false,
		},
		"ReturnsFalseWhenMajorIsIncreased": {
			was:      version.SemVer{Major: 7, Minor: 1, Patch: 9573},
			now:      version.SemVer{Major: 8, Minor: 2, Patch: 0},
			expected: false,
		},
		"ReturnsFalseWhenPatchIsNotReset": {
			was:      version.SemVer{Major: 365, Minor: 19, Patch: 4},
			now:      version.SemVer{Major: 365, Minor: 20, Patch: 4},
			expected: false,
		},
	}

	for name, testCase := range testCases {
		tc := testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			actual := version.IsValidMinor(tc.was, tc.now)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

//nolint:dupl
func TestIsValidPatch(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		was      version.SemVer
		now      version.SemVer
		expected bool
	}{
		"ReturnsTrueForValidPatch": {
			was:      version.SemVer{Major: 1, Minor: 0, Patch: 4},
			now:      version.SemVer{Major: 1, Minor: 0, Patch: 5},
			expected: true,
		},
		"ReturnsFalseWhenPatchTooHigh": {
			was:      version.SemVer{Major: 0, Minor: 1, Patch: 4},
			now:      version.SemVer{Major: 0, Minor: 1, Patch: 6},
			expected: false,
		},
		"ReturnsFalseWhenMajorIsIncreased": {
			was:      version.SemVer{Major: 0, Minor: 1, Patch: 4},
			now:      version.SemVer{Major: 1, Minor: 1, Patch: 5},
			expected: false,
		},
		"ReturnsFalseWhenMinorIsIncreased": {
			was:      version.SemVer{Major: 0, Minor: 1, Patch: 4},
			now:      version.SemVer{Major: 0, Minor: 2, Patch: 5},
			expected: false,
		},
	}

	for name, testCase := range testCases {
		tc := testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			actual := version.IsValidPatch(tc.was, tc.now)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
