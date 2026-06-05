package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tx3stn/vrsn/internal/config"
	"github.com/tx3stn/vrsn/internal/files"
	"github.com/tx3stn/vrsn/internal/flags"
	"github.com/tx3stn/vrsn/internal/git"
	"github.com/tx3stn/vrsn/internal/logger"
)

// NewCmdGet creates the get command.
func NewCmdGet() *cobra.Command {
	shortDescription := "Get the current semantic version."

	cmd := &cobra.Command{
		RunE: runGet,
		//nolint:perfsprint
		Long: fmt.Sprintf(`%s

Prints the current version so it can easily be used in scripts or checked for
info, e.g.:

  version=$(vrsn get)

By default the version is read from the version file in the current directory,
use the --file flag to read from a specific file, or the --git-tag flag to read
the latest git tag.

When multiple files are configured with the files option in the config file the
version found in each file is printed on a separate line.`, shortDescription),
		Short:         shortDescription,
		SilenceErrors: true,
		SilenceUsage:  true,
		Use:           "get",
	}

	cmd.Flags().
		BoolVar(
			&flags.GitTag,
			"git-tag",
			false,
			"Read the current version from the latest git tag rather than a version file.",
		)

	return cmd
}

// runGet is the entrypoint for the get command.
func runGet(ccmd *cobra.Command, args []string) error {
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
	log.Debugf("get command args: %s", args)

	if conf.Bump.GitTag {
		tag, err := git.LatestTag(curDir)
		if err != nil {
			return fmt.Errorf("error getting latest tag: %w", err)
		}

		log.Info(tag)

		return nil
	}

	versionFiles, err := resolveVersionFiles(curDir, conf.Files, log, true)
	if err != nil {
		return fmt.Errorf("error locating version file: %w", err)
	}

	return printVersionsInFiles(curDir, versionFiles, log)
}

// printVersionsInFiles prints the version found in the version files.
// A single file prints the bare version so it can easily be used in scripts,
// multiple files print a file: version line per file.
func printVersionsInFiles(curDir string, versionFiles []string, log logger.Basic) error {
	for _, versionFile := range versionFiles {
		version, err := files.GetVersionFromFile(curDir, versionFile)
		if err != nil {
			return fmt.Errorf("error getting version from file %s: %w", versionFile, err)
		}

		if len(versionFiles) == 1 {
			log.Info(version)

			return nil
		}

		log.Infof("%s: %s", versionFile, version)
	}

	return nil
}
