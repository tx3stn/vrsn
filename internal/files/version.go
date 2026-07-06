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
	// secondary describes an additional value updated alongside the primary
	// version (e.g. android:versionCode). It is nil for the single-field
	// formats and only applied when a value is supplied to the writer.
	secondary *secondaryField
}

// secondaryField is an extra value written alongside the primary version. Its
// regex must have the same 4-group shape as versionRegex, i.e.
// (prefix)(key=")(value)(suffix), so the writer can replace the value by index.
type secondaryField struct {
	lineMatcher   func(string) bool
	notFoundError error
	regex         *regexp.Regexp
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

var bazelModMatcher = versionFileMatcher{
	lineMatcher:    tomlMatcher.lineMatcher,
	notFoundError:  ErrGettingVersionFromBazelModule,
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

// androidManifestMatcher extracts the version from the android:versionName
// attribute in an AndroidManifest.xml file. It optionally updates the
// android:versionCode attribute when a version code is supplied to the writer.
var androidManifestMatcher = versionFileMatcher{
	lineMatcher: func(line string) bool {
		return strings.Contains(line, "android:versionName")
	},
	notFoundError:  ErrGettingVersionFromAndroidManifest,
	singleLineFile: false,
	versionRegex: regexp.MustCompile(
		`(.*)(android:versionName\s*=\s*")(?P<semver>v*\d+\.\d+\.\d+)(".*)`,
	),
	secondary: &secondaryField{
		lineMatcher: func(line string) bool {
			return strings.Contains(line, "android:versionCode")
		},
		notFoundError: ErrGettingVersionCodeFromAndroidManifest,
		regex:         regexp.MustCompile(`(.*)(android:versionCode\s*=\s*")(\d+)(".*)`),
	},
}

// versionFileMatchers contains the utilities to extract and update the
// version from each supported version file.
var versionFileMatchers = map[string]versionFileMatcher{
	"BUILD.bazel":      bazelMatcher,
	"MODULE.bazel":     bazelModMatcher,
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

// patternMatchers holds matchers for version files identified by a filename
// glob rather than an exact name, e.g. AndroidManifest.xml and its variants
// (AndroidManifest.debug.xml). They are consulted only when the exact-name
// map has no entry.
var patternMatchers = []struct {
	pattern string
	matcher versionFileMatcher
}{
	{pattern: "AndroidManifest*.xml", matcher: androidManifestMatcher},
}

// lookupVersionFileMatcher resolves the matcher for a base filename, checking
// exact names first then filename patterns.
func lookupVersionFileMatcher(name string) (versionFileMatcher, bool) {
	if matcher, exists := versionFileMatchers[name]; exists {
		return matcher, true
	}

	for _, pm := range patternMatchers {
		if matched, _ := filepath.Match(pm.pattern, name); matched {
			return pm.matcher, true
		}
	}

	return versionFileMatcher{}, false
}

// getVersionMatcher gets the relevant versionFileMatcher config for the
// provided input file, falling back to the best effort matcher if there is no
// config for a file with that name.
func getVersionMatcher(inputFile string) versionFileMatcher {
	// Split dir and file to support relative paths provided with `--file` CLI flag.
	_, file := filepath.Split(inputFile)

	matcher, exists := lookupVersionFileMatcher(file)
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
	opts WriteOptions,
) ([]string, error) {
	if v.singleLineFile {
		return []string{opts.NewVersion}, nil
	}

	// The secondary field is only updated when the matcher defines one and a
	// value is supplied, so ordinary formats are unaffected.
	updateSecondary := v.secondary != nil && opts.AndroidVersionCode != ""
	foundVersion := false
	foundSecondary := false
	allLines := []string{}

	for scanner.Scan() {
		lineText := scanner.Text()

		// Only replace the first matching line, mirroring getVersion which
		// reads the first match. Later matches can be unrelated, e.g.
		// dependency version constraints in Cargo.toml or pyproject.toml.
		if !foundVersion && v.lineMatcher(lineText) {
			lineText = v.versionRegex.ReplaceAllString(
				lineText,
				fmt.Sprintf(`${1}${2}%s${4}`, opts.NewVersion),
			)
			foundVersion = true
		}

		// The secondary value (e.g. android:versionCode) may share the version
		// line or be on its own, so it is checked independently of the primary.
		if updateSecondary && !foundSecondary && v.secondary.lineMatcher(lineText) {
			lineText = v.secondary.regex.ReplaceAllString(
				lineText,
				fmt.Sprintf(`${1}${2}%s${4}`, opts.AndroidVersionCode),
			)
			foundSecondary = true
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

	if updateSecondary && !foundSecondary {
		return []string{}, v.secondary.notFoundError
	}

	return allLines, nil
}
