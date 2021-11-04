package fs

import (
	"errors"
	"io"
	"io/ioutil"
	"os"

	"github.com/complex1io/library/logging"
)

var (
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

	// Rename renames a file or a directory.
	Rename(src string, dst string) error

	// Stat returns a file info.
	Stat(filename string) (FileInfo, error)

	// Temp creates a new temporary file in the directory dir, see ioutil.TempFile.
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
func New(logging logging.Logging) FileSystem {
	logger := logging.Logger("fs")
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

// Rename renames a file or a directory.
func (fs *fs) Rename(src string, dst string) error {
	return os.Rename(src, dst)
}

// Stat returns a file info.
func (fs *fs) Stat(filename string) (FileInfo, error) {
	return os.Stat(filename)
}

// Temp creates a new temporary file in the directory dir, see ioutil.TempFile.
func (fs *fs) TempFile(dir, pattern string) (File, error) {
	f, err := ioutil.TempFile(dir, pattern)
	if err != nil {
		return nil, err
	}
	return newFile(f), nil
}
