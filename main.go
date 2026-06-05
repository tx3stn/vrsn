// Package main is the entrypoint of the CLI.
package main

import (
	"fmt"
	"os"

	"github.com/tx3stn/vrsn/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
