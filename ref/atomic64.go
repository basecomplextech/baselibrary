// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package ref

import "sync/atomic"

// Atomic64 is an atomic int64 reference counter with a released bit.
// The reference cannot be acquired if the released bit has been set.
//
//	Usage:
//
//	r := Atomic64{}
//
//	acquired := r.Acquire()
//	if acquired {
//		return
//	}
//
//	released := r.Release()
//	if !released {
//		return
//	}
//
//	free()
type Atomic64 struct {
	refs atomic.Int64
}

// Init initializes the reference counter.
func (r *Atomic64) Init(refs int64) {
	r.refs.Store(refs)
}

// Acquire increments the refcount and returs true if the reference has been acquired,
// or false if the reference has been released already.
func (r *Atomic64) Acquire() (acquired bool) {
	v := r.refs.Add(1)

	_, released := unpackAtomic64(v)
	if released {
		return false
	}
	return true
}

// Retain increments the refcount, panics if the reference has been released already.
func (r *Atomic64) Retain() {
	v := r.refs.Add(1)

	_, released := unpackAtomic64(v)
	if released {
		panic("retain of already released reference")
	}
}

// Release decrements the refcount and returns true if the reference has been released,
// or false if the reference is still alive.
func (r *Atomic64) Release() (released bool) {
	v := r.refs.Add(-1)

	// Check alive or released already
	refs, released := unpackAtomic64(v)
	switch {
	case refs < 0:
		panic("release of already released reference")
	case refs > 0 || released:
		return false
	}

	// Release only once
	return r.refs.CompareAndSwap(0, releasedBit64)
}

// Refcount returns the current refcount.
func (r *Atomic64) Refcount() int64 {
	v := r.refs.Load()
	refs, _ := unpackAtomic64(v)
	return refs
}

// private

const releasedBit64 = int64(1 << 62)

func unpackAtomic64(r int64) (refs int64, released bool) {
	released = r&releasedBit64 != 0
	refs = r &^ releasedBit64
	return
}
