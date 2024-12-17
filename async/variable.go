// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package async

import (
	"sync"

	"github.com/basecomplextech/baselibrary/collect/chans"
	"github.com/basecomplextech/baselibrary/opt"
	"github.com/basecomplextech/baselibrary/status"
)

// [Experimental] Variable is an asynchronous variable which can be set, cleared, or failed.
type Variable[T any] interface {
	// Get returns the current value/error, or false if pending.
	Get() (T, bool, status.Status)

	// GetWait returns the current value or waits for it or an error.
	GetWait(ctx Context) (T, status.Status)

	// Set

	// Clear clears the variable.
	Clear()

	// Complete sets the value and status.
	Complete(value T, st status.Status)

	// Fail sets the error.
	Fail(st status.Status)

	// Set sets the value.
	Set(value T)

	// Wait

	// Wait returns a channel that is closed when the variable is set or failed.
	Wait() <-chan struct{}
}

// NewVariable returns a new pending async variable.
func NewVariable[T any]() Variable[T] {
	return newVariable[T]()
}

// internal

var _ Variable[any] = (*variable[any])(nil)

type variable[T any] struct {
	mu sync.RWMutex

	// either done is true or promise is set
	value T
	done  bool
	st    status.Status

	promise opt.Opt[Promise[T]]
}

func newVariable[T any]() *variable[T] {
	return &variable[T]{
		promise: opt.New(NewPromise[T]()),
	}
}

// Get returns the current value/error, or false if pending.
func (v *variable[T]) Get() (T, bool, status.Status) {
	v.mu.RLock()
	defer v.mu.RUnlock()

	return v.value, v.done, v.st
}

// GetWait returns the current value or waits for it or an error.
func (v *variable[T]) GetWait(ctx Context) (T, status.Status) {
	// Get value or promise
	value, ok, st, promise := v.get()
	if ok {
		return value, st
	}

	// Await promise
	p, _ := promise.Unwrap()
	select {
	case <-ctx.Wait():
		return value, st
	case <-p.Wait():
		return p.Result()
	}
}

// Set

// Clear clears the variable.
func (v *variable[T]) Clear() {
	v.mu.Lock()
	defer v.mu.Unlock()

	// Clear value
	var zero T
	v.value = zero
	v.done = false
	v.st = status.None

	// Make promise
	p, ok := v.promise.Unwrap()
	if !ok || p.Done() {
		v.promise.Set(NewPromise[T]())
	}
}

// Complete sets the value and status.
func (v *variable[T]) Complete(value T, st status.Status) {
	v.mu.Lock()
	defer v.mu.Unlock()

	// Set value
	v.value = value
	v.done = true
	v.st = st

	// Resolve promise
	p, ok := v.promise.Unwrap()
	if ok {
		p.Complete(value, st)
	}
}

// Fail sets the error.
func (v *variable[T]) Fail(st status.Status) {
	var zero T
	v.Complete(zero, st)
}

// Set sets the value.
func (v *variable[T]) Set(value T) {
	v.Complete(value, status.OK)
}

// Wait

// Wait returns a channel that is closed when the variable is set or failed.
func (v *variable[T]) Wait() <-chan struct{} {
	_, ok, _, promise := v.get()
	if ok {
		return chans.Closed()
	}

	p, _ := promise.Unwrap()
	return p.Wait()
}

// private

func (v *variable[T]) get() (T, bool, status.Status, opt.Opt[Promise[T]]) {
	v.mu.RLock()
	defer v.mu.RUnlock()

	return v.value, v.done, v.st, v.promise
}
