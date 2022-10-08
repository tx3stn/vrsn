package version

// Error is the error type.
type Error uint

const (
	// ErrConvertingToInt is the error thrown when a version part cannot be
	// converted to a string.
	ErrConvertingToInt Error = iota + 1
	// ErrInvalidBump is the error when the version has changed but to a value
	// that is not valid sem ver.
	ErrInvalidBump
	// ErrNoVersionParts is the error when the version string does not contain any
	// '.' to split into version parts.
	ErrNoVersionParts
	// ErrNumVersionParts is the error if the semantic version does not contain
	// three parts separated by '.'.
	ErrNumVersionParts
	// ErrVersionNotBumped is the error when the version has not been bumped.
	ErrVersionNotBumped
)

// Error returns the error string for the error enum.
func (e Error) Error() string {
	switch e {
	case ErrConvertingToInt:
		return "error converting version part to int"

	case ErrInvalidBump:
		return "invalid version bump"

	case ErrNumVersionParts:
		return "invalid number of version parts"

	case ErrNoVersionParts:
		return "version string does not contain any . splitting version segments"

	case ErrVersionNotBumped:
		return "version has not been bumped"

	default:
		return "unknown error"
	}
}
