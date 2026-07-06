package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/tx3stn/vrsn/internal/config"
	"github.com/tx3stn/vrsn/internal/files"
	"github.com/tx3stn/vrsn/internal/flags"
	"github.com/tx3stn/vrsn/internal/git"
	"github.com/tx3stn/vrsn/internal/logger"
	"github.com/tx3stn/vrsn/internal/prompt"
	"github.com/tx3stn/vrsn/internal/template"
	"github.com/tx3stn/vrsn/internal/version"
)

// NewCmdBump creates the bump command.
func NewCmdBump() *cobra.Command {
	shortDescription := "Increment the current semantic version with a valid patch, major or minor bump."

	cmd := &cobra.Command{
		Args: cobra.MatchAll(cobra.MaximumNArgs(1), cobra.OnlyValidArgs),
		RunE: runBump,
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
		BoolVar(
			&flags.AndroidVersionCode,
			"android-version-code",
			false,
			"Also bump android:versionCode in AndroidManifest files, derived from the new "+
				"version as MAJOR*10000+MINOR*100+PATCH.",
		)

	cmd.Flags().
		BoolVar(&flags.Commit, "commit", false, "Commit the updated version file after bumping.")

	cmd.Flags().
		StringVar(
			&flags.CommitMsg,
			"commit-msg",
			"bump version",
			"Customise the commit message used when committing the version bump. "+
				"Supports Go template syntax with the {{.Version}} variable for the new version.",
		)

	cmd.Flags().
		BoolVar(
			&flags.GitTag,
			"git-tag",
			false,
			"Use git tags rather than a version file. "+
				"Combine with --file and --commit to bump the version file, commit it and tag the commit.",
		)

	cmd.Flags().
		StringVar(
			&flags.TagMsg,
			"tag-msg",
			"",
			"Customise the tag message used when adding the version tag. "+
				"Supports Go template syntax with the {{.Version}} variable for the new version.",
		)

	return cmd
}

// runBump is the entrypoint for the bump command.
func runBump(ccmd *cobra.Command, args []string) error {
	// TODO: support color option.
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
	log.Debugf("bump command args: %s", args)

	err = ValidateBumpOpts(conf.Bump.GitTag, conf.Files, conf.Bump.Commit)
	if err != nil {
		return err
	}

	if conf.Bump.GitTag && len(conf.Files) == 0 {
		return bumpGitTag(curDir, args, log, conf.Bump.TagMsg)
	}

	newVersion, err := bumpVersionFile(curDir, args, log, conf)
	if err != nil {
		return err
	}

	// When --git-tag is combined with --file the version file has been
	// bumped and committed above, so the tag points at the bump commit.
	if conf.Bump.GitTag {
		if err := applyGitTag(curDir, newVersion, conf.Bump.TagMsg); err != nil {
			return err
		}

		log.Infof("git tag %s added", newVersion)
	}

	return nil
}

// ValidateBumpOpts checks the combination of bump options is valid.
func ValidateBumpOpts(gitTag bool, versionFiles []string, commit bool) error {
	if gitTag && len(versionFiles) > 0 && !commit {
		return ErrGitTagFileNoCommit
	}

	return nil
}

// bumpVersionFile finds the version files, bumps the version in them and
// optionally commits the change, returning the new version.
func bumpVersionFile(
	curDir string,
	args []string,
	log logger.Basic,
	conf config.Config,
) (string, error) {
	versionFiles, err := resolveVersionFiles(curDir, conf.Files, log, true)
	if err != nil {
		return "", err
	}

	currentVersion, err := files.GetVersionsFromFiles(curDir, versionFiles, log)
	if err != nil {
		return "", fmt.Errorf("error getting version from files: %w", err)
	}

	newVersion, err := getNewVersion(currentVersion, args)
	if err != nil {
		return "", err
	}

	// Render the commit message before writing so an invalid template errors
	// before any files are changed.
	commitMsg := ""
	if conf.Bump.Commit {
		commitMsg, err = template.Render(conf.Bump.CommitMsg, newVersion)
		if err != nil {
			return "", fmt.Errorf("error rendering commit message: %w", err)
		}
	}

	writeOpts := files.WriteOptions{NewVersion: newVersion}

	// The version code is derived from the new semver, so it is computed once
	// and only when requested, then applied to any AndroidManifest files.
	if conf.Bump.AndroidVersionCode {
		parsed, parseErr := version.Parse(newVersion)
		if parseErr != nil {
			return "", fmt.Errorf(
				"error parsing new version for android version code: %w",
				parseErr,
			)
		}

		writeOpts.AndroidVersionCode = strconv.Itoa(parsed.AndroidVersionCode())
	}

	for _, versionFile := range versionFiles {
		if err := files.WriteVersionToFile(curDir, versionFile, writeOpts); err != nil {
			return "", fmt.Errorf("error writing version to file %s: %w", versionFile, err)
		}

		log.Debugf("bumped version in %s", versionFile)
	}

	log.Infof("version bumped from %s to %s", currentVersion, newVersion)

	if conf.Bump.Commit {
		if err := commitVersionFiles(curDir, versionFiles, commitMsg, log); err != nil {
			return "", err
		}
	}

	return newVersion, nil
}

// commitVersionFiles stages the bumped version files and commits them all in
// a single commit.
func commitVersionFiles(
	curDir string,
	versionFiles []string,
	commitMsg string,
	log logger.Basic,
) error {
	addOutput, err := git.Add(curDir, versionFiles...)
	if err != nil {
		log.Infof("git add output: %s", addOutput)

		return fmt.Errorf("error git adding files: %w", err)
	}

	commitOutput, err := git.Commit(curDir, commitMsg, versionFiles...)
	if err != nil {
		log.Infof("git commit output: %s", commitOutput)

		return fmt.Errorf("error git committing files: %w", err)
	}

	log.Infof("version file committed")

	return nil
}

// applyGitTag adds the new version as an annotated git tag, defaulting the tag
// message when one isn't provided.
func applyGitTag(curDir string, newVersion string, tagMsg string) error {
	if tagMsg == "" {
		tagMsg = "Release " + newVersion
	}

	renderedMsg, err := template.Render(tagMsg, newVersion)
	if err != nil {
		return fmt.Errorf("error rendering tag message: %w", err)
	}

	if err := git.AddTag(curDir, newVersion, renderedMsg); err != nil {
		return fmt.Errorf("error adding tag: %w", err)
	}

	return nil
}

func getNewVersion(currentVersion string, args []string) (string, error) {
	if len(args) > 0 {
		options, err := version.GetBumpOptions(currentVersion)
		if err != nil {
			return "", fmt.Errorf("error getting bump options: %w", err)
		}

		newVersion, err := options.SelectedIncrement(args[0])
		if err != nil {
			return "", fmt.Errorf("error getting selected increment: %w", err)
		}

		return newVersion, nil
	}

	newVersion, err := prompt.NewBumpSelector().Select(currentVersion)
	if err != nil {
		return "", fmt.Errorf("error selecting bump type: %w", err)
	}

	return newVersion, nil
}

func bumpGitTag(curDir string, args []string, log logger.Basic, tagMsg string) error {
	currentVersion, err := git.LatestTag(curDir)
	if err != nil {
		return fmt.Errorf("error getting latest tag: %w", err)
	}

	log.Debugf("current git tag version: %s", currentVersion)

	newVersion, err := getNewVersion(currentVersion, args)
	if err != nil {
		return err
	}

	if err := applyGitTag(curDir, newVersion, tagMsg); err != nil {
		return err
	}

	log.Infof("git tag version bumped from %s to %s", currentVersion, newVersion)

	return nil
}
