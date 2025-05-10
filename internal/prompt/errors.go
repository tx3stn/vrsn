package prompt

// Error is the error type.
type Error uint8

const (
	// ErrGettingBumpOptions is the error returned when something goes wrong getting
	// the version bump options for the specified string.
	ErrGettingBumpOptions Error = iota + 1
)

// Error returns the message string for the given error.
func (e Error) Error() string {
	switch e {
	case ErrGettingBumpOptions:
		return "error getting bump options"

	default:
		return "unknown error"
	}
}
