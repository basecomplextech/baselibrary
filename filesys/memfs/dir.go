// Copyright 2022 Ivan Korobkov. All rights reserved.

package memfs

import (
	"errors"

	"github.com/basecomplextech/baselibrary/filesys"
)

type memDir struct {
	fs *memFS

	name    string
	parent  *memDir
	entries map[string]memEntry
}

func newMemDir(fs *memFS, name string, parent *memDir) *memDir {
	return &memDir{
		fs: fs,

		name:    name,
		parent:  parent,
		entries: make(map[string]memEntry),
	}
}

// Size returns the file size in bytes.
func (d *memDir) Size() (int64, error) {
	panic("not implemented")
}

// Methods

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
func (d *memDir) Readdir(n int) ([]filesys.FileInfo, error) {
	d.fs.mu.RLock()
	defer d.fs.mu.RUnlock()

	var infos []filesys.FileInfo
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
func (d *memDir) Stat() (filesys.FileInfo, error) {
	d.fs.mu.RLock()
	defer d.fs.mu.RUnlock()

	info := d.getInfo()
	return info, nil
}

// Sync syncs the file to disk.
func (d *memDir) Sync() error {
	d.fs.mu.Lock()
	defer d.fs.mu.Unlock()

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

func (d *memDir) isDir() bool {
	return true
}

func (d *memDir) isEmpty() bool {
	return len(d.entries) == 0
}

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

func (d *memDir) getParent() *memDir {
	return d.parent
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
		return nil, filesys.ErrExist
	}

	f := newMemFile(d.fs, name, d)
	d.entries[name] = f
	return f, nil
}

func (d *memDir) makeDir(name string) (*memDir, error) {
	_, ok := d.entries[name]
	if ok {
		return nil, filesys.ErrExist
	}

	dir := newMemDir(d.fs, name, d)
	d.entries[name] = dir
	return dir, nil
}

func (d *memDir) makePath(names ...string) (*memDir, error) {
	if len(names) == 0 {
		return nil, filesys.ErrInvalid
	}
	name := names[0]

	var dir *memDir
	if e, ok := d.entries[name]; ok {
		dir, ok = e.(*memDir)
		if !ok {
			return nil, errors.New("path is not a directory")
		}
	} else {
		var err error
		dir, err = d.makeDir(name)
		if err != nil {
			return nil, err
		}
	}

	if len(names) == 1 {
		return dir, nil
	}

	return dir.makePath(names[1:]...)
}

func (d *memDir) open() error {
	return nil
}

func (d *memDir) remove(name string) error {
	_, ok := d.entries[name]
	if !ok {
		return filesys.ErrNotExist
	}

	delete(d.entries, name)
	return nil
}

func (d *memDir) addEntry(entry memEntry, name string) {
	d.entries[name] = entry
	entry.move(name, d)
}

func (d *memDir) removeEntry(e memEntry) {
	name := e.getName()
	delete(d.entries, name)
}
