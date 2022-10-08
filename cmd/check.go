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
	"github.com/tx3stn/vrsn/internal/version"
)

// NewCmdCheck creates the check command.
func NewCmdCheck() *cobra.Command {
	shortDescription := "Check the semantic version has been correctly incremented."

	cmd := &cobra.Command{
		RunE: func(ccmd *cobra.Command, args []string) error {
			// TODO: support color option.
			log := logger.NewBasic(false, flags.Verbose)
			curDir, err := os.Getwd()
			if err != nil {
				return err
			}

			if flags.Was != "" && flags.Now != "" {
				return validateAndCompare(log, flags.Was, flags.Now)
			}

			currentBranch, err := git.CurrentBranch(curDir)
			if err != nil {
				return err
			}

			log.Debugf("current branch: %s", currentBranch)

			versionFileFinder := files.VersionFileFinder{
				ErrorOnNoFilesFound: false,
				FileFlag:            flags.VersionFile,
				Logger:              log,
				SearchDir:           curDir,
			}

			versionFile, err := versionFileFinder.Find()
			if err != nil {
				return err
			}

			if flags.Now == "" {
				if versionFile == "" {
					log.Info("no version files found in directory and no --now flag provided")
					return errors.New("please either pass version with --now flag or run inside a directory that uses a version file")
				}

				log.Debugf("reading current version from %s", versionFile)

				flags.Now, err = files.GetVersionFromFile(curDir, versionFile)
				if err != nil {
					return err
				}
			}

			if flags.Was == "" {
				if versionFile == "" {
					log.Info("no version files found in directory and no --was flag provided")
					return errors.New("please either pass version with --was flag or run inside a directory that uses a version file")
				}

				if currentBranch == flags.BaseBranch {
					return errors.Errorf("currently on the %s branch and no --was value supplied, unable to compare versions", flags.BaseBranch)
				}

				log.Debugf("reading previous version from %s on branch %s", versionFile, flags.BaseBranch)

				baseBranchVersion, err := git.VersionAtBranch(curDir, flags.BaseBranch, versionFile)
				if err != nil {
					return err
				}

				flags.Was, err = files.GetVersionFromString(versionFile, baseBranchVersion)
				if err != nil {
					return err
				}
			}

			return validateAndCompare(log, flags.Was, flags.Now)
		},
		Long: fmt.Sprintf(`%s

Detects if you are on a branch that is not the repository's base branch so the
current version can be read from the git history.
If you're on a branch that is not the repository's base branch just run:

  vrsn check

That's all you need!

You can also use the --was and --now flags to compare the versions so you can
read them from A N Y W H E R E.
`, shortDescription),
		Short:         shortDescription,
		SilenceErrors: true,
		SilenceUsage:  true,
		Use:           "check",
	}
	cmd.Flags().StringVar(&flags.BaseBranch, "base-branch", "main", "Name of the base branch used when auto detecting version changes.")
	cmd.Flags().StringVar(&flags.Was, "was", "", "The previous semantic version (if passing for direct comparison).")
	cmd.PersistentFlags().StringVar(&flags.Now, "now", "", "The current semantic version (if passing for direct comparison).")
	return cmd
}

func validateAndCompare(log logger.Basic, was string, now string) error {
	if err := flags.Validate(flags.Was, flags.Now); err != nil {
		return err
	}

	log.Infof("was: %s", flags.Was)
	log.Infof("now: %s", flags.Now)

	if err := version.Compare(flags.Was, flags.Now); err != nil {
		return err
	}

	log.Info("valid version bump")

	return nil
}
