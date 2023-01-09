package alloc

func init() {
	initBlockClasses()
	initGlobal()
}

func initGlobal() {
	global = newAllocator()
}
