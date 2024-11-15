// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package ref

import "sync/atomic"

// Atomic32 is an atomic int32 reference counter with a released bit.
// The reference cannot be acquired if the released bit has been set.
//
//	Usage:
//
//	r := Atomic32{}
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
type Atomic32 struct {
	refs atomic.Int32
}

// Init initializes the reference counter.
func (r *Atomic32) Init(refs int32) {
	r.refs.Store(refs)
}

// Acquire increments the refcount and returs true if the reference has been acquired,
// or false if the reference has been released already.
func (r *Atomic32) Acquire() (acquired bool) {
	v := r.refs.Add(1)

	_, released := unpackAtomic32(v)
	if released {
		return false
	}
	return true
}

// Release decrements the refcount and returns true if the reference has been released,
// or false if the reference is still alive.
func (r *Atomic32) Release() (released bool) {
	v := r.refs.Add(-1)

	// Check alive or released already
	refs, releasedAlready := unpackAtomic32(v)
	switch {
	case refs < 0:
		panic("release of already released reference")
	case refs > 0 || releasedAlready:
		return false
	}

	// Release only once
	return r.refs.CompareAndSwap(0, releasedBit32)
}

// Refcount returns the current refcount.
func (r *Atomic32) Refcount() int32 {
	v := r.refs.Load()
	refs, _ := unpackAtomic32(v)
	return refs
}

// private

const releasedBit32 = int32(1 << 30)

func unpackAtomic32(r int32) (refs int32, released bool) {
	released = r&releasedBit32 != 0
	refs = r &^ releasedBit32
	return
}
