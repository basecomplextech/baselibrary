package filesys

import (
	"io/ioutil"
	"os"
)

type filesys struct{}

func newFS() *filesys {
	return &filesys{}
}

// Create creates a file in the fstem.
func (fs *filesys) Create(name string) (File, error) {
	f, err := os.Create(name)
	if err != nil {
		return nil, err
	}
	return newFile(f), nil
}

// Exists returns true if the file/directory exists.
func (fs *filesys) Exists(name string) (bool, error) {
	_, err := os.Lstat(name)
	switch {
	case os.IsNotExist(err):
		return false, nil
	case err != nil:
		return false, nil
	}
	return true, nil
}

// MakeDir creates a directory in the fstem.
func (fs *filesys) MakeDir(name string, perm os.FileMode) error {
	return os.Mkdir(name, perm)
}

// MakePath creates a directory path and all parents that does not exist.
func (fs *filesys) MakePath(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

// Open opens the named file for reading, see os.Open.
func (fs *filesys) Open(filename string) (File, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	return newFile(f), nil
}

// OpenFile is the generalized open call, see os.OpenFile.
func (fs *filesys) OpenFile(filename string, flag int, perm os.FileMode) (File, error) {
	f, err := os.OpenFile(filename, flag, perm)
	if err != nil {
		return nil, err
	}
	return newFile(f), nil
}

// Remove removes a file or an empty directory.
func (fs *filesys) Remove(filename string) error {
	return os.Remove(filename)
}

// RemoveAll removes a path and any children it contains.
func (fs *filesys) RemoveAll(filename string) error {
	return os.RemoveAll(filename)
}

// Rename renames a file or a directory.
func (fs *filesys) Rename(src string, dst string) error {
	return os.Rename(src, dst)
}

// Stat returns a file info.
func (fs *filesys) Stat(filename string) (FileInfo, error) {
	return os.Stat(filename)
}

// MkdirTemp creates a new temporary directory in the directory dir
// and returns the pathname of the new directory.
func (fs *filesys) TempDir(dir, pattern string) (name string, err error) {
	return os.MkdirTemp(dir, pattern)
}

// TempFile creates a new temporary file in the directory dir,
// opens the file for reading and writing, and returns the resulting *os.File.
func (fs *filesys) TempFile(dir, pattern string) (File, error) {
	f, err := ioutil.TempFile(dir, pattern)
	if err != nil {
		return nil, err
	}
	return newFile(f), nil
}
