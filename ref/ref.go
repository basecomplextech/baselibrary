package ref

import "sync/atomic"

// Ref is an atomic countable reference.
type Ref interface {
	// Acquire increments the count and returns true on success.
	Acquire() bool

	// Retain increments the count, panics when the count was 0.
	Retain() int32

	// Release decrements the count and releases the object when the count is 0.
	Release() int32

	// Refcount returns the current reference count.
	Refcount() int32
}

func New(r Releaser) Ref {
	return &refImpl{
		r:     r,
		count: 1,
	}
}

func Empty() Ref {
	return New(nil)
}

func NewEmpty() Ref {
	return New(nil)
}

type Releaser interface {
	Released()
}

type refImpl struct {
	r     Releaser
	count int32
}

func (r *refImpl) Acquire() bool {
	v := atomic.AddInt32(&r.count, 1)
	return v > 1
}

func (r *refImpl) Retain() int32 {
	v := atomic.AddInt32(&r.count, 1)
	if v <= 1 {
		panic("cannot retain an already released reference")
	}
	return v
}

func (r *refImpl) Release() int32 {
	v := atomic.AddInt32(&r.count, -1)
	switch {
	case v > 0:
		return v
	case v < 0:
		panic("cannot release an already released reference")
	}

	if r.r != nil {
		r.r.Released()
	}
	return v
}

func (r *refImpl) Refcount() int32 {
	return atomic.LoadInt32(&r.count)
}
