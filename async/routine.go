// Copyright 2022 Ivan Korobkov. All rights reserved.

package async

import (
	"github.com/basecomplextech/baselibrary/status"
)

// Routine is an async routine which returns the result as a future, recovers on panics,
// and can be cancelled.
type Routine[T any] interface {
	Future[T]

	// Stop requests the routine to stop and returns a wait channel.
	Stop() <-chan struct{}
}

// RoutineDyn is a routine interface without generics, i.e. Routine[?].
type RoutineDyn interface {
	FutureDyn

	Stop() <-chan struct{}
}

// Go

// Go runs a function in a new routine, recovers on panics.
func Go(fn func(ctx Context) status.Status) Routine[struct{}] {
	r := newRoutine[struct{}]()

	go func() {
		defer r.ctx.Free()
		defer func() {
			e := recover()
			if e == nil {
				return
			}

			st := status.Recover(e)
			r.result.Complete(struct{}{}, st)
		}()

		st := fn(r.ctx)
		r.result.Complete(struct{}{}, st)
	}()

	return r
}

// GoPool runs a function in a pool, recovers on panics.
func GoPool(pool Pool, fn func(ctx Context) status.Status) Routine[struct{}] {
	r := newRoutine[struct{}]()

	pool.Go(func() {
		defer r.ctx.Free()
		defer func() {
			e := recover()
			if e == nil {
				return
			}

			st := status.Recover(e)
			r.result.Complete(struct{}{}, st)
		}()

		st := fn(r.ctx)
		r.result.Complete(struct{}{}, st)
	})

	return r
}

// Call

// Call calls a function in a new routine, and returns its result, recovers on panics.
func Call[T any](fn func(ctx Context) (T, status.Status)) Routine[T] {
	r := newRoutine[T]()

	go func() {
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

		result, st := fn(r.ctx)
		r.result.Complete(result, st)
	}()

	return r
}

// CallPool calls a function in a pool, and returns its result, recovers on panics.
func CallPool[T any](pool Pool, fn func(ctx Context) (T, status.Status)) Routine[T] {
	r := newRoutine[T]()

	pool.Go(func() {
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

		result, st := fn(r.ctx)
		r.result.Complete(result, st)
	})

	return r
}

// Exited

// Exited returns a routine which has exited with the given result and status.
func Exited[T any](result T, st status.Status) Routine[T] {
	r := newRoutine[T]()
	r.result.Complete(result, st)
	return r
}

// internal

var _ Routine[any] = (*routine[any])(nil)

type routine[T any] struct {
	ctx    Context
	result promise[T]
}

func newRoutine[T any]() *routine[T] {
	return &routine[T]{
		ctx:    newContext(nil /* no parent */),
		result: newPromiseEmbedded[T](),
	}
}

// Stop requests the future to stop and returns a wait channel.
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
