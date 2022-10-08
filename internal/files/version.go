package files

import (
	"bufio"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
)

type versionFileMatcher struct {
	lineMatcher    func(string) bool
	notFoundError  error
	singleLineFile bool
	versionRegex   string
}

var tomlMatcher = versionFileMatcher{
	lineMatcher: func(line string) bool {
		return strings.Contains(line, "version =")
	},
	notFoundError:  ErrGettingVersionFromTOML,
	singleLineFile: false,
	versionRegex:   `(.*)(version *=* *"*)(?P<semver>\d+.\d+.\d+)(.*)`,
}

// not toml files, but version string is same format.
var gradleMatcher = versionFileMatcher{
	lineMatcher:    tomlMatcher.lineMatcher,
	notFoundError:  ErrGettingVersionFromBuildGradle,
	singleLineFile: false,
	versionRegex:   tomlMatcher.versionRegex,
}

// versionFileMatchers contains the utilies to extract and update the version
// from the version file.
func versionFileMatchers() map[string]versionFileMatcher {
	return map[string]versionFileMatcher{
		"build.gradle":     gradleMatcher,
		"build.gradle.kts": gradleMatcher,
		"Cargo.toml":       tomlMatcher,
		"CMakeLists.txt": {
			lineMatcher: func(line string) bool {
				return strings.Contains(line, "project(")
			},
			notFoundError:  ErrGettingVersionFromCMakeLists,
			singleLineFile: false,
			versionRegex:   `(project\(.*)(VERSION ){1}(?P<semver>\d+.\d+.\d+)(.*\))`,
		},
		"package.json": {
			lineMatcher: func(line string) bool {
				return strings.Contains(line, `"version": "`)
			},
			notFoundError:  ErrGettingVersionFromPackageJSON,
			singleLineFile: false,
			versionRegex:   `(.*)("version": *"){1}(?P<semver>\d+.\d+.\d+)(".*)`,
		},
		"pyproject.toml": tomlMatcher,
		"setup.py": {
			lineMatcher: func(line string) bool {
				return strings.Contains(line, `version=`)
			},
			notFoundError:  ErrGettingVersionFromSetupPy,
			singleLineFile: false,
			versionRegex:   `(.*)(version=['"])(?P<semver>\d+.\d+.\d+)(.*)`,
		},
		"VERSION": {
			lineMatcher: func(line string) bool {
				// single line file so nothing to match on.
				return true
			},
			notFoundError:  ErrGettingVersionFromVERSION,
			singleLineFile: true,
			versionRegex:   `(.*)(?P<semver>\d+.\d+.\d+)(.*)`,
		},
	}
}

// getVersionMatcher gets the relevant versionFileMatcher config for the
// provided input file or errors if there is no config for a file with that name.
func getVersionMatcher(inputFile string) (versionFileMatcher, error) {
	// Split dir and file to support relative paths provided with `--file` CLI flag.
	_, file := filepath.Split(inputFile)

	matcher, exists := versionFileMatchers()[file]
	if !exists {
		return versionFileMatcher{}, fmt.Errorf("%s is not a supported version file type", file)
	}

	return matcher, nil
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
			re := regexp.MustCompile(v.versionRegex)
			result := make(map[string]string)

			match := re.FindStringSubmatch(lineText)
			if match == nil {
				return "", v.notFoundError
			}

			for i, name := range re.SubexpNames() {
				if i != 0 && name != "" {
					result[name] = match[i]
				}
			}

			semver, exists := result["semver"]
			if !exists {
				return "", v.notFoundError
			}

			return semver, nil
		}
	}

	return "", v.notFoundError
}

func (v versionFileMatcher) updateVersionInPlace(scanner *bufio.Scanner, newVersion string) ([]string, error) {
	if v.singleLineFile {
		return []string{newVersion}, nil
	}

	foundVersion := false
	allLines := []string{}

	for scanner.Scan() {
		lineText := scanner.Text()

		if v.lineMatcher(lineText) {
			re := regexp.MustCompile(v.versionRegex)
			newVersionLine := re.ReplaceAllString(lineText, fmt.Sprintf(`${1}${2}%s${4}`, newVersion))
			allLines = append(allLines, newVersionLine)
			foundVersion = true
			continue
		}

		allLines = append(allLines, lineText)
	}

	if !foundVersion {
		return []string{}, v.notFoundError
	}

	return allLines, nil
}
