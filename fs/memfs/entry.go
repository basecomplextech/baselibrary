package memfs

import "github.com/epochtimeout/baselibrary/fs"

type memEntry interface {
	fs.File

	isDir() bool
	isEmpty() bool

	getInfo() *memInfo
	getName() string
	getParent() *memDir

	open() error
	delete() error
	move(newName string, newParent *memDir) error
}
