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
	// ErrGitTagFileNoCommit is the error when '--git-tag' and '--file' are combined
	// but commit is disabled, so there is no version bump commit to tag.
	ErrGitTagFileNoCommit
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

	case ErrGitTagFileNoCommit:
		return "cannot combine --git-tag with --file unless commit is enabled (the tag must point at the version bump commit)"

	default:
		return "unknown error"
	}
}
