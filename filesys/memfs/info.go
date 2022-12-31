package memfs

import (
	"time"

	"github.com/complex1tech/baselibrary/filesys"
)

var _ filesys.FileInfo = (*memInfo)(nil)

type memInfo struct {
	name string
	size int64
	dir  bool
}

// Name is a base name of the file.
func (i *memInfo) Name() string {
	return i.name
}

// Size is a length in bytes for regular files; system-dependent for others
func (i *memInfo) Size() int64 {
	return i.size
}

// File mode bits.
func (i *memInfo) Mode() filesys.FileMode {
	return 0
}

// Modification time.
func (i *memInfo) ModTime() time.Time {
	return time.Time{}
}

// Abbreviation for Mode().IsDir().
func (i *memInfo) IsDir() bool {
	return i.dir
}

// Underlying data source (can return nil)
func (i *memInfo) Sys() any {
	panic("not implemented")
}
