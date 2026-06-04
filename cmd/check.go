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
	"github.com/tx3stn/vrsn/internal/version"
)

// NewCmdCheck creates the check command.
//
//nolint:cyclop,funlen
func NewCmdCheck() *cobra.Command {
	shortDescription := "Check the semantic version has been correctly incremented."

	cmd := &cobra.Command{
		RunE: func(ccmd *cobra.Command, args []string) error {
			// TODO: support color option.
			conf, err := config.Get(flags.ConfigFile)
			if err != nil {
				return fmt.Errorf("error getting config: %w", err)
			}

			log := logger.NewBasic(false, conf.Verbose)

			curDir, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("error getting current working directory: %w", err)
			}

			log.Debugf("config: %+v", conf)
			log.Debugf("check command args: %s", args)

			if flags.Was != "" && flags.Now != "" {
				return validateAndCompare(log, flags.Was, flags.Now)
			}

			currentBranch, err := git.CurrentBranch(curDir)
			if err != nil {
				return fmt.Errorf("error getting current git branch: %w", err)
			}

			log.Debugf("current branch: %s", currentBranch)

			versionFiles, err := resolveVersionFiles(curDir, conf.Files, log, false)
			if err != nil {
				return fmt.Errorf("error locating version file: %w", err)
			}

			if flags.Now == "" {
				if len(versionFiles) == 0 {
					log.Info("no version files found in directory and no --now flag provided")

					return ErrNoNowOrFile
				}

				flags.Now, err = files.GetVersionsFromFiles(curDir, versionFiles, log)
				if err != nil {
					return fmt.Errorf("error reading version from files: %w", err)
				}
			}

			if flags.Was == "" {
				if len(versionFiles) == 0 {
					log.Info("no version files found in directory and no --was flag provided")

					return ErrNoWasOrFile
				}

				if currentBranch == conf.Check.BaseBranch {
					return fmt.Errorf(
						"%w: base branch: %s",
						ErrCantCompareVersionsOnBranch,
						conf.Check.BaseBranch,
					)
				}

				flags.Was, err = getWasVersionFromFiles(
					curDir,
					conf.Check.BaseBranch,
					versionFiles,
					log,
				)
				if err != nil {
					return err
				}
			}

			return validateAndCompare(log, flags.Was, flags.Now)
		},
		//nolint:perfsprint
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
	cmd.Flags().
		StringVar(
			&flags.BaseBranch,
			"base-branch",
			"main",
			"Name of the base branch used when auto detecting version changes.",
		)

	cmd.Flags().
		StringVar(&flags.Was, "was", "", "The previous semantic version (if passing for direct comparison).")
	cmd.PersistentFlags().
		StringVar(&flags.Now, "now", "", "The current semantic version (if passing for direct comparison).")

	return cmd
}

// getWasVersionFromFiles reads the version each of the files contained at the
// base branch and returns the common version they all had.
// The version found in each file is debug logged, and if the versions do not
// all match an ErrVersionsDoNotMatch error is returned.
func getWasVersionFromFiles(
	curDir string,
	baseBranch string,
	versionFiles []string,
	log logger.Basic,
) (string, error) {
	versions := make([]string, 0, len(versionFiles))

	for _, versionFile := range versionFiles {
		baseBranchVersion, err := git.VersionAtBranch(curDir, baseBranch, versionFile)
		if err != nil {
			return "", fmt.Errorf("error getting version at branch: %w", err)
		}

		was, err := files.GetVersionFromString(versionFile, baseBranchVersion)
		if err != nil {
			return "", fmt.Errorf("error parsing the version from string: %w", err)
		}

		log.Debugf("file %s has version %s on branch %s", versionFile, was, baseBranch)

		versions = append(versions, was)
	}

	for _, version := range versions {
		if version != versions[0] {
			return "", files.ErrVersionsDoNotMatch
		}
	}

	return versions[0], nil
}

func validateAndCompare(log logger.Basic, was string, now string) error {
	if err := flags.Validate(was, now); err != nil {
		return fmt.Errorf("error validating flags: %w", err)
	}

	log.Infof("was: %s", was)
	log.Infof("now: %s", now)

	if err := version.Compare(was, now); err != nil {
		return fmt.Errorf("error comparing versions: %w", err)
	}

	log.Info("valid version bump")

	return nil
}
