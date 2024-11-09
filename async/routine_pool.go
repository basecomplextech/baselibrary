// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package async

import (
	"github.com/basecomplextech/baselibrary/async/internal/routinepool"
	"github.com/basecomplextech/baselibrary/status"
)

// RoutinePool allows to reuse goroutines with preallocated big stacks.
type RoutinePool = routinepool.RoutinePool

// NewRoutinePool returns a new goroutine pool.
//
// Use [GoPool] and [RunPool] to run functions in the pool.
func NewRoutinePool() RoutinePool {
	return routinepool.New()
}

// GoPool runs a function in a pool, recovers on panics.
func GoPool(pool RoutinePool, fn func(ctx Context) status.Status) Routine[struct{}] {
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

// RunPool runs a function in a routine pool, and returns the result, recovers on panics.
func RunPool[T any](pool RoutinePool, fn func(ctx Context) (T, status.Status)) Routine[T] {
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
