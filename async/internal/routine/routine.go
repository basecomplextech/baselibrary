// Copyright 2022 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package routine

import (
	"github.com/basecomplextech/baselibrary/async/internal/context"
	"github.com/basecomplextech/baselibrary/async/internal/future"
	"github.com/basecomplextech/baselibrary/async/internal/promise"
	"github.com/basecomplextech/baselibrary/status"
)

// Routine is an async routine which returns the result as a future, recovers on panics,
// and can be cancelled.
type Routine[T any] interface {
	future.Future[T]

	// Stop requests the routine to stop and returns a wait channel.
	Stop() <-chan struct{}
}

// RoutineDyn is a routine interface without generics, i.e. Routine[?].
type RoutineDyn interface {
	future.FutureDyn

	// Stop requests the routine to stop and returns a wait channel.
	Stop() <-chan struct{}
}

// Go

// Go runs a function in a new routine, recovers on panics.
func Go(fn func(ctx context.Context) status.Status) Routine[struct{}] {
	fn1 := func(ctx context.Context) (struct{}, status.Status) {
		return struct{}{}, fn(ctx)
	}

	r := newRoutine(fn1)
	go r.Run()
	return r
}

// Run runs a function in a new routine, and returns the result, recovers on panics.
func Run[T any](fn func(ctx context.Context) (T, status.Status)) Routine[T] {
	r := newRoutine(fn)
	go r.Run()
	return r
}

// Exited returns a routine which has exited with the given result and status.
func Exited[T any](result T, st status.Status) Routine[T] {
	r := newRoutine[T](nil)
	r.result.Complete(result, st)
	return r
}

// internal

var _ Routine[any] = (*routine[any])(nil)

type routine[T any] struct {
	ctx    context.CancelContext
	fn     func(ctx context.Context) (T, status.Status)
	result promise.Promise[T]
}

func newRoutine[T any](fn func(ctx context.Context) (T, status.Status)) *routine[T] {
	return &routine[T]{
		ctx:    context.New(),
		result: promise.New[T](),
		fn:     fn,
	}
}

func newRoutineEmbedded[T any](fn func(ctx context.Context) (T, status.Status)) routine[T] {
	return routine[T]{
		ctx:    context.New(),
		result: promise.New[T](),
		fn:     fn,
	}
}

// Done returns true if the routine is complete.
func (r *routine[T]) Done() bool {
	return r.result.Done()
}

// Stop requests the routine to stop and returns a wait channel.
func (r *routine[T]) Stop() <-chan struct{} {
	r.ctx.Cancel()
	return r.result.Wait()
}

// Wait returns a channel which is closed when the future is complete.
func (r *routine[T]) Wait() <-chan struct{} {
	return r.result.Wait()
}

// Result returns a value and a status.
func (r *routine[T]) Result() (T, status.Status) {
	return r.result.Result()
}

// Status returns a status or none.
func (r *routine[T]) Status() status.Status {
	return r.result.Status()
}

// private

func (r *routine[T]) Run() {
	defer r.ctx.Free()
	defer func() {
		e := recover()
		if e == nil {
			return
		}

		var zero T
		st := status.Recover(e)
		r.result.Complete(zero, st)
	}()

	result, st := r.fn(r.ctx)
	r.result.Complete(result, st)
}
