// Copyright 2022 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package flag

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// Flag is a read-only boolean flag that can be waited on until set.
//
// Example:
//
//	serving := async.UnsetFlag()
//
//	func handle(ctx Context, req *request) {
//		if !serving.Get() { // just to show Get in example
//			select {
//			case <-ctx.Wait():
//				return ctx.Status()
//			case <-serving.Wait():
//			}
//		}
//
//		// ... handle request
//	}
type Flag interface {
	// IsSet returns true if the flag is set.
	// The method uses an atomic boolean internally and is non-blocking.
	IsSet() bool

	// Wait waits for the flag to be set.
	Wait() <-chan struct{}
}

// MutFlag is a routine-safe boolean flag that can be set, reset, and waited on until set.
//
// Example:
//
//	serving := async.UnsetFlag()
//
//	func serve() {
//		s.serving.Set()
//		defer s.serving.Unset()
//
//		// ... start server ...
//	}
//
//	func handle(ctx Context, req *request) {
//		select {
//		case <-ctx.Wait():
//			return ctx.Status()
//		case <-serving.Wait():
//		}
//
//		// ... handle request
//	}
type MutFlag interface {
	Flag

	// Set sets the flag and notifies the waiters.
	Set()

	// Unset unsets the flag and replaces its wait channel with an open one.
	Unset()
}

// SetFlag returns a new set flag.
func SetFlag() MutFlag {
	f := UnsetFlag()
	f.Set()
	return f
}

// UnsetFlag returns a new unset flag.
func UnsetFlag() MutFlag {
	return newFlag()
}

// ReverseFlag returns a new flag which reverses the original one.
func ReverseFlag(f Flag) Flag {
	switch f := f.(type) {
	case *flag:
		return newReverseFlag(f)
	case *reverseFlag:
		return f.src
	}

	panic(fmt.Sprintf("unsupported flag type %T", f))
}

// internal

type flag struct {
	mu      sync.Mutex
	set     atomic.Bool
	setChan chan struct{} // closed when set

	reverse []*reverseFlag
}

func newFlag() *flag {
	return &flag{
		setChan: make(chan struct{}),
	}
}

// IsSet returns true if the flag is set.
// The method uses an atomic boolean internally and is non-blocking.
func (f *flag) IsSet() bool {
	return f.set.Load()
}

// Set sets the flag and notifies the waiters.
func (f *flag) Set() {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.set.Load() {
		return
	}

	f.set.Store(true)
	close(f.setChan)

	for _, rf := range f.reverse {
		rf.srcSet()
	}
}

// Unset unsets the flag and replaces its wait channel with an open one.
func (f *flag) Unset() {
	f.mu.Lock()
	defer f.mu.Unlock()

	if !f.set.Load() {
		return
	}

	f.set.Store(false)
	f.setChan = make(chan struct{})

	for _, rf := range f.reverse {
		rf.srcUnset()
	}
}

// Wait waits for the flag to be set.
func (f *flag) Wait() <-chan struct{} {
	f.mu.Lock()
	defer f.mu.Unlock()

	return f.setChan
}

// reverse

var _ Flag = (*reverseFlag)(nil)

type reverseFlag struct {
	src *flag

	mu    sync.Mutex
	unset chan struct{} // closed (set) when source unset
}

func newReverseFlag(src *flag) *reverseFlag {
	f := &reverseFlag{
		src:   src,
		unset: make(chan struct{}),
	}

	src.mu.Lock()
	defer src.mu.Unlock()

	src.reverse = append(src.reverse, f)

	if !src.set.Load() {
		close(f.unset)
	}
	return f
}

// IsSet returns true if the flag is set.
// The method uses an atomic boolean internally and is non-blocking.
func (f *reverseFlag) IsSet() bool {
	return !f.src.IsSet()
}

// Wait waits for the flag to be set.
func (f *reverseFlag) Wait() <-chan struct{} {
	f.mu.Lock()
	defer f.mu.Unlock()

	return f.unset
}

// srcSet unsets the reverse flag.
func (f *reverseFlag) srcSet() {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.unset = make(chan struct{})
}

// srcUnset sets the reverse flag.
func (f *reverseFlag) srcUnset() {
	f.mu.Lock()
	defer f.mu.Unlock()

	close(f.unset)
}
