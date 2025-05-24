package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tx3stn/vrsn/internal/config"
	"github.com/tx3stn/vrsn/internal/files"
	"github.com/tx3stn/vrsn/internal/flags"
	"github.com/tx3stn/vrsn/internal/git"
	"github.com/tx3stn/vrsn/internal/logger"
	"github.com/tx3stn/vrsn/internal/version"
)

// NewCmdCheck creates the check command.
//
//nolint:gocognit,cyclop,funlen
func NewCmdCheck() *cobra.Command {
	shortDescription := "Check the semantic version has been correctly incremented."

	cmd := &cobra.Command{
		RunE: func(ccmd *cobra.Command, args []string) error {
			// TODO: support color option.
			conf, err := config.Get()
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

			versionFileFinder := files.VersionFileFinder{
				ErrorOnNoFilesFound: false,
				FileFlag:            flags.VersionFile,
				Logger:              log,
				SearchDir:           curDir,
			}

			versionFile, err := versionFileFinder.Find()
			if err != nil {
				return fmt.Errorf("error locating version file: %w", err)
			}

			if flags.Now == "" {
				if versionFile == "" {
					log.Info("no version files found in directory and no --now flag provided")

					return ErrNoNowOrFile
				}

				log.Debugf("reading current version from %s", versionFile)

				flags.Now, err = files.GetVersionFromFile(curDir, versionFile)
				if err != nil {
					return fmt.Errorf("error reading version from file: %w", err)
				}
			}

			if flags.Was == "" {
				if versionFile == "" {
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

				log.Debugf(
					"reading previous version from %s on branch %s",
					versionFile,
					conf.Check.BaseBranch,
				)

				baseBranchVersion, err := git.VersionAtBranch(
					curDir,
					conf.Check.BaseBranch,
					versionFile,
				)
				if err != nil {
					return fmt.Errorf("error getting version at branch: %w", err)
				}

				flags.Was, err = files.GetVersionFromString(versionFile, baseBranchVersion)
				if err != nil {
					return fmt.Errorf("error parsing the version from string: %w", err)
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
		String(
			"base-branch",
			"main",
			"Name of the base branch used when auto detecting version changes.",
		)

	if err := viper.BindPFlag("base-branch", cmd.Flags().Lookup("base-branch")); err != nil {
		fmt.Printf("error binding --base-branch flag: %s", err)
		os.Exit(1)
	}

	cmd.Flags().
		StringVar(&flags.Was, "was", "", "The previous semantic version (if passing for direct comparison).")
	cmd.PersistentFlags().
		StringVar(&flags.Now, "now", "", "The current semantic version (if passing for direct comparison).")

	return cmd
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
