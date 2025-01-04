// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package routine

import (
	"github.com/basecomplextech/baselibrary/async/internal"
	"github.com/basecomplextech/baselibrary/async/internal/context"
	"github.com/basecomplextech/baselibrary/async/internal/future"
	"github.com/basecomplextech/baselibrary/status"
)

// RoutineGroup is a group of routines of the same type.
// Use [RoutineGroupDyn] for a group of routines of different types.
type RoutineGroup[T any] []Routine[T]

// Await waits for the completion of all routines in the group.
// The method returns the context status if the context is cancelled.
func (g RoutineGroup[T]) Await(ctx context.Context) status.Status {
	return future.AwaitAll(ctx, g...)
}

// AwaitAny waits for the completion of any routine in the group, and returns its result.
// The method returns -1 and the context status if the context is cancelled.
func (g RoutineGroup[T]) AwaitAny(ctx context.Context) (T, int, status.Status) {
	return future.AwaitAny(ctx, g...)
}

// AwaitError waits for any failure of the routines in the group, and returns the error.
// The method returns ok if all routines are successful.
func (g RoutineGroup[T]) AwaitError(ctx context.Context) (int, status.Status) {
	return future.AwaitError(ctx, g...)
}

// AwaitResults waits for the completion of all routines in the group, and returns the results.
// The method returns nil and the context status if the context is cancelled.
func (g RoutineGroup[T]) AwaitResults(ctx context.Context) ([]future.Result[T], status.Status) {
	return future.AwaitResults(ctx, g...)
}

// Values returns the values of all results in the group.
func (g RoutineGroup[T]) Values() []T {
	values := make([]T, 0, len(g))
	for _, r := range g {
		value, _ := r.Result()
		values = append(values, value)
	}
	return values
}

// Results returns the results of all routines in the group.
func (g RoutineGroup[T]) Results() []future.Result[T] {
	results := make([]future.Result[T], 0, len(g))
	for _, r := range g {
		value, st := r.Result()
		result := future.Result[T]{value, st}
		results = append(results, result)
	}
	return results
}

// Statuses returns the statuses of all routines in the group.
func (g RoutineGroup[T]) Statuses() []status.Status {
	result := make([]status.Status, 0, len(g))
	for _, r := range g {
		st := r.Status()
		result = append(result, st)
	}
	return result
}

// Stop stops all routines in the group.
func (g RoutineGroup[T]) Stop() {
	internal.StopAll(g...)
}

// StopWait stops and awaits all routines in the group.
func (g RoutineGroup[T]) StopWait() {
	internal.StopWaitAll(g...)
}
