package files

import (
	"bufio"
	"fmt"
	"io"
	"path/filepath"
	"regexp"
	"strings"
)

// maxLineBytes is the max line length the version file scanners support,
// larger than the bufio default so long lines (e.g. in a minified
// package.json) error cleanly rather than being truncated.
const maxLineBytes = 1024 * 1024

// newScanner creates a line scanner for reading version files.
func newScanner(reader io.Reader) *bufio.Scanner {
	scanner := bufio.NewScanner(reader)
	scanner.Buffer(make([]byte, 0, bufio.MaxScanTokenSize), maxLineBytes)

	return scanner
}

type versionFileMatcher struct {
	lineMatcher    func(string) bool
	notFoundError  error
	singleLineFile bool
	versionRegex   *regexp.Regexp
}

// tomlVersionLine matches a version key at the start of the line so
// dependency constraints like `tokio = { version = "1.0.0" }` are ignored.
var tomlVersionLine = regexp.MustCompile(`^\s*version\s*=`)

var tomlMatcher = versionFileMatcher{
	lineMatcher:    tomlVersionLine.MatchString,
	notFoundError:  ErrGettingVersionFromTOML,
	singleLineFile: false,
	versionRegex:   regexp.MustCompile(`(.*)(version\s*=\s*['"]?)(?P<semver>v*\d+\.\d+\.\d+)(.*)`),
}

// not a toml file, but version attribute is same format.
var bazelMatcher = versionFileMatcher{
	lineMatcher:    tomlMatcher.lineMatcher,
	notFoundError:  ErrGettingVersionFromBuildBazel,
	singleLineFile: false,
	versionRegex:   tomlMatcher.versionRegex,
}

// not toml files, but version string is same format.
var gradleMatcher = versionFileMatcher{
	lineMatcher:    tomlMatcher.lineMatcher,
	notFoundError:  ErrGettingVersionFromBuildGradle,
	singleLineFile: false,
	versionRegex:   tomlMatcher.versionRegex,
}

// bestEffortRegex matches toml style `version = X` lines with single, double
// or no quotes.
var bestEffortRegex = regexp.MustCompile(`(.*)(version\s*=\s*['"]?)(?P<semver>v*\d+\.\d+\.\d+)(.*)`)

// bestEffortMatcher is the fallback for files explicitly provided with the
// --file flag that don't match any of the supported version files.
// Using the regex as the lineMatcher means a line only matches when the
// version can actually be extracted from it.
var bestEffortMatcher = versionFileMatcher{
	lineMatcher:    bestEffortRegex.MatchString,
	notFoundError:  ErrGettingVersionBestEffort,
	singleLineFile: false,
	versionRegex:   bestEffortRegex,
}

// versionFileMatchers contains the utilities to extract and update the
// version from each supported version file.
var versionFileMatchers = map[string]versionFileMatcher{
	"BUILD.bazel":      bazelMatcher,
	"build.gradle":     gradleMatcher,
	"build.gradle.kts": gradleMatcher,
	"Cargo.toml":       tomlMatcher,
	"CMakeLists.txt": {
		lineMatcher: func(line string) bool {
			return strings.Contains(line, "project(")
		},
		notFoundError:  ErrGettingVersionFromCMakeLists,
		singleLineFile: false,
		versionRegex: regexp.MustCompile(
			`(project\(.*)(VERSION\s+)(?P<semver>v*\d+\.\d+\.\d+)(.*\))`,
		),
	},
	"package.json": {
		lineMatcher: func(line string) bool {
			return strings.Contains(line, `"version":`)
		},
		notFoundError:  ErrGettingVersionFromPackageJSON,
		singleLineFile: false,
		versionRegex:   regexp.MustCompile(`(.*)("version":\s*")(?P<semver>v*\d+\.\d+\.\d+)(".*)`),
	},
	"pyproject.toml": tomlMatcher,
	"setup.py": {
		lineMatcher: func(line string) bool {
			return strings.Contains(line, `version=`)
		},
		notFoundError:  ErrGettingVersionFromSetupPy,
		singleLineFile: false,
		versionRegex:   regexp.MustCompile(`(.*)(version=['"])(?P<semver>v*\d+\.\d+\.\d+)(.*)`),
	},
	"VERSION": {
		lineMatcher: func(line string) bool {
			// single line file so nothing to match on.
			return true
		},
		notFoundError:  ErrGettingVersionFromVERSION,
		singleLineFile: true,
		versionRegex:   regexp.MustCompile(`(.*)(?P<semver>v*\d+\.\d+\.\d+)(.*)`),
	},
}

// getVersionMatcher gets the relevant versionFileMatcher config for the
// provided input file, falling back to the best effort matcher if there is no
// config for a file with that name.
func getVersionMatcher(inputFile string) versionFileMatcher {
	// Split dir and file to support relative paths provided with `--file` CLI flag.
	_, file := filepath.Split(inputFile)

	matcher, exists := versionFileMatchers[file]
	if !exists {
		return bestEffortMatcher
	}

	return matcher
}

func (v versionFileMatcher) getVersion(scanner *bufio.Scanner) (string, error) {
	for scanner.Scan() {
		lineText := scanner.Text()

		if v.singleLineFile && (lineText == "" || lineText == "\n") {
			return "", v.notFoundError
		}

		if v.singleLineFile {
			return lineText, nil
		}

		if v.lineMatcher(lineText) {
			semver, found := v.extractVersion(lineText)
			if !found {
				return "", v.notFoundError
			}

			return semver, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error reading version file: %w", err)
	}

	return "", v.notFoundError
}

// extractVersion pulls the semver capture group out of the version line.
func (v versionFileMatcher) extractVersion(lineText string) (string, bool) {
	match := v.versionRegex.FindStringSubmatch(lineText)
	if match == nil {
		return "", false
	}

	semverIndex := v.versionRegex.SubexpIndex("semver")
	if semverIndex == -1 || match[semverIndex] == "" {
		return "", false
	}

	return match[semverIndex], true
}

func (v versionFileMatcher) updateVersionInPlace(
	scanner *bufio.Scanner,
	newVersion string,
) ([]string, error) {
	if v.singleLineFile {
		return []string{newVersion}, nil
	}

	foundVersion := false
	allLines := []string{}

	for scanner.Scan() {
		lineText := scanner.Text()

		// Only replace the first matching line, mirroring getVersion which
		// reads the first match. Later matches can be unrelated, e.g.
		// dependency version constraints in Cargo.toml or pyproject.toml.
		if !foundVersion && v.lineMatcher(lineText) {
			newVersionLine := v.versionRegex.ReplaceAllString(
				lineText,
				fmt.Sprintf(`${1}${2}%s${4}`, newVersion),
			)
			allLines = append(allLines, newVersionLine)
			foundVersion = true

			continue
		}

		allLines = append(allLines, lineText)
	}

	// A scan error means the file was only partially read, so writing
	// allLines back would truncate the file.
	if err := scanner.Err(); err != nil {
		return []string{}, fmt.Errorf("error reading version file: %w", err)
	}

	if !foundVersion {
		return []string{}, v.notFoundError
	}

	return allLines, nil
}
