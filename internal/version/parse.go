// Package version holds logic for validating an interacting with a single version.
package version

import (
	"fmt"
	"strconv"
	"strings"
)

// SemVer holds the details of the semantic version parts.
type SemVer struct {
	Major  int
	Minor  int
	Patch  int
	Prefix string
}

const (
	prefix      = "v"
	semVerParts = 3
)

// Parse checks the input string is a valid semantic version and
// parses it into a SemVer struct.
func Parse(version string) (SemVer, error) {
	if !strings.Contains(version, ".") {
		return SemVer{}, ErrNoVersionParts
	}

	parts := strings.Split(version, ".")
	if len(parts) != semVerParts {
		return SemVer{}, ErrNumVersionParts
	}

	pre := ""

	if majorPart, found := strings.CutPrefix(parts[0], prefix); found {
		parts[0] = majorPart
		pre = prefix
	}

	major, err := parsePart(parts[0], "major")
	if err != nil {
		return SemVer{}, err
	}

	minor, err := parsePart(parts[1], "minor")
	if err != nil {
		return SemVer{}, err
	}

	patch, err := parsePart(parts[2], "patch")
	if err != nil {
		return SemVer{}, err
	}

	return SemVer{
		Major:  major,
		Minor:  minor,
		Patch:  patch,
		Prefix: pre,
	}, nil
}

// parsePart converts a single version part to an int, rejecting anything
// that isn't a non-negative number.
func parsePart(part string, name string) (int, error) {
	num, err := strconv.Atoi(part)
	if err != nil {
		return 0, fmt.Errorf("%w: %s version: %w", ErrConvertingToInt, name, err)
	}

	if num < 0 {
		return 0, fmt.Errorf("%w: %s version: negative number", ErrConvertingToInt, name)
	}

	return num, nil
}

// MajorBump increments the major version by 1.
func (s *SemVer) MajorBump() {
	s.Major++
	s.Minor = 0
	s.Patch = 0
}

// MinorBump increments the minor version by 1.
func (s *SemVer) MinorBump() {
	s.Minor++
	s.Patch = 0
}

// PatchBump increments the patch version by 1.
func (s *SemVer) PatchBump() {
	s.Patch++
}

// String returns the string representation of a SemVer struct.
func (s *SemVer) String() string {
	return fmt.Sprintf("%s%d.%d.%d", s.Prefix, s.Major, s.Minor, s.Patch)
}

const (
	androidMajorMultiplier = 10000
	androidMinorMultiplier = 100
)

// AndroidVersionCode derives an Android versionCode integer from the semantic
// version using the conventional MAJOR*10000 + MINOR*100 + PATCH scheme. This
// reserves two digits each for the minor and patch parts, so it assumes both
// are below 100.
func (s *SemVer) AndroidVersionCode() int {
	return s.Major*androidMajorMultiplier + s.Minor*androidMinorMultiplier + s.Patch
}
