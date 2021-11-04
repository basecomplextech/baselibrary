package memfs

import (
	"os"

	"github.com/complex1io/library/fs"
	"github.com/spf13/afero"
)

var _ fs.FileSystem = (*memfs)(nil)

// New returns a new in-memory fstem.
func New() fs.FileSystem {
	return newMemFS()
}

type memfs struct {
	fs afero.Fs
}

func newMemFS() *memfs {
	fs := afero.NewMemMapFs()
	return &memfs{fs: fs}
}

// Create creates a file in the fstem.
func (fs *memfs) Create(name string) (fs.File, error) {
	f, err := fs.fs.Create(name)
	if err != nil {
		return nil, err
	}
	return newFile(f), nil
}

// Mkdir creates a directory in the fstem.
func (fs *memfs) Mkdir(name string, perm os.FileMode) error {
	return fs.fs.Mkdir(name, perm)
}

// MkdirAll creates a directory path and all parents that does not exist.
func (fs *memfs) MkdirAll(path string, perm os.FileMode) error {
	return fs.fs.MkdirAll(path, perm)
}

// Open opens the named file for reading, see os.Open.
func (fs *memfs) Open(filename string) (fs.File, error) {
	f, err := fs.fs.Open(filename)
	if err != nil {
		return nil, err
	}
	return newFile(f), nil
}

// OpenFile is the generalized open call, see os.OpenFile.
func (fs *memfs) OpenFile(filename string, flag int, perm os.FileMode) (fs.File, error) {
	f, err := fs.fs.OpenFile(filename, flag, perm)
	if err != nil {
		return nil, err
	}
	return newFile(f), nil
}

// Remove removes a file or an empty directory.
func (fs *memfs) Remove(filename string) error {
	return fs.fs.Remove(filename)
}

// Rename renames a file or a directory.
func (fs *memfs) Rename(src string, dst string) error {
	return fs.fs.Rename(src, dst)
}

// Stat returns a file info.
func (fs *memfs) Stat(filename string) (fs.FileInfo, error) {
	return fs.fs.Stat(filename)
}

// Temp creates a new temporary file in the directory dir, see ioutil.TempFile.
func (fs *memfs) TempFile(dir, pattern string) (fs.File, error) {
	f, err := afero.TempFile(fs.fs, dir, pattern)
	if err != nil {
		return nil, err
	}
	return newFile(f), nil
}
