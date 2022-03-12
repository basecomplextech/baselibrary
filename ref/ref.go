package ref

import "sync/atomic"

// Ref is an atomic reference which implements the Counter interface.
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

// Counter represents an object which supports reference counting.
type Counter interface {
	// Retain increments refcount, panics when count is 0.
	Retain() int32

	// Release decrements refcount and releases the object if the count is 0.
	Release() int32

	// Refcount returns the current refcount.
	Refcount() int32
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
func Retain[C Counter](counter C) C {
	counter.Retain()
	return counter
}

// RetainAll retains all references.
func RetainAll[C Counter](counters ...C) {
	for _, counter := range counters {
		counter.Retain()
	}
}

// ReleaseAll releases all references.
func ReleaseAll[C Counter](counters ...C) {
	for _, counter := range counters {
		counter.Release()
	}
}

// Swap releases an old reference and returns a new one.
//
// Usage:
//	new := table.Clone()
//	s.table = Swap(s.table, new)
//
func Swap[C Counter](old C, new C) C {
	old.Release()
	return new
}

var _ Counter = (*Ref)(nil)

// Retain increments refcount, panics when count is 0.
func (a *Ref) Retain() int32 {
	v := atomic.AddInt32(&a.refs, 1)
	if v <= 1 {
		panic("cannot retain already released reference")
	}
	return v
}

// Release decrements refcount and releases the object if the count is 0.
func (a *Ref) Release() int32 {
	v := atomic.AddInt32(&a.refs, -1)
	switch {
	case v > 0:
		return v
	case v < 0:
		panic("cannot release already released reference")
	}

	if a.r != nil {
		a.r.Released()
	}
	return v
}

// Refcount returns the current refcount.
func (a *Ref) Refcount() int32 {
	return atomic.LoadInt32(&a.refs)
}
