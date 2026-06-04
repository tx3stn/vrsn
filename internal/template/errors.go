package template

// Error is the error type.
type Error uint

const (
	// ErrParsingTemplate is the error when the message template cannot be
	// parsed.
	ErrParsingTemplate Error = iota + 1
	// ErrRenderingTemplate is the error when the message template cannot be
	// rendered, such as when it references an unsupported variable.
	ErrRenderingTemplate
)

// Error returns the error string for the error enum.
func (e Error) Error() string {
	switch e {
	case ErrParsingTemplate:
		return "error parsing message template"

	case ErrRenderingTemplate:
		return "error rendering message template"

	default:
		return "unknown error"
	}
}
