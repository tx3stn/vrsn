// Package cmd contains all of the CLI commands.
package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tx3stn/vrsn/internal/flags"
)

// Version is the CLI version set via linker flags at build time.
//
//nolint:gochecknoglobals
var Version string

//nolint:gochecknoglobals
var rootCmd = &cobra.Command{
	RunE: func(ccmd *cobra.Command, args []string) error {
		err := ccmd.Help()
		if err != nil {
			return fmt.Errorf("error getting cobra help: %w", err)
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

	//nolint:wrapcheck
	return rootCmd.ExecuteContext(ctx)
}

//nolint:gochecknoinits
func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.AddCommand(NewCmdCheck())
	rootCmd.AddCommand(NewCmdBump())

	rootCmd.PersistentFlags().
		Bool("verbose", false, "display verbose output for more detail on what the command is doing")

	if err := viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose")); err != nil {
		fmt.Printf("error binding --verbose flag: %s", err)
		os.Exit(1)
	}

	rootCmd.PersistentFlags().
		StringVar(&flags.VersionFile, "file", "", "specify the path to the version file (if not in current directory)")

	rootCmd.PersistentFlags().
		StringVar(&flags.ConfigFile, "config", "", "override the config file location")

	if err := viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config")); err != nil {
		fmt.Printf("error binding --config flag: %s", err)
		os.Exit(1)
	}
}

func initConfig() {
	if flags.ConfigFile == "" {
		viper.SetConfigName("vrsn")
		viper.SetConfigType("toml")
		viper.AddConfigPath("$XDG_CONFIG_DIR/")
		viper.AddConfigPath("$HOME/.config")
	} else {
		viper.SetConfigFile(flags.ConfigFile)
	}

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("error trying to read config file: %s", err)
		os.Exit(1)
	}
}
