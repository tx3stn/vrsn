// Package prompt contains logic for prompting user interaction.
package prompt

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/tx3stn/vrsn/internal/version"
)

// BumpTypeSelectorFunc is the type def for the selector func used in the BumpSelector struct.
type BumpTypeSelectorFunc func(version.BumpOptions) (string, error)

// BumpSelector is a utility struct to enable mocking calls of the selector prompt for easier testability.
type BumpSelector struct {
	SelectorFunc BumpTypeSelectorFunc
}

// NewBumpSelector creates a new instance of the bump selector.
func NewBumpSelector() BumpSelector {
	return BumpSelector{
		SelectorFunc: selectBumpType,
	}
}

// Select prompts the user to select a bump type.
func (b BumpSelector) Select(currentVersion string) (string, error) {
	versionOptions, err := version.GetBumpOptions(currentVersion)
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrGettingBumpOptions, err)
	}

	selected, err := b.SelectorFunc(versionOptions)
	if err != nil {
		return "", err
	}

	//nolint:wrapcheck
	return versionOptions.SelectedIncrement(selected)
}

// selectBumpType prompts the user to select the type of version increment they
// wish to use.
func selectBumpType(opts version.BumpOptions) (string, error) {
	var selected string

	prompt := huh.NewSelect[string]().
		Options(huh.NewOptions(opts.PromptOptions()...)...).
		Title("select version bump type:").
		Value(&selected)

	if err := prompt.Run(); err != nil {
		return "", fmt.Errorf("error prompting for bump type: %w", err)
	}

	fmt.Println(strings.ReplaceAll(prompt.View(), "\n", ""))

	return selected, nil
}
