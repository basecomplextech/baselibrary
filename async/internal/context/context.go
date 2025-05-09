// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package context

import (
	context_ "context"
	"sync/atomic"
	"time"

	"github.com/basecomplextech/baselibrary/collect/chans"
	"github.com/basecomplextech/baselibrary/ref"
	"github.com/basecomplextech/baselibrary/status"
)

// Context is an async cancellation context.
//
// Usage:
//
//	ctx := NewContext()
//	defer ctx.Free()
type Context interface {
	// Done returns true if the context is cancelled.
	Done() bool

	// Wait returns a channel which is closed when the context is cancelled.
	Wait() <-chan struct{}

	// Status returns a cancellation status or OK.
	Status() status.Status

	// Callbacks

	// AddCallback adds a callback.
	AddCallback(c Callback)

	// RemoveCallback removes a callback.
	RemoveCallback(c Callback)

	// Internal

	// Free cancels and releases the context.
	Free()
}

// CancelContext is a cancellable async context.
type CancelContext interface {
	Context

	// Cancel cancels the context.
	Cancel()
}

// Callback is called when the context is cancelled.
type Callback interface {
	// OnCancelled is called when the context is cancelled.
	OnCancelled(status.Status)
}

// New

// New returns a new cancellable context.
func New() CancelContext {
	return newContext(nil /* no parent */)
}

// No returns a non-cancellable background context.
func No() Context {
	return no
}

// Cancelled returns a cancelled context.
func Cancelled() CancelContext {
	return done
}

// Timeout

// Timeout returns a context with a timeout.
func Timeout(timeout time.Duration) Context {
	return newContextTimeout(nil /* no parent */, timeout)
}

// Deadline returns a context with a deadline.
func Deadline(deadline time.Time) Context {
	timeout := time.Until(deadline)
	return newContextTimeout(nil /* no parent */, timeout)
}

// Next

// Next returns a child context.
func Next(parent Context) CancelContext {
	return newContext(parent)
}

// NextTimeout returns a child context with a timeout.
func NextTimeout(parent Context, timeout time.Duration) Context {
	return newContextTimeout(parent, timeout)
}

// NextDeadline returns a child context with a deadline.
func NextDeadline(parent Context, deadline time.Time) Context {
	timeout := time.Until(deadline)
	return newContextTimeout(parent, timeout)
}

// Standard

// Std returns a standard library context from an async one.
func Std(ctx Context) context_.Context {
	return newStdContext(ctx)
}

// internal

var _ CancelContext = (*context)(nil)

type context struct {
	refs  ref.Atomic32 // 1 by default, with released bit
	freed atomic.Bool  // free only once
	state atomic.Pointer[state]
}

func newContext(parent Context) *context {
	s := newState(parent)

	x := &context{}
	x.refs.Init(1)
	x.state.Store(s)

	// Maybe add callback
	if parent != nil {
		parent.AddCallback(x)
	}
	return x
}

func newContextTimeout(parent Context, timeout time.Duration) *context {
	x := newContext(parent)

	// Maybe already timed out
	if timeout <= 0 {
		x.timeout()
		return x
	}

	// Start timer
	timer := time.AfterFunc(timeout, x.timeout)
	s := x.state.Load()
	s.timer.set(timer)
	return x
}

// Cancel cancels the context.
func (x *context) Cancel() {
	x.cancel(status.Cancelled)
}

// Done returns true if the context is cancelled.
func (x *context) Done() bool {
	s, ok := x.acquire()
	if !ok {
		return true
	}
	defer x.release()

	done, _ := s.result.get()
	return done
}

// Wait returns a channel which is closed when the context is cancelled.
func (x *context) Wait() <-chan struct{} {
	s, ok := x.acquire()
	if !ok {
		return chans.Closed()
	}
	defer x.release()

	return s.result.wait()
}

// Status returns a cancellation status or OK.
func (x *context) Status() status.Status {
	s, ok := x.acquire()
	if !ok {
		return status.Cancelled
	}
	defer x.release()

	_, st := s.result.get()
	return st
}

// Callbacks

// AddCallback adds a callback.
func (x *context) AddCallback(c Callback) {
	s, ok := x.acquire()
	if !ok {
		c.OnCancelled(status.Cancelled)
		return
	}
	defer x.release()

	s.callbacks.add(c, &s.result)
}

// RemoveCallback removes a callback.
func (x *context) RemoveCallback(c Callback) {
	s, ok := x.acquire()
	if !ok {
		return
	}
	defer x.release()

	ok, _ = s.result.get()
	if ok {
		return
	}

	s.callbacks.remove(c)
}

// Internal

// Free cancels and releases the context.
func (x *context) Free() {
	ok := x.freed.CompareAndSwap(false, true)
	if !ok {
		return
	}

	x.cancel(status.Cancelled)
	x.release()
}

// Parent

// OnCancelled is called when the context is cancelled.
func (x *context) OnCancelled(st status.Status) {
	x.cancel(st)
}

// internal

func (x *context) cancel(st status.Status) {
	s, ok := x.acquire()
	if !ok {
		return
	}
	defer x.release()

	s.cancel(x, st)
}

func (x *context) timeout() {
	s, ok := x.acquire()
	if !ok {
		return
	}
	defer x.release()

	s.cancel(x, status.Timeout)
}

// private

// acquire increments refs and returns the state, or immediately releases it if released.
func (x *context) acquire() (*state, bool) {
	acquired := x.refs.Acquire()
	if acquired {
		s := x.state.Load()
		return s, true
	}

	// Release immediately
	x.release()
	return nil, false
}

// release decrements refs and returns the state to the pool if refs reach zero.
func (x *context) release() {
	released := x.refs.Release()
	if !released {
		return
	}

	// Release state
	s := x.state.Swap(nil)
	s.reset()
	releaseState(s)
}
