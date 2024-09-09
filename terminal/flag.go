// Copyright 2023 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package terminal

import (
	"os"
	"strings"
)

func hasFlag(flag string) bool {
	args := os.Args
	return hasFlagArgs(args, flag)
}

func hasFlagArgs(args []string, flag string) bool {
	// Prefix the flag with the necessary dashes
	var prefix string
	if !strings.HasPrefix(flag, "-") {
		if len(flag) == 1 {
			prefix = "-"
		} else {
			prefix = "--"
		}
	}

	// Check flag position
	pos := flagIndexOf(args, prefix+flag)
	if pos == -1 {
		return false
	}

	// Check terminator "--" position, stop parsing after it
	term := flagIndexOf(args, "--")
	if term == -1 {
		return true
	}
	return pos < term
}

func flagIndexOf(ss []string, s string) int {
	for i, el := range ss {
		if el == s {
			return i
		}
	}
	return -1
}
