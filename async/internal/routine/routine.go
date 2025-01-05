// Copyright 2022 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package routine

import (
	"sync"

	"github.com/basecomplextech/baselibrary/async/internal/context"
	"github.com/basecomplextech/baselibrary/async/internal/future"
	"github.com/basecomplextech/baselibrary/opt"
	"github.com/basecomplextech/baselibrary/status"
)

// Routine is an async routine which returns the result as a future, recovers on panics,
// and can be cancelled.
type Routine[T any] interface {
	future.Future[T]

	// Stop requests the routine to stop and returns a wait channel.
	Stop() <-chan struct{}

	// OnStop adds an on stop callback, or returns false if the routine is already stopped.
	OnStop(fn func(Routine[T])) bool
}

// RoutineDyn is a routine interface without generics, i.e. Routine[?].
type RoutineDyn interface {
	future.FutureDyn

	// Stop requests the routine to stop and returns a wait channel.
	Stop() <-chan struct{}
}

// RoutineVoid is a routine which has no result.
type RoutineVoid = Routine[struct{}]

// Go

// Go runs a function in a new routine, recovers on panics.
func Go(fn func(ctx context.Context) status.Status) RoutineVoid {
	fn1 := func(ctx context.Context) (struct{}, status.Status) {
		return struct{}{}, fn(ctx)
	}

	r := newRoutine(fn1)
	go r.run()
	return r
}

// Run runs a function in a new routine, and returns the result, recovers on panics.
func Run[T any](fn func(ctx context.Context) (T, status.Status)) Routine[T] {
	r := newRoutine(fn)
	go r.run()
	return r
}

// Exited returns a routine which has exited with the given result and status.
func Exited[T any](result T, st status.Status) Routine[T] {
	r := newRoutine[T](nil)
	r.complete(result, st)
	return r
}

// internal

var _ Routine[any] = (*routine[any])(nil)

type routine[T any] struct {
	ctx context.CancelContext
	fn  func(ctx context.Context) (T, status.Status)

	mu       sync.Mutex
	wait     chan struct{}
	callback opt.Opt[func(Routine[T])]

	st     status.Status
	done   bool
	result T
}

func newRoutine[T any](fn func(ctx context.Context) (T, status.Status)) *routine[T] {
	return &routine[T]{
		ctx: context.New(),
		fn:  fn,

		wait: make(chan struct{}),
	}
}

// Future

// Done returns true if the routine is complete.
func (r *routine[T]) Done() bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.done
}

// Wait returns a channel which is closed when the future is complete.
func (r *routine[T]) Wait() <-chan struct{} {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.wait
}

// Result returns a value and a status.
func (r *routine[T]) Result() (T, status.Status) {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.result, r.st
}

// Status returns a status or none.
func (r *routine[T]) Status() status.Status {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.st
}

// Routine

// Stop requests the routine to stop and returns a wait channel.
func (r *routine[T]) Stop() <-chan struct{} {
	r.ctx.Cancel()

	r.mu.Lock()
	defer r.mu.Unlock()

	return r.wait
}

// OnStop adds an on stop callback, or returns false if the routine is already stopped.
func (r *routine[T]) OnStop(fn func(Routine[T])) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.done {
		return false
	}

	cb := fn
	prev, ok := r.callback.Unwrap()
	if ok {
		cb = func(r Routine[T]) {
			defer fn(r)

			prev(r)
		}
	}

	r.callback = opt.New(cb)
	return true
}

// private

func (r *routine[T]) run() {
	defer r.ctx.Free()
	defer func() {
		if e := recover(); e != nil {
			st := status.Recover(e)
			r.reject(st)
		}
	}()

	result, st := r.fn(r.ctx)
	r.complete(result, st)
}

func (r *routine[T]) reject(st status.Status) {
	var zero T
	r.complete(zero, st)
}

func (r *routine[T]) complete(result T, st status.Status) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.done {
		return
	}

	r.st = st
	r.done = true
	r.result = result
	close(r.wait)

	if cb, ok := r.callback.Unwrap(); ok {
		cb(r)
	}
}
