package memfs

import "github.com/epochtimeout/baselibrary/fs"

type memEntry interface {
	fs.File

	getInfo() *memInfo
	getName() string
	getParent() *memDir

	open() error
	delete() error
	move(newName string, newParent *memDir) error
}
