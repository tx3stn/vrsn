// Package logger implements a basic logger to support the core level of output
// and a verbose mode, and maybe some color in future.
package logger

import "fmt"

// Basic is a basic logger to provide the simple functionality to allow people
// to run vrsn in verbose mode and disable/enable color support.
type Basic struct {
	UseColor bool
	Verbose  bool
}

// NewBasic creates an instance of the Basic logger.
func NewBasic(color bool, verbose bool) Basic {
	return Basic{
		UseColor: color,
		Verbose:  verbose,
	}
}

// Debug is a log that will only be displayed if the `vrsn` command is run in
// verbose mode.
func (b Basic) Debug(msg string) {
	if b.Verbose {
		fmt.Println(msg)
	}
}

// Debugf is a log that will only be displayed if the `vrsn` command is run in
// verbose mode, with support for variables.
func (b Basic) Debugf(msg string, args ...interface{}) {
	b.Debug(fmt.Sprintf(msg, args...))
}

// Info is an info level log which will be default always be displayed.
func (b Basic) Info(msg string) {
	fmt.Println(msg)
}

// Infof is an info level log which will be default always be displayed with
// support for variables.
func (b Basic) Infof(msg string, args ...interface{}) {
	b.Info(fmt.Sprintf(msg, args...))
}
