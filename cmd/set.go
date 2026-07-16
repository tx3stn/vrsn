package cmd

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tx3stn/vrsn/internal/config"
	"github.com/tx3stn/vrsn/internal/flags"
	"github.com/tx3stn/vrsn/internal/logger"
	"github.com/tx3stn/vrsn/internal/version"
)

// setSuffixRegex matches the optional suffix set accepts after the first '-'
// (e.g. dev, rc1, fix-this): a non-empty run of letters, digits and hyphens.
//

var setSuffixRegex = regexp.MustCompile(`^[a-zA-Z0-9-]+$`)

// NewCmdSet creates the set command.
func NewCmdSet() *cobra.Command {
	shortDescription := "Set the semantic version in the version file(s) directly."

	cmd := &cobra.Command{
		Args: cobra.ExactArgs(1),
		RunE: runSet,
		//nolint:perfsprint
		Long: fmt.Sprintf(`%s

Pass the version to write as an argument, e.g.:

  vrsn set 2.0.0

Unlike bump, set does not check that the version is a valid increment of the
current one, so it can jump to an arbitrary version or even move backwards. The
version must be MAJOR.MINOR.PATCH, optionally with a "-" suffix of letters,
digits and hyphens (e.g. 1.2.3-dev). It only updates the version file(s); it
does not commit or tag.`, shortDescription),
		Short:         shortDescription,
		SilenceErrors: true,
		SilenceUsage:  true,
		Use:           "set <version>",
	}

	cmd.Flags().
		BoolVar(
			&flags.AndroidVersionCode,
			"android-version-code",
			false,
			"Also set android:versionCode in AndroidManifest files, derived from the "+
				"version as MAJOR*10000+MINOR*100+PATCH.",
		)

	return cmd
}

// runSet is the entrypoint for the set command.
func runSet(ccmd *cobra.Command, args []string) error {
	conf, err := config.Get(flags.ConfigFile, ccmd.Flags())
	if err != nil {
		return fmt.Errorf("error getting config: %w", err)
	}

	log := logger.NewBasic(false, conf.Verbose)

	curDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting current working directory: %w", err)
	}

	log.Debugf("config: %+v", conf)
	log.Debugf("set command args: %s", args)

	_, err = writeVersion(curDir, args, log, conf, writeConfig{
		resolve:            getSetVersion,
		verb:               "set",
		androidVersionCode: conf.Set.AndroidVersionCode,
	})

	return err
}

// getSetVersion validates the supplied version and returns its canonical form.
// It ignores the current version so, unlike bump, it performs no
// increment-validity check. It additionally accepts an optional "-<suffix>"
// marker (e.g. 1.2.3-dev) that the numeric-only bump and check do not.
func getSetVersion(_ string, args []string) (string, error) {
	core, suffix, hasSuffix := strings.Cut(args[0], "-")

	parsed, err := version.Parse(core)
	if err != nil {
		return "", fmt.Errorf("error parsing version: %w", err)
	}

	result := parsed.String()

	if hasSuffix {
		if !setSuffixRegex.MatchString(suffix) {
			return "", ErrInvalidVersionSuffix
		}

		result += "-" + suffix
	}

	return result, nil
}
