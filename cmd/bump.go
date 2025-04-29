package cmd

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/tx3stn/vrsn/internal/files"
	"github.com/tx3stn/vrsn/internal/flags"
	"github.com/tx3stn/vrsn/internal/git"
	"github.com/tx3stn/vrsn/internal/logger"
	"github.com/tx3stn/vrsn/internal/prompt"
	"github.com/tx3stn/vrsn/internal/version"
)

// NewCmdBump creates the bump command.
// TODO: split this out into smaller chunks and remove nolint.
//
//nolint:cyclop,funlen
func NewCmdBump() *cobra.Command {
	shortDescription := "Increment the current semantic version with a valid patch, major or minor bump."

	cmd := &cobra.Command{
		Args: cobra.OnlyValidArgs,
		RunE: func(ccmd *cobra.Command, args []string) error {
			// TODO: support color option.
			log := logger.NewBasic(false, flags.Verbose)
			curDir, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("error getting current working directory: %w", err)
			}

			log.Debugf("bump command args: %s", args)

			versionFileFinder := files.VersionFileFinder{
				FileFlag:  flags.VersionFile,
				Logger:    log,
				SearchDir: curDir,
			}

			versionFile, err := versionFileFinder.Find()
			if err != nil {
				return fmt.Errorf("error finding version file: %w", err)
			}

			currentVersion, err := files.GetVersionFromFile(curDir, versionFile)
			if err != nil {
				return fmt.Errorf("errpr getting version from file: %w", err)
			}

			var newVersion string
			if len(args) > 0 {
				options, err := version.GetBumpOptions(currentVersion)
				if err != nil {
					return fmt.Errorf("error getting bump options: %w", err)
				}

				newVersion, err = options.SelectedIncrement(args[0])
				if err != nil {
					return fmt.Errorf("error getting selected increment: %w", err)
				}
			} else {
				newVersion, err = prompt.SelectBumpType(currentVersion)
				if err != nil {
					return fmt.Errorf("error selecting bump type: %w", err)
				}
			}

			if err := files.WriteVersionToFile(curDir, versionFile, newVersion); err != nil {
				return fmt.Errorf("error writing version to file: %w", err)
			}

			log.Infof("version bumped from %s to %s", currentVersion, newVersion)

			if flags.Commit {
				addOutput, err := git.Add(curDir, versionFile)
				if err != nil {
					return errors.Wrapf(err, "git add output: %s", addOutput)
				}

				commitOutput, err := git.Commit(curDir, versionFile, flags.CommitMsg)
				if err != nil {
					return errors.Wrapf(err, "git add output: %s", commitOutput)
				}

				log.Infof("version file committed")
			}

			return nil
		},
		//nolint:perfsprint
		Long: fmt.Sprintf(`%s

Pass the increment type directly as an argument to the command, e.g.:

  vrsn bump patch

Or use the interactive prompt to select the increment you want.
The semantic version in the version file will be updated in place.`, shortDescription),
		Short:         shortDescription,
		SilenceErrors: true,
		SilenceUsage:  true,
		Use:           "bump",
		ValidArgs:     []string{"patch", "major", "minor"},
	}

	cmd.Flags().
		BoolVar(&flags.Commit, "commit", false, "Commit the updated version file after bumping.")
	cmd.Flags().
		StringVar(
			&flags.CommitMsg,
			"commit-msg",
			"bump version",
			"Customise the commit message used when committing the version bump.",
		)

	return cmd
}
