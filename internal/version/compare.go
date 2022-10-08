package version

// Compare compares the provided versions to see if the increase is a valid
// semver increment.
func Compare(wasInput string, nowInput string) error {
	if wasInput == nowInput {
		return ErrVersionNotBumped
	}

	was, err := Parse(wasInput)
	if err != nil {
		return err
	}

	now, err := Parse(nowInput)
	if err != nil {
		return err
	}

	if IsValidPatch(was, now) {
		return nil
	}

	if IsValidMinor(was, now) {
		return nil
	}

	if IsValidMajor(was, now) {
		return nil
	}

	return ErrInvalidBump
}
