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
	fn1 := func(ctx Context) (struct{}, status.Status) {
		return struct{}{}, fn(ctx)
	}

	r := newRoutine[struct{}](fn1)
	pool.Run(r)
	return r
}

// RunPool runs a function in a routine pool, and returns the result, recovers on panics.
func RunPool[T any](pool RoutinePool, fn func(ctx Context) (T, status.Status)) Routine[T] {
	r := newRoutine[T](fn)
	pool.Run(r)
	return r
}
