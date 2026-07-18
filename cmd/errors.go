package cmd

// Error is the error type.
type Error uint

const (
	// ErrNoNowOrFile is the error when no version file can be found in the directory
	// and the '--now' flag was not passed.
	ErrNoNowOrFile Error = iota + 1
	// ErrNoWasOrFile is the error when no version file can be found in the directory
	// and the '--was' flag was not passed.
	ErrNoWasOrFile
	// ErrCantCompareVersionsOnBranch is the error when you are on the base branch and
	// no '--was' flag was passed so there is nothing to compare.
	ErrCantCompareVersionsOnBranch
	// ErrInvalidVersionSuffix is the error when the suffix after the first '-'
	// in a version passed to set contains characters other than letters, digits
	// and hyphens.
	ErrInvalidVersionSuffix
)

// Error returns the error string for the error enum.
func (e Error) Error() string {
	switch e {
	case ErrNoNowOrFile:
		return "please pass version with --now flag or run inside a directory that contains a version file"

	case ErrNoWasOrFile:
		return "please pass version with --was flag or run inside a directory that contains a version file"

	case ErrCantCompareVersionsOnBranch:
		return "on base branch with no --was flag supplied, nothing to compare"

	case ErrInvalidVersionSuffix:
		return "version suffix must contain only letters, digits and hyphens"

	default:
		return "unknown error"
	}
}
