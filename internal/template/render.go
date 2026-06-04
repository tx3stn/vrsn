// Package template handles rendering of user provided message templates.
package template

import (
	"bytes"
	"fmt"
	"text/template"
)

// Data holds the variables available to message templates.
type Data struct {
	Version string
}

// Render renders the provided message template, exposing the new version via
// the {{.Version}} template variable. Messages that don't use any template
// syntax are returned unchanged.
func Render(msg string, version string) (string, error) {
	tmpl, err := template.New("message").Option("missingkey=error").Parse(msg)
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrParsingTemplate, err)
	}

	var rendered bytes.Buffer
	if err := tmpl.Execute(&rendered, Data{Version: version}); err != nil {
		return "", fmt.Errorf("%w: %w", ErrRenderingTemplate, err)
	}

	return rendered.String(), nil
}
