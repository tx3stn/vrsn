package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tx3stn/vrsn/internal/config"
	"github.com/tx3stn/vrsn/internal/flags"
)

func TestGetFiles(t *testing.T) {
	testCases := map[string]struct {
		configFile      string
		versionFileFlag string
		expectedFiles   []string
	}{
		"ReturnsFilesFromConfigFile": {
			configFile:      "testdata/with-files/vrsn.toml",
			versionFileFlag: "",
			expectedFiles:   []string{"VERSION", "package.json"},
		},
		"ConfigFilesTakePrecedenceOverFileFlag": {
			configFile:      "testdata/with-files/vrsn.toml",
			versionFileFlag: "Cargo.toml",
			expectedFiles:   []string{"VERSION", "package.json"},
		},
		"ReturnsFileFlagWhenConfigFileHasNoFiles": {
			configFile:      "testdata/xdg/vrsn.toml",
			versionFileFlag: "Cargo.toml",
			expectedFiles:   []string{"Cargo.toml"},
		},
		"ReturnsFileFlagWhenNoConfigFileFound": {
			configFile:      "",
			versionFileFlag: "Cargo.toml",
			expectedFiles:   []string{"Cargo.toml"},
		},
		"ReturnsNoFilesWhenNoConfigFileAndNoFileFlag": {
			configFile:      "",
			versionFileFlag: "",
			expectedFiles:   nil,
		},
	}

	for name, testCase := range testCases {
		tc := testCase

		t.Run(name, func(t *testing.T) {
			// Ensure no real config file is picked up when no --config flag
			// is passed, t.Setenv also prevents the tests running in
			// parallel which keeps the mutation of the global flag var safe.
			t.Setenv("XDG_CONFIG_DIR", "")
			t.Setenv("XDG_CONFIG_HOME", "")
			t.Setenv("HOME", "")

			originalVersionFile := flags.VersionFile
			flags.VersionFile = tc.versionFileFlag

			t.Cleanup(func() {
				flags.VersionFile = originalVersionFile
			})

			conf, err := config.Get(tc.configFile, nil)
			require.NoError(t, err)
			assert.Equal(t, tc.expectedFiles, conf.Files)
		})
	}
}

// changedFlags is a test double for the FlagChecker interface.
type changedFlags map[string]bool

func (c changedFlags) Changed(name string) bool {
	return c[name]
}

func TestGetMergesConfigFileAndFlags(t *testing.T) {
	testCases := map[string]struct {
		configFile   string
		changed      changedFlags
		flagCommit   bool
		expectedBump config.BumpOpts
	}{
		"MissingConfigKeysKeepFlagDefaults": {
			configFile: "testdata/partial-bump/vrsn.toml",
			changed:    changedFlags{},
			flagCommit: false,
			expectedBump: config.BumpOpts{
				Commit:    true,
				CommitMsg: "bump version",
				GitTag:    false,
				TagMsg:    "",
			},
		},
		"ChangedFlagsOverrideConfigValues": {
			configFile: "testdata/with-files/vrsn.toml",
			changed:    changedFlags{"commit": true},
			flagCommit: true,
			expectedBump: config.BumpOpts{
				Commit:    true,
				CommitMsg: "bump version",
				GitTag:    false,
				TagMsg:    "",
			},
		},
		"UnchangedFlagsDoNotOverrideConfigValues": {
			configFile: "testdata/with-files/vrsn.toml",
			changed:    changedFlags{},
			flagCommit: true,
			expectedBump: config.BumpOpts{
				Commit:    false,
				CommitMsg: "bump version",
				GitTag:    false,
				TagMsg:    "",
			},
		},
	}

	for name, testCase := range testCases {
		tc := testCase

		t.Run(name, func(t *testing.T) {
			// t.Setenv also prevents the tests running in parallel which
			// keeps the mutation of the global flag vars safe.
			t.Setenv("XDG_CONFIG_DIR", "")
			t.Setenv("XDG_CONFIG_HOME", "")
			t.Setenv("HOME", "")

			originalCommit := flags.Commit
			originalCommitMsg := flags.CommitMsg
			originalBaseBranch := flags.BaseBranch
			flags.Commit = tc.flagCommit
			// Simulate the defaults cobra applies when registering the flags.
			flags.CommitMsg = "bump version"
			flags.BaseBranch = "main"

			t.Cleanup(func() {
				flags.Commit = originalCommit
				flags.CommitMsg = originalCommitMsg
				flags.BaseBranch = originalBaseBranch
			})

			conf, err := config.Get(tc.configFile, tc.changed)
			require.NoError(t, err)
			assert.Equal(t, tc.expectedBump, conf.Bump)
			assert.Equal(t, "main", conf.Check.BaseBranch)
		})
	}
}

func TestFindConfigFile(t *testing.T) {
	testCases := map[string]struct {
		chdir           string
		xdgEnvValue     string
		xdgHomeEnvValue string
		homeEnvValue    string
		expected        string
		expectedError   error
	}{
		"ReturnsCurrentDirectoryFileWhenExists": {
			chdir:         "testdata/project",
			xdgEnvValue:   "testdata/xdg/",
			homeEnvValue:  "testdata/home/",
			expected:      "vrsn.toml",
			expectedError: nil,
		},
		"ReturnsXdgFileWhenExists": {
			xdgEnvValue:   "testdata/xdg/",
			homeEnvValue:  "testdata/home/",
			expected:      "testdata/xdg/vrsn.toml",
			expectedError: nil,
		},
		"ReturnsXdgConfigHomeFileWhenExists": {
			xdgHomeEnvValue: "testdata/xdg/",
			homeEnvValue:    "testdata/home/",
			expected:        "testdata/xdg/vrsn.toml",
			expectedError:   nil,
		},
		"ReturnsHomeFileWhenExists": {
			xdgEnvValue:   "",
			homeEnvValue:  "testdata/home/",
			expected:      "testdata/home/.config/vrsn.toml",
			expectedError: nil,
		},
		"ReturnsEmptyStringWhenNoEnvVarsAreSet": {
			xdgEnvValue:   "",
			homeEnvValue:  "",
			expected:      "",
			expectedError: nil,
		},
	}

	for name, testCase := range testCases {
		tc := testCase

		t.Run(name, func(t *testing.T) {
			t.Setenv("XDG_CONFIG_DIR", tc.xdgEnvValue)
			t.Setenv("XDG_CONFIG_HOME", tc.xdgHomeEnvValue)
			t.Setenv("HOME", tc.homeEnvValue)

			if tc.chdir != "" {
				t.Chdir(tc.chdir)
			}

			file, err := config.FindConfigFile()
			require.ErrorIs(t, err, tc.expectedError)
			assert.Equal(t, tc.expected, file)
		})
	}
}
