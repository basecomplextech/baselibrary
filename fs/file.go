package fs

import (
	"os"
	"path/filepath"
	"sync"
	"syscall"

	"github.com/edsrzf/mmap-go"
	"github.com/epochtimeout/baselibrary/errors2"
)

// File is an extended file interface for *os.File.
type File interface {
	// Filename returns a file name, not a path as in *os.File.
	Filename() string

	// Name returns a file path as in *os.File.
	Name() string

	// Path returns *os.File.Name() as presented to Open.
	Path() string

	// Size returns the file size in bytes.
	Size() (int64, error)

	// Methods

	// Close closes the file.
	Close() error

	// Map maps the file into memory and returns its data.
	Map() ([]byte, error)

	// Read reads data from the file into p.
	Read(p []byte) (n int, err error)

	// ReadAt reads data from the file at an offset into p.
	ReadAt(p []byte, off int64) (n int, err error)

	// Readdir reads and returns the directory entries, upto n entries if n > 0.
	Readdir(count int) ([]os.FileInfo, error)

	// Readdirnames reads and returns the directory entries, upto n entries if n > 0.
	Readdirnames(n int) ([]string, error)

	// Seek sets the file offset.
	Seek(offset int64, whence int) (int64, error)

	// Stat returns a file info.
	Stat() (os.FileInfo, error)

	// Sync syncs the file to disk.
	Sync() error

	// Truncate truncates the file to a given length.
	Truncate(size int64) error

	// Write writes data to the file at the current offset.
	Write(p []byte) (n int, err error)

	// WriteAt writes data to the file at an offset.
	WriteAt(p []byte, off int64) (n int, err error)

	// WriteString writes a string to the file at the current offset.
	WriteString(s string) (ret int, err error)
}

// FileInfo describes a file and is returned by Stat and Lstat.
type FileInfo = os.FileInfo

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
	return errors2.Combine(err0, err1)
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
