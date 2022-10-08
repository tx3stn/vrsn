package flags

// Error is the error type.
type Error uint

const (
	// ErrNoValues is the error when both was or now are not supplied.
	ErrNoValues Error = iota
	// ErrNoNowValue is the error when no --now value is supplied.
	ErrNoNowValue
	// ErrNoWasValue is the error when no --was value iss supplied.
	ErrNoWasValue
)

// Error returns the error string for the error enum.
func (e Error) Error() string {
	switch e {
	case ErrNoNowValue:
		return "no --now value provided"

	case ErrNoValues:
		return "no values provided for --was and --now"

	case ErrNoWasValue:
		return "no --was value provided"

	default:
		return "unknown error"
	}
}
