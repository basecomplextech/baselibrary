package fs

import (
	"errors"
	"io"
	"io/ioutil"
	"os"

	"github.com/epochtimeout/baselibrary/logging"
)

var (
	ErrExist    = os.ErrExist
	ErrNotExist = os.ErrNotExist

	ErrClosed           = errors.New("operation on a closed file")
	ErrClosedFileWriter = errors.New("operation on a closed file writer")
)

type FileInfo = os.FileInfo

type FileSystem interface {
	// Create creates a file in the fstem.
	Create(name string) (File, error)

	// Mkdir creates a directory in the fstem.
	Mkdir(name string, perm os.FileMode) error

	// MkdirAll creates a directory path and all parents that does not exist.
	MkdirAll(path string, perm os.FileMode) error

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

// File is an extended file interface for *os.File.
type File interface {
	// Filename returns a file name, not a path as in *os.File.
	Filename() string

	// Ext returns a file extension.
	Ext() string

	// Path returns *os.File.Name() as presented to Open.
	Path() string

	// Map maps the file into memory and returns its data.
	Map() ([]byte, error)

	// Size returns the file size in bytes.
	Size() (int64, error)

	// os.File methods

	io.Closer
	io.Reader
	io.ReaderAt
	io.Seeker
	io.Writer
	io.WriterAt

	Name() string
	Readdir(count int) ([]os.FileInfo, error)
	Readdirnames(n int) ([]string, error)
	Stat() (os.FileInfo, error)
	Sync() error
	Truncate(size int64) error
	WriteString(s string) (ret int, err error)
}

// FileWriter writes a single file.
type FileWriter interface {
	io.Closer
	io.Writer
	Done() (File, error)
}

// New returns a disk file system.
func New(logger logging.Logger) FileSystem {
	return newFileSystem(logger)
}

type fs struct {
	logger logging.Logger
}

func newFileSystem(logger logging.Logger) *fs {
	return &fs{logger: logger}
}

// Create creates a file in the fstem.
func (fs *fs) Create(name string) (File, error) {
	f, err := os.Create(name)
	if err != nil {
		return nil, err
	}
	return newFile(f), nil
}

// Mkdir creates a directory in the fstem.
func (fs *fs) Mkdir(name string, perm os.FileMode) error {
	return os.Mkdir(name, perm)
}

// MkdirAll creates a directory path and all parents that does not exist.
func (fs *fs) MkdirAll(path string, perm os.FileMode) error {
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
