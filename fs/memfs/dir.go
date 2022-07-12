package memfs

import (
	"os"
	"path/filepath"

	"github.com/epochtimeout/baselibrary/fs"
)

var _ fs.File = (*memDir)(nil)

type memDir struct {
	fs *memFS

	name    string
	parent  *memDir
	entries map[string]memEntry

	refs    int
	opened  bool
	deleted bool
}

func newMemDir(fs *memFS, name string, parent *memDir) *memDir {
	return &memDir{
		fs: fs,

		name:    name,
		parent:  parent,
		entries: make(map[string]memEntry),
	}
}

// Filename returns a file name, not a path as in *os.File.
func (d *memDir) Filename() string {
	return d.name
}

// Name returns a file path as in *os.File.
func (d *memDir) Name() string {
	d.fs.mu.RLock()
	defer d.fs.mu.RUnlock()

	return d.getPath()
}

// Path returns *os.File.Name() as presented to Open.
func (d *memDir) Path() string {
	d.fs.mu.RLock()
	defer d.fs.mu.RUnlock()

	return d.getPath()
}

// Size returns the file size in bytes.
func (d *memDir) Size() (int64, error) {
	panic("not implemented")
}

// Methods

// Close closes the file.
func (d *memDir) Close() error {
	d.fs.mu.Lock()
	defer d.fs.mu.Unlock()

	if !d.opened {
		return os.ErrClosed
	}

	d.refs--
	if d.refs > 0 {
		return nil
	}

	d.opened = false
	return nil
}

// Map maps the file into memory and returns its data.
func (d *memDir) Map() ([]byte, error) {
	panic("not implemented")
}

// Read reads data from the file into p.
func (d *memDir) Read(p []byte) (n int, err error) {
	panic("not implemented")
}

// ReadAt reads data from the file at an offset into p.
func (d *memDir) ReadAt(p []byte, off int64) (n int, err error) {
	panic("not implemented")
}

// Readdir reads and returns the directory entries, upto n entries if n > 0.
func (d *memDir) Readdir(n int) ([]os.FileInfo, error) {
	d.fs.mu.RLock()
	defer d.fs.mu.RUnlock()

	if !d.opened {
		return nil, os.ErrClosed
	}

	var infos []os.FileInfo
	for _, entry := range d.entries {
		if n > 0 && len(infos) >= n {
			break
		}

		info := entry.getInfo()
		infos = append(infos, info)
	}

	return infos, nil
}

// Readdirnames reads and returns the directory entries, upto n entries if n > 0.
func (d *memDir) Readdirnames(n int) ([]string, error) {
	d.fs.mu.RLock()
	defer d.fs.mu.RUnlock()

	if !d.opened {
		return nil, os.ErrClosed
	}

	var names []string
	for _, entry := range d.entries {
		if n > 0 && len(names) >= n {
			break
		}

		name := entry.getName()
		names = append(names, name)
	}

	return names, nil
}

// Seek sets the file offset.
func (d *memDir) Seek(offset int64, whence int) (int64, error) {
	panic("not implemented")
}

// Stat returns a file info.
func (d *memDir) Stat() (os.FileInfo, error) {
	d.fs.mu.RUnlock()
	defer d.fs.mu.RUnlock()

	if !d.opened {
		return nil, os.ErrClosed
	}

	info := d.getInfo()
	return info, nil
}

// Sync syncs the file to disk.
func (d *memDir) Sync() error {
	d.fs.mu.Lock()
	defer d.fs.mu.Unlock()

	if !d.opened {
		return os.ErrClosed
	}
	return nil
}

// Truncate truncates the file to a given length.
func (d *memDir) Truncate(size int64) error {
	panic("not implemented")
}

// Write writes data to the file at the current offset.
func (d *memDir) Write(p []byte) (n int, err error) {
	panic("not implemented")
}

// WriteAt writes data to the file at an offset.
func (d *memDir) WriteAt(p []byte, off int64) (n int, err error) {
	panic("not implemented")
}

// WriteString writes a string to the file at the current offset.
func (d *memDir) WriteString(s string) (ret int, err error) {
	panic("not implemented")
}

// entry

func (d *memDir) getInfo() *memInfo {
	return &memInfo{
		name: d.name,
		size: 0,
		dir:  true,
	}
}

func (d *memDir) getName() string {
	return d.name
}

func (d *memDir) getPath() string {
	names := []string{d.name}
	for parent := d.parent; parent != nil; parent = parent.parent {
		names = append([]string{parent.name}, names...)
	}

	return filepath.Join(names...)
}

func (d *memDir) getParent() *memDir {
	return d.parent
}

func (d *memDir) open() error {
	if d.deleted {
		return os.ErrNotExist
	}

	d.refs++
	d.opened = true
	return nil
}

func (d *memDir) delete() error {
	d.parent = nil
	d.deleted = true
	return nil
}

func (d *memDir) move(newName string, newParent *memDir) error {
	d.name = newName
	d.parent = newParent
	return nil
}

// internal

func (d *memDir) find(name string) (memEntry, bool) {
	e, ok := d.entries[name]
	return e, ok
}

func (d *memDir) findDir(name string) (*memDir, bool) {
	e, ok := d.entries[name]
	if !ok {
		return nil, false
	}
	dir, ok := e.(*memDir)
	return dir, ok
}

func (d *memDir) findPath(names ...string) (memEntry, bool) {
	if len(names) == 0 {
		return nil, false
	}
	if len(names) == 1 {
		return d.find(names[0])
	}

	name := names[0]
	dir, ok := d.findDir(name)
	if !ok {
		return nil, false
	}

	return dir.findPath(names[1:]...)
}

func (d *memDir) create(name string) (*memFile, error) {
	_, ok := d.entries[name]
	if ok {
		return nil, os.ErrExist
	}

	f := newMemFile(d.fs, name, d)
	d.entries[name] = f
	return f, nil
}

func (d *memDir) makeDir(name string) (*memDir, error) {
	_, ok := d.entries[name]
	if ok {
		return nil, os.ErrExist
	}

	dir := newMemDir(d.fs, name, d)
	d.entries[name] = dir
	return dir, nil
}

func (d *memDir) makePath(names ...string) (*memDir, error) {
	if len(names) == 0 {
		return nil, os.ErrInvalid
	}

	name := names[0]
	dir, err := d.makeDir(name)
	if err != nil {
		return nil, err
	}
	if len(names) == 1 {
		return dir, nil
	}

	return dir.makePath(names[1:]...)
}

func (d *memDir) remove(name string) error {
	e, ok := d.entries[name]
	if !ok {
		return os.ErrNotExist
	}

	delete(d.entries, name)
	return e.delete()
}

func (d *memDir) removePath(names ...string) error {
	return nil
}

func (d *memDir) addEntry(entry memEntry, name string) {
	e, ok := d.entries[name]
	if ok {
		e.delete()
	}

	d.entries[name] = entry
	entry.move(name, d)
}

func (d *memDir) removeEntry(e memEntry) {
	name := e.getName()
	delete(d.entries, name)
}
