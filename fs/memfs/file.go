package memfs

import (
	"path/filepath"

	"github.com/complex1io/library/fs"
	"github.com/spf13/afero"
	"github.com/spf13/afero/mem"
)

var _ fs.File = (*file)(nil)

type file struct {
	*mem.File
}

func newFile(f afero.File) *file {
	mf := f.(*mem.File)
	return &file{File: mf}
}

// Filename returns a file name, not a path as in *os.File.
func (f *file) Filename() string {
	path := f.Path()
	_, name := filepath.Split(path)
	return name
}

// Ext returns a file extension.
func (f *file) Ext() string {
	path := f.Path()
	return filepath.Ext(path)
}

// Path returns *os.File.Name() as presented to Open.
func (f *file) Path() string {
	return f.File.Name()
}

// Map maps the file into memory and returns its data.
func (f *file) Map() ([]byte, error) {
	// TODO: Map file
	// return f.File.Map()
	panic("not implemented")
}

// Size returns the file size in bytes.
func (f *file) Size() (int64, error) {
	stat, err := f.Stat()
	if err != nil {
		return 0, err
	}

	size := stat.Size()
	return size, nil
}
