package terminal

import (
	"os"

	"github.com/mattn/go-isatty"
)

// FileDescriptor is the interface that wraps the file descriptor method.
type FileDescriptor interface {
	Fd() uintptr
}

// CheckColor checks that the file descriptor is a terminal and it supports basic color.
// If the environment variables force no color, then returns false.
func CheckColor(fd FileDescriptor) bool {
	if !isatty.IsTerminal(fd.Fd()) {
		return false
	}
	return SupportsColor()
}

// SupportsColor returns true if the environment supports basic color.
func SupportsColor() bool {
	// check disabled
	switch {
	case hasFlag("no-color"):
		return false
	case hasFlag("no-colors"):
		return false
	case hasFlag("color=false"):
		return false
	case hasFlag("color=never"):
		return false
	}

	// check color
	switch {
	case checkForceColor():
		return true
	case hasFlag("color"):
		return true
	case hasFlag("colors"):
		return true
	case hasFlag("color=true"):
		return true
	case hasFlag("color=always"):
		return true
	case hasFlag("color=16m"):
		return true
	case hasFlag("color=full"):
		return true
	case hasFlag("color=256"):
		return true
	case hasFlag("color=truecolor"):
		return true
	}
	return false
}

func checkForceColor() bool {
	_, ok := os.LookupEnv("FORCE_COLOR")
	if !ok {
		return false
	}

	s := os.Getenv("FORCE_COLOR")
	if s == "false" {
		return false
	}
	return true
}
