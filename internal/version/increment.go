package version

import (
	"errors"
	"fmt"
	"strings"
)

// BumpOptions contains details about the bump options.
type BumpOptions struct {
	Major string
	Minor string
	Patch string
}

// GetBumpOptions returns the possible valid version bump options from the
// input string.
func GetBumpOptions(inputVersion string) (BumpOptions, error) {
	parsed, err := Parse(inputVersion)
	if err != nil {
		return BumpOptions{}, err
	}

	major := parsed
	major.MajorBump()

	minor := parsed
	minor.MinorBump()

	patch := parsed
	patch.PatchBump()

	return BumpOptions{
		Patch: patch.ToString(),
		Minor: minor.ToString(),
		Major: major.ToString(),
	}, nil
}

// PromptOptions returns the options formatted for a user prompt.
func (b BumpOptions) PromptOptions() []string {
	return []string{
		b.formattedPatch(),
		b.formattedMinor(),
		b.formattedMajor(),
	}
}

// SelectedIncrement gets ust the version number from the user selected prompt.
func (b BumpOptions) SelectedIncrement(increment string) (string, error) {
	if strings.Contains(increment, "patch") {
		return b.Patch, nil
	}
	if strings.Contains(increment, "minor") {
		return b.Minor, nil
	}
	if strings.Contains(increment, "major") {
		return b.Major, nil
	}

	return "", errors.New("invalid increment type")
}

func (b BumpOptions) formattedMajor() string {
	return fmt.Sprintf("major (%s)", b.Major)
}

func (b BumpOptions) formattedMinor() string {
	return fmt.Sprintf("minor (%s)", b.Minor)
}

func (b BumpOptions) formattedPatch() string {
	return fmt.Sprintf("patch (%s)", b.Patch)
}
