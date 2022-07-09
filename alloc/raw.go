package alloc

import "unsafe"

// Raw is a raw allocated memory segment.
type Raw struct {
	Ptr  uintptr
	Size int
}

// Bytes returns a memory segment as a byte slice.
func (r Raw) Bytes() []byte {
	uptr := *(*unsafe.Pointer)(unsafe.Pointer(&r.Ptr))
	return unsafe.Slice((*byte)(uptr), r.Size)
}
