package fs

import (
	"io/ioutil"
	"os"
)

var (
	ErrExist    = os.ErrExist
	ErrNotExist = os.ErrNotExist
)

type FileSystem interface {
	// Create creates a file in the file system.
	Create(name string) (File, error)

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
	// and returns the pathname of the new directory.
	TempDir(dir, pattern string) (name string, err error)

	// TempFile creates a new temporary file in the directory dir,
	// opens the file for reading and writing, and returns the resulting *os.File.
	TempFile(dir, pattern string) (File, error)
}

// New returns a disk file system.
func New() FileSystem {
	return newFS()
}

type fs struct{}

func newFS() *fs {
	return &fs{}
}

// Create creates a file in the fstem.
func (fs *fs) Create(name string) (File, error) {
	f, err := os.Create(name)
	if err != nil {
		return nil, err
	}
	return newFile(f), nil
}

// MakeDir creates a directory in the fstem.
func (fs *fs) MakeDir(name string, perm os.FileMode) error {
	return os.Mkdir(name, perm)
}

// MakePath creates a directory path and all parents that does not exist.
func (fs *fs) MakePath(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

// Open opens the named file for reading, see os.Open.
func (fs *fs) Open(filename string) (File, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	return newFile(f), nil
}

// OpenFile is the generalized open call, see os.OpenFile.
func (fs *fs) OpenFile(filename string, flag int, perm os.FileMode) (File, error) {
	f, err := os.OpenFile(filename, flag, perm)
	if err != nil {
		return nil, err
	}
	return newFile(f), nil
}

// Remove removes a file or an empty directory.
func (fs *fs) Remove(filename string) error {
	return os.Remove(filename)
}

// RemoveAll removes a path and any children it contains.
func (fs *fs) RemoveAll(filename string) error {
	return os.RemoveAll(filename)
}

// Rename renames a file or a directory.
func (fs *fs) Rename(src string, dst string) error {
	return os.Rename(src, dst)
}

// Stat returns a file info.
func (fs *fs) Stat(filename string) (FileInfo, error) {
	return os.Stat(filename)
}

// MkdirTemp creates a new temporary directory in the directory dir
// and returns the pathname of the new directory.
func (fs *fs) TempDir(dir, pattern string) (name string, err error) {
	return os.MkdirTemp(dir, pattern)
}

// TempFile creates a new temporary file in the directory dir,
// opens the file for reading and writing, and returns the resulting *os.File.
func (fs *fs) TempFile(dir, pattern string) (File, error) {
	f, err := ioutil.TempFile(dir, pattern)
	if err != nil {
		return nil, err
	}
	return newFile(f), nil
}
