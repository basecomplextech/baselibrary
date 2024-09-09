// Copyright 2023 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package terminal

import (
	"os"
	"regexp"

	"github.com/mattn/go-isatty"
)

// FileDescriptor is the interface that wraps the file descriptor method.
type FileDescriptor interface {
	Fd() uintptr
}

// CheckColor checks that the file descriptor is a terminal and it supports basic color.
// If the environment variables force no color, then returns false.
func CheckColor(fd FileDescriptor) bool {
	tty := isatty.IsTerminal(fd.Fd())
	if !tty {
		return false
	}
	return SupportsColor()
}

// SupportsColor returns true if the environment supports basic color.
func SupportsColor() bool {
	if checkColorDisabled() {
		return false
	}
	if checkColorFlags() {
		return true
	}
	return checkColorEnv()
}

func checkForceColor() bool {
	s, ok := os.LookupEnv("FORCE_COLOR")
	switch {
	case !ok:
		return false
	case s == "false":
		return false
	}
	return true
}

func checkColorDisabled() bool {
	switch {
	case hasFlag("no-color"):
		return true
	case hasFlag("no-colors"):
		return true
	case hasFlag("color=false"):
		return true
	case hasFlag("color=never"):
		return true
	}
	return false
}

func checkColorFlags() bool {
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

func checkColorEnv() bool {
	// Dump terminal
	term := os.Getenv("TERM")
	if term == "dumb" {
		return true
	}

	// Ci servers
	if _, ci := os.LookupEnv("CI"); ci {
		var names = []string{"TRAVIS", "CIRCLECI", "APPVEYOR", "GITLAB_CI", "GITHUB_ACTIONS",
			"BUILDKITE", "DRONE"}
		for _, name := range names {
			_, ok := os.LookupEnv(name)
			if ok {
				return true
			}
		}

		if os.Getenv("CI_NAME") == "codeship" {
			return true
		}

		return false
	}

	// Color terminal
	if _, ok := os.LookupEnv("COLORTERM"); ok {
		return true
	}

	// Terminal program
	term, ok := os.LookupEnv("TERM_PROGRAM")
	if ok {
		switch term {
		case "iTerm.app":
			return true
		case "Apple_Terminal":
			return true
		}
	}

	// Term256
	term, _ = os.LookupEnv("TERM")
	var term256Regex = regexp.MustCompile("(?i)-256(color)?$")
	if term256Regex.MatchString(term) {
		return true
	}

	// Term basic
	var termBasicRegex = regexp.MustCompile("(?i)^screen|^xterm|^vt100|^vt220|^rxvt|color|ansi|cygwin|linux")
	if termBasicRegex.MatchString(term) {
		return true
	}

	return false
}
