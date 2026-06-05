package git

// Error is the error type.
type Error uint

const (
	// ErrNoGitTags is the error when no version tags are found in the
	// repository.
	ErrNoGitTags Error = iota + 1
)

// Error returns the error string for the error enum.
func (e Error) Error() string {
	switch e {
	case ErrNoGitTags:
		return "no git tags found"

	default:
		return "unknown error"
	}
}
