package memfs

import (
	"path/filepath"
	"strings"
)

func splitPath(path string) []string {
	if path == "" {
		return []string{}
	}

	path = strings.TrimLeft(path, string(filepath.Separator))
	return strings.Split(path, string(filepath.Separator))
}
