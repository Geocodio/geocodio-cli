// Package ui provides terminal UI components for the Geocodio CLI.
package ui

import (
	"os"

	"golang.org/x/term"
)

// IsTTY returns true if the given file is connected to a terminal.
func IsTTY(f *os.File) bool {
	return term.IsTerminal(int(f.Fd()))
}

// StdoutIsTTY returns true if stdout is connected to a terminal.
func StdoutIsTTY() bool {
	return IsTTY(os.Stdout)
}

// StderrIsTTY returns true if stderr is connected to a terminal.
func StderrIsTTY() bool {
	return IsTTY(os.Stderr)
}

// ColorEnabled returns true if colored output should be used.
// It respects the NO_COLOR and FORCE_COLOR environment variables.
// NO_COLOR takes precedence over FORCE_COLOR.
// If neither is set, it falls back to checking if stdout is a TTY.
func ColorEnabled() bool {
	if _, ok := os.LookupEnv("NO_COLOR"); ok {
		return false
	}
	if _, ok := os.LookupEnv("FORCE_COLOR"); ok {
		return true
	}
	return StdoutIsTTY()
}
