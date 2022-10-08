// Package cmd contains all of the CLI commands.
package cmd

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/tx3stn/vrsn/internal/flags"
)

// Version is the CLI version set via linker flags at build time.
var Version string

var rootCmd = &cobra.Command{
	RunE: func(ccmd *cobra.Command, args []string) error {
		err := ccmd.Help()
		if err != nil {
			return err
		}
		return nil
	},
	Short:   "A single tool for all of your semantic versioning needs.",
	Use:     "vrsn",
	Version: Version,
}

// Execute executes the root command.
func Execute() error {
	ctx := context.Background()
	return rootCmd.ExecuteContext(ctx)
}

func init() {
	rootCmd.AddCommand(NewCmdCheck())
	rootCmd.AddCommand(NewCmdBump())
	rootCmd.PersistentFlags().BoolVar(&flags.Verbose, "verbose", false, "display verbose output for more detail on what the command is doing")
	rootCmd.PersistentFlags().StringVar(&flags.VersionFile, "file", "", "specify the path to the version file (if not in current directory)")
}
