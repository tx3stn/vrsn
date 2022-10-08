package flags

// Validate checks that both required values have been supplied.
func Validate(was string, now string) error {
	if was == "" && now == "" {
		return ErrNoValues
	}

	if was == "" {
		return ErrNoWasValue
	}

	if now == "" {
		return ErrNoNowValue
	}

	return nil
}
