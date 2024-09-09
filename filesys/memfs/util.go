// Copyright 2022 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package memfs

import (
	"path/filepath"
	"strings"
)

func splitPath(path string) []string {
	if path == "" {
		return []string{}
	}

	path = cleanPath(path)
	return strings.Split(path, string(filepath.Separator))
}

func cleanPath(path string) string {
	path = filepath.Clean(path)
	path = strings.TrimLeft(path, ".")
	path = strings.TrimLeft(path, string(filepath.Separator))
	path = strings.TrimRight(path, string(filepath.Separator))
	return path
}
