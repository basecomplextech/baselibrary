package fs

import "os"

var _ File = (*memFile)(nil)

type memFile struct{}

func newMemFile(fs *memFS, name string, parent *memDir) *memFile {
	return nil
}

// Filename returns a file name, not a path as in *os.File.
func (f *memFile) Filename() string {
	return ""
}

// Name returns a file path as in *os.File.
func (f *memFile) Name() string {
	return ""
}

// Path returns *os.File.Name() as presented to Open.
func (f *memFile) Path() string {
	return ""
}

// Size returns the file size in bytes.
func (f *memFile) Size() (int64, error) {
	return 0, nil
}

// Methods

// Close closes the file.
func (f *memFile) Close() error {
	return nil
}

// Map maps the file into memory and returns its data.
func (f *memFile) Map() ([]byte, error) {
	return nil, nil
}

// Read reads data from the file into p.
func (f *memFile) Read(p []byte) (n int, err error) {
	return 0, nil
}

// ReadAt reads data from the file at an offset into p.
func (f *memFile) ReadAt(p []byte, off int64) (n int, err error) {
	return 0, nil
}

// Readdir reads and returns the directory entries, upto n entries if n > 0.
func (f *memFile) Readdir(count int) ([]os.FileInfo, error) {
	return nil, nil
}

// Readdirnames reads and returns the directory entries, upto n entries if n > 0.
func (f *memFile) Readdirnames(n int) ([]string, error) {
	return nil, nil
}

// Seek sets the file offset.
func (f *memFile) Seek(offset int64, whence int) (int64, error) {
	return 0, nil
}

// Stat returns a file info.
func (f *memFile) Stat() (os.FileInfo, error) {
	return nil, nil
}

// Sync syncs the file to disk.
func (f *memFile) Sync() error {
	return nil
}

// Truncate truncates the file to a given length.
func (f *memFile) Truncate(size int64) error {
	return nil
}

// Write writes data to the file at the current offset.
func (f *memFile) Write(p []byte) (n int, err error) {
	return 0, nil
}

// WriteAt writes data to the file at an offset.
func (f *memFile) WriteAt(p []byte, off int64) (n int, err error) {
	return 0, nil
}

// WriteString writes a string to the file at the current offset.
func (f *memFile) WriteString(s string) (ret int, err error) {
	return 0, nil
}

// entry

func (f *memFile) getInfo() *memInfo {
	return nil
}

func (f *memFile) getName() string {
	return ""
}

func (f *memFile) getParent() *memDir {
	return nil
}

func (f *memFile) open() error {
	return nil
}

func (f *memFile) delete() error {
	return nil
}

func (f *memFile) move(newName string, newParent *memDir) error {
	return nil
}
