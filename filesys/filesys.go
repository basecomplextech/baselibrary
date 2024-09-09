// Copyright 2022 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package filesys

import (
	"os"

	"github.com/basecomplextech/baselibrary/system"
)

type FileSystem interface {
	// Create creates a file in the file system.
	Create(name string) (File, error)

	// Exists returns true if the file/directory exists.
	Exists(name string) (bool, error)

	// MakeDir creates a directory in the file system.
	MakeDir(name string, perm os.FileMode) error

	// MakePath creates a directory path and all parents that does not exist.
	MakePath(path string, perm os.FileMode) error

	// Open opens the named file for reading, see os.Open.
	Open(filename string) (File, error)

	// OpenFile is the generalized open call, see os.OpenFile.
	OpenFile(filename string, flag int, perm os.FileMode) (File, error)

	// Remove removes a file or an empty directory.
	Remove(filename string) error

	// RemoveAll removes a path and any children it contains.
	RemoveAll(filename string) error

	// Rename renames a file or a directory.
	Rename(src string, dst string) error

	// Stat returns a file info.
	Stat(filename string) (FileInfo, error)

	// MkdirTemp creates a new temporary directory in the directory dir
	// And returns the pathname of the new directory.
	TempDir(dir, pattern string) (name string, err error)

	// TempFile creates a new temporary file in the directory dir,
	// Opens the file for reading and writing, and returns the resulting *os.File.
	TempFile(dir, pattern string) (File, error)

	// Usage returns a disk usage info of a directory.
	Usage(path string) (system.DiskInfo, error)
}

// New returns a standard file system.
func New() FileSystem {
	return newFS()
}

// Default returns the default file system.
func Default() FileSystem {
	return defaultFS
}

// private

var defaultFS = newFS()
