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
	// ErrGettingVersionFromBuildBazel is the error when a version attribute can't
	// be found inside a BUILD.bazel file.
	ErrGettingVersionFromBuildBazel
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
	// ErrGettingVersionBestEffort is the error when a version can't be found in
	// an unsupported file using best effort matching.
	ErrGettingVersionBestEffort
	// ErrGettingFilesInDirectory is the error returned when getting files in a directory fails.
	ErrGettingFilesInDirectory
	// ErrFileNotFound is the error returned when the specified version file cannot be found.
	ErrFileNotFound
	// ErrFileIsDirectory is the error returned when the specified version file is a directory.
	ErrFileIsDirectory
	// ErrVersionsDoNotMatch is the error returned when multiple version files
	// contain different versions so a single bump cannot be applied.
	ErrVersionsDoNotMatch
)

// Error returns the error string for the error enum.
//
//nolint:cyclop
func (e Error) Error() string {
	switch e {
	case ErrNoVersionFilesInDir:
		return "no version files found in directory"

	case ErrMultipleVersionFiles:
		return "multiple version files found in directory, use the --file flag to select the specific file to use"

	case ErrGettingVersionFromBuildBazel:
		return "unable to read version from BUILD.bazel"

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

	case ErrGettingVersionBestEffort:
		return "unable to read version using best effort matching, file is not a supported version file type"

	case ErrGettingFilesInDirectory:
		return "error getting files in directory"

	case ErrFileNotFound:
		return "file not found"

	case ErrFileIsDirectory:
		return "file is a directory"

	case ErrVersionsDoNotMatch:
		return "version files do not contain matching versions, run with --verbose to see the version in each file"

	default:
		return "unknown error"
	}
}
