// Copyright 2025 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package routine

import (
	"sync"

	"github.com/basecomplextech/baselibrary/async/internal/context"
	"github.com/basecomplextech/baselibrary/async/internal/future"
	"github.com/basecomplextech/baselibrary/async/internal/promise"
	"github.com/basecomplextech/baselibrary/status"
)

// Routine is an async routine which returns the result as a future, recovers on panics,
// and can be cancelled.
type Routine[T any] interface {
	future.Future[T]

	// Start start the routine, if not started or stopped yet.
	Start()

	// Stop requests the routine to stop and returns a wait channel.
	// The method does not call the on-stop callbacks if the routine has not started.
	Stop() <-chan struct{}

	// OnStop adds a callback which is called when the routine stops.
	// The callback is called by the routine goroutine, and is not called
	// if the routine has not started.
	OnStop(fn func(Routine[T])) bool
}

// RoutineDyn is a routine interface without generics, i.e. Routine[?].
type RoutineDyn interface {
	future.FutureDyn

	// Start starts the routine, if not stopped already.
	Start()

	// Stop requests the routine to stop and returns a wait channel.
	Stop() <-chan struct{}
}

// RoutineVoid is a routine which has no result.
type RoutineVoid = Routine[struct{}]

// Func

type (
	// Func is a function which returns the result.
	Func[T any] = func(ctx context.Context) (T, status.Status)

	// Func1 is a single argument function which returns the result.
	Func1[T any, A any] = func(ctx context.Context, arg A) (T, status.Status)

	// FuncVoid is a function which returns no result.
	FuncVoid = func(ctx context.Context) status.Status

	// FuncVoid1 is a single argument function which returns no result
	FuncVoid1[A any] = func(ctx context.Context, arg A) status.Status
)

// New

// New returns a new routine, but does not start it.
func New[T any](fn Func[T]) Routine[T] {
	return newRoutine(fn)
}

// NewVoid returns a new routine without a result, but does not start it.
func NewVoid(fn func(ctx context.Context) status.Status) RoutineVoid {
	fn1 := func(ctx context.Context) (struct{}, status.Status) {
		return struct{}{}, fn(ctx)
	}

	return newRoutine(fn1)
}

// Run

// Run runs a function in a new routine, and returns the result, recovers on panics.
func Run[T any](fn Func[T]) Routine[T] {
	r := newRoutine(fn)
	r.Start()
	return r
}

// Run1 runs a function in a new routine, and returns the result, recovers on panics.
func Run1[T any, A any](fn Func1[T, A], arg A) Routine[T] {
	fn1 := func(ctx context.Context) (T, status.Status) {
		return fn(ctx, arg)
	}

	r := newRoutine(fn1)
	r.Start()
	return r
}

// RunVoid

// RunVoid runs a procedure in a new routine, recovers on panics.
func RunVoid(fn FuncVoid) RoutineVoid {
	fn1 := func(ctx context.Context) (struct{}, status.Status) {
		return struct{}{}, fn(ctx)
	}

	r := newRoutine(fn1)
	r.Start()
	return r
}

// RunVoid1 runs a procedure in a new routine, recovers on panics.
func RunVoid1[A any](fn FuncVoid1[A], arg A) RoutineVoid {
	fn1 := func(ctx context.Context) (struct{}, status.Status) {
		return struct{}{}, fn(ctx, arg)
	}

	r := newRoutine(fn1)
	r.Start()
	return r
}

// Stopped

// Stopped returns a routine which has stopped with the given result and status.
func Stopped[T any](result T, st status.Status) Routine[T] {
	r := newRoutine[T](nil)
	r.complete(result, st)
	return r
}

// StoppedVoid returns a routine without a result which has stopped with the given status.
func StoppedVoid(st status.Status) RoutineVoid {
	r := newRoutine[struct{}](nil)
	r.complete(struct{}{}, st)
	return r
}

// internal

var _ Routine[any] = (*routine1[any])(nil)

type routine1[T any] struct {
	ctx context.CancelContext
	fn  Func[T]

	mu       sync.Mutex
	promise  promise.Promise[T]
	callback func(Routine[T]) // maybe nil

	start bool // start has been called
	stop  bool // stop has been called
}

func newRoutine[T any](fn Func[T]) *routine1[T] {
	ctx := context.New()
	promise := promise.New[T]()

	return &routine1[T]{
		ctx:     ctx,
		fn:      fn,
		promise: promise,
	}
}

// Future

// Done returns true if the future is complete.
func (r *routine1[T]) Done() bool {
	return r.promise.Done()
}

// Wait returns a channel which is closed when the future is complete.
func (r *routine1[T]) Wait() <-chan struct{} {
	return r.promise.Wait()
}

// Result returns a value and a status.
func (r *routine1[T]) Result() (T, status.Status) {
	return r.promise.Result()
}

// Status returns a status or none.
func (r *routine1[T]) Status() status.Status {
	return r.promise.Status()
}

// Routine

// Start start the routine, if not started or stopped yet.
func (r *routine1[T]) Start() {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.start || r.stop {
		return
	}

	r.start = true
	go r.run()
}

// Stop requests the routine to stop and returns a wait channel.
// The method does not call the on-stop callbacks if the routine has not started.
func (r *routine1[T]) Stop() <-chan struct{} {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Stop already requested
	if r.stop {
		return r.promise.Wait()
	}

	// Reject if not started
	if !r.start {
		r.promise.Reject(status.Cancelled)
		return r.promise.Wait()
	}

	// Cancel context and return wait
	r.ctx.Cancel()
	r.stop = true
	return r.promise.Wait()
}

// OnStop adds a callback which is called when the routine stops.
// The callback is called by the routine goroutine, and is not called
// if the routine has not started.
func (r *routine1[T]) OnStop(fn func(Routine[T])) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Return false if done
	if r.promise.Done() {
		return false
	}

	// First callback
	if r.callback == nil {
		r.callback = fn
		return true
	}

	// Chain callbacks
	last := r.callback
	r.callback = func(r Routine[T]) {
		defer fn(r)
		last(r)
	}
	return true
}

// private

func (r *routine1[T]) run() {
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

func (r *routine1[T]) reject(st status.Status) {
	var zero T
	r.complete(zero, st)
}

func (r *routine1[T]) complete(result T, st status.Status) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Complete promise
	ok := r.promise.Complete(result, st)
	if !ok {
		return
	}

	// Notify callback
	if r.callback != nil {
		r.callback(r)
	}
}
