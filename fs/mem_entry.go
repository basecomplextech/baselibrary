package fs

type memEntry interface {
	File

	getInfo() *memInfo
	getName() string
	getParent() *memDir

	open() error
	delete() error
	move(newName string, newParent *memDir) error
}
