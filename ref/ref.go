package ref

import "sync/atomic"

// Count represents an object which supports reference counting.
type Count interface {
	// Retain increments refcount, panics when count is 0.
	Retain() int32

	// Release decrements refcount and releases the object if the count is 0.
	Release() int32
}

// Ref is an atomic reference which implements the Count interface.
type Ref struct {
	r    Releaser
	refs int32
}

// New returns a new atomic reference.
func New(r Releaser) Ref {
	return Ref{
		r:    r,
		refs: 1,
	}
}

// Releaser is called when reference count reaches zero.
type Releaser interface {
	Released()
}

// Retain retains and returns a reference.
//
// Usage:
//	tree.table = Retain(table)
//
func Retain[C Count](count C) C {
	count.Retain()
	return count
}

// RetainAll retains all references.
func RetainAll[C Count](counts ...C) {
	for _, count := range counts {
		count.Retain()
	}
}

// ReleaseAll releases all references.
func ReleaseAll[C Count](counts ...C) {
	for _, count := range counts {
		count.Release()
	}
}

// Swap releases an old reference and returns a new one.
//
// Usage:
//	new := table.Clone()
//	s.table = Swap(s.table, new)
//
func Swap[C Count](old C, new C) C {
	old.Release()
	return new
}

var _ Count = (*Ref)(nil)

// Retain increments refcount, panics when count is 0.
func (r *Ref) Retain() int32 {
	v := atomic.AddInt32(&r.refs, 1)
	if v <= 1 {
		panic("cannot retain already released reference")
	}
	return v
}

// Release decrements refcount and releases the object if the count is 0.
func (r *Ref) Release() int32 {
	v := atomic.AddInt32(&r.refs, -1)
	switch {
	case v > 0:
		return v
	case v < 0:
		panic("cannot release already released reference")
	}

	if r.r != nil {
		r.r.Released()
	}
	return v
}

// Refcount returns the current refcount.
func (r *Ref) Refcount() int32 {
	return atomic.LoadInt32(&r.refs)
}
