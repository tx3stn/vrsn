// Package prompt contains logic for prompting user interaction.
package prompt

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/tx3stn/vrsn/internal/version"
)

// SelectBumpType prompts the user to select the type of version increment they
// wish to use.
func SelectBumpType(currentVersion string) (string, error) {
	versionOptions, err := version.GetBumpOptions(currentVersion)
	if err != nil {
		return "", fmt.Errorf("error getting bump options: %w", err)
	}

	answer := struct {
		Selected string `survey:"bump"`
	}{}

	err = survey.Ask([]*survey.Question{
		{
			Name: "bump",
			Prompt: &survey.Select{
				Message: "select version bump type:",
				Options: versionOptions.PromptOptions(),
			},
		},
	}, &answer)
	if err != nil {
		return "", fmt.Errorf("error prompting to selection version bump type: %w", err)
	}

	//nolint:wrapcheck
	return versionOptions.SelectedIncrement(answer.Selected)
}
