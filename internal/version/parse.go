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

	majorPart := parts[0]
	if majorPart[0:1] == prefix {
		majorPart = parts[0][1:]
		pre = prefix
	}

	major, err := strconv.Atoi(majorPart)
	if err != nil {
		return SemVer{}, fmt.Errorf("%w: major version :%w", ErrConvertingToInt, err)
	}

	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		return SemVer{}, fmt.Errorf("%w: minor version :%w", ErrConvertingToInt, err)
	}

	patch, err := strconv.Atoi(parts[2])
	if err != nil {
		return SemVer{}, fmt.Errorf("%w: patch version :%w", ErrConvertingToInt, err)
	}

	return SemVer{
		Major:  major,
		Minor:  minor,
		Patch:  patch,
		Prefix: pre,
	}, nil
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

// ToString returns the string representation of a SemVer struct.
func (s *SemVer) ToString() string {
	return fmt.Sprintf("%s%d.%d.%d", s.Prefix, s.Major, s.Minor, s.Patch)
}
