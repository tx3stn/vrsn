package files

// Error is the error type.
type Error uint

const (
	// ErrNoVersionFilesInDir is the error when no version files are found.
	ErrNoVersionFilesInDir Error = iota
	// ErrMultipleVersionFiles is the error when there are multiple valid version
	// file types found in a directory.
	ErrMultipleVersionFiles
	// ErrGettingVersionFromCMakeLists is the error when the version can't be
	// found inside a CMakeLists.txt file.
	ErrGettingVersionFromCMakeLists
	// ErrGettingVersionFromBuildGradle is the error when the a version key can't
	// be found inside a build.gradle or build.gradle.kts file.
	ErrGettingVersionFromBuildGradle
	// ErrGettingVersionFromPackageJSON is the error when a version key can't be
	// found inside a package.json file.
	ErrGettingVersionFromPackageJSON
	// ErrGettingVersionFromSetupPy is the error when a version key can't be found
	// inside a setup.py file.
	ErrGettingVersionFromSetupPy
	// ErrGettingVersionFromTOML is the error when a version key can't be found
	// inside a toml file.
	ErrGettingVersionFromTOML
	// ErrGettingVersionFromVERSION is the error when the VERSION file is empty.
	ErrGettingVersionFromVERSION
)

// Error returns the error string for the error enum.
func (e Error) Error() string {
	switch e {
	case ErrNoVersionFilesInDir:
		return "no version files found in directory"

	case ErrMultipleVersionFiles:
		return "multiple version files found in directory, use the --file flag to select the specific file to use"

	case ErrGettingVersionFromBuildGradle:
		return "unable to read version from build.gradle"

	case ErrGettingVersionFromCMakeLists:
		return "unable to read version from CMakeLists.txt"

	case ErrGettingVersionFromPackageJSON:
		return "unable to read version from package.json"

	case ErrGettingVersionFromSetupPy:
		return "unable to read version from setup.py"

	case ErrGettingVersionFromTOML:
		return "unable to read version from toml file"

	case ErrGettingVersionFromVERSION:
		return "unable to read version from VERSION file"

	default:
		return "unknown error"
	}
}
