package filesys

import (
	"errors"
	"os"
	"path/filepath"
	"sync"
	"syscall"

	"github.com/edsrzf/mmap-go"
)

var _ File = (*file)(nil)

type file struct {
	*os.File

	mu   sync.Mutex
	mmap mmap.MMap
}

// newFile wraps a file into a fstem file wrapper.
func newFile(f *os.File) File {
	return &file{File: f}
}

// Close closes the file.
func (f *file) Close() error {
	f.mu.Lock()
	defer f.mu.Unlock()

	var err0 error
	if f.mmap != nil {
		err0 = f.mmap.Unmap()
		f.mmap = nil
	}

	err1 := f.File.Close()
	return errors.Join(err0, err1)
}

// Filename returns a file name, not a path as in *os.File.
func (f *file) Filename() string {
	_, name := filepath.Split(f.Name())
	return name
}

// Path returns *os.File.Name() as presented to Open.
func (f *file) Path() string {
	return f.File.Name()
}

// Map maps the file into memory and returns its data.
func (f *file) Map() ([]byte, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.mmap != nil {
		return f.mmap, nil
	}

	var err error
	f.mmap, err = mmap.Map(f.File, mmap.RDONLY, 0)
	if err != nil {
		return nil, err
	}

	return f.mmap, nil
}

// Size returns the file size in bytes.
func (f *file) Size() (int64, error) {
	info, err := f.Stat()
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}

// Sync syncs the file to disk.
func (f *file) Sync() error {
	fd := f.Fd()
	return syscall.Fsync(int(fd))
}
