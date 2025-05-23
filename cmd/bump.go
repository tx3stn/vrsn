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
	"github.com/tx3stn/vrsn/internal/prompt"
	"github.com/tx3stn/vrsn/internal/version"
)

// NewCmdBump creates the bump command.
// TODO: split this out into smaller chunks and remove nolint.
//
//nolint:funlen,cyclop
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

			conf, err := config.Get()
			if err != nil {
				return fmt.Errorf("error getting config: %w", err)
			}

			log.Debugf("config: %+v", conf)

			log.Debugf("bump command args: %s", args)

			if flags.GitTag {
				return bumpGitTag(curDir, args, log)
			}

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

			newVersion, err := getNewVersion(currentVersion, args)
			if err != nil {
				return err
			}

			if err := files.WriteVersionToFile(curDir, versionFile, newVersion); err != nil {
				return fmt.Errorf("error writing version to file: %w", err)
			}

			log.Infof("version bumped from %s to %s", currentVersion, newVersion)

			if conf.Commit {
				addOutput, err := git.Add(curDir, versionFile)
				if err != nil {
					log.Infof("git add output: %s", addOutput)

					return fmt.Errorf("error git adding files: %w", err)
				}

				commitOutput, err := git.Commit(curDir, versionFile, conf.CommitMsg)
				if err != nil {
					log.Infof("git commit output: %s", commitOutput)

					return fmt.Errorf("error git committing files: %w", err)
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

	cmd.Flags().Bool("commit", false, "Commit the updated version file after bumping.")

	if err := viper.BindPFlag("commit", cmd.Flags().Lookup("commit")); err != nil {
		fmt.Printf("error binding commit flag: %s", err)
		os.Exit(1)
	}

	cmd.Flags().
		String(
			"commit-msg",
			"bump version",
			"Customise the commit message used when committing the version bump.",
		)

	if err := viper.BindPFlag("commit-msg", cmd.Flags().Lookup("commit-msg")); err != nil {
		fmt.Printf("error binding commit-msg flag: %s", err)
		os.Exit(1)
	}

	cmd.Flags().
		BoolVar(&flags.GitTag, "git-tag", false, "Use git tags rather than a version file.")
	cmd.Flags().
		StringVar(
			&flags.TagMsg,
			"tag-msg",
			"",
			"Customise the tag message used when adding the version tag.",
		)

	return cmd
}

func getNewVersion(currentVersion string, args []string) (string, error) {
	var newVersion string

	var err error

	if len(args) > 0 {
		options, err := version.GetBumpOptions(currentVersion)
		if err != nil {
			return "", fmt.Errorf("error getting bump options: %w", err)
		}

		newVersion, err = options.SelectedIncrement(args[0])
		if err != nil {
			return "", fmt.Errorf("error getting selected increment: %w", err)
		}
	} else {
		bump := prompt.NewBumpSelector()

		newVersion, err = bump.Select(currentVersion)
		if err != nil {
			return "", fmt.Errorf("error selecting bump type: %w", err)
		}
	}

	return newVersion, nil
}

func bumpGitTag(curDir string, args []string, log logger.Basic) error {
	currentVersion, err := git.LatestTag(curDir)
	if err != nil {
		return fmt.Errorf("error getting latest tag: %w", err)
	}

	log.Debugf("current git tag version: %s", currentVersion)

	newVersion, err := getNewVersion(currentVersion, args)
	if err != nil {
		return err
	}

	if flags.TagMsg == "" {
		flags.TagMsg = "Release " + newVersion
	}

	if err := git.AddTag(curDir, newVersion, flags.TagMsg); err != nil {
		return fmt.Errorf("error adding tag: %w", err)
	}

	log.Debugf("new git tag added: %s", newVersion)

	return nil
}
