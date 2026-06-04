package template_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tx3stn/vrsn/internal/template"
)

func TestRender(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		message       string
		version       string
		expected      string
		expectedError error
	}{
		"ReturnsMessageWithoutTemplateSyntaxUnchanged": {
			message:       "bump version",
			version:       "1.0.0",
			expected:      "bump version",
			expectedError: nil,
		},
		"RendersVersionVariable": {
			message:       "release {{.Version}}",
			version:       "1.2.3",
			expected:      "release 1.2.3",
			expectedError: nil,
		},
		"RendersVersionVariableWithPrefix": {
			message:       "bump version to {{.Version}}",
			version:       "v1.2.3",
			expected:      "bump version to v1.2.3",
			expectedError: nil,
		},
		"RendersVersionVariableWithSpacing": {
			message:       "release {{ .Version }}",
			version:       "0.1.0",
			expected:      "release 0.1.0",
			expectedError: nil,
		},
		"ReturnsErrorForInvalidTemplateSyntax": {
			message:       "release {{.Version",
			version:       "1.2.3",
			expected:      "",
			expectedError: template.ErrParsingTemplate,
		},
		"ReturnsErrorForUnsupportedVariable": {
			message:       "release {{.NotAThing}}",
			version:       "1.2.3",
			expected:      "",
			expectedError: template.ErrRenderingTemplate,
		},
	}

	for name, testCase := range testCases {
		tc := testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			rendered, err := template.Render(tc.message, tc.version)
			require.ErrorIs(t, err, tc.expectedError)
			assert.Equal(t, tc.expected, rendered)
		})
	}
}
