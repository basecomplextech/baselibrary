// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package async

import (
	"github.com/basecomplextech/baselibrary/status"
)

// RoutineGroup is a group of routines of the same type.
// Use [RoutineGroupDyn] for a group of routines of different types.
type RoutineGroup[T any] []Routine[T]

// Await waits for the completion of all routines in the group.
// The method returns the context status if the context is cancelled.
func (g RoutineGroup[T]) Await(ctx Context) status.Status {
	return awaitAll(ctx, g...)
}

// AwaitAny waits for the completion of any routine in the group, and returns its result.
// The method returns -1 and the context status if the context is cancelled.
func (g RoutineGroup[T]) AwaitAny(ctx Context) (T, int, status.Status) {
	return awaitAny(ctx, g...)
}

// AwaitError waits for any failure of the routines in the group, and returns the error.
// The method returns ok if all routines are successful.
func (g RoutineGroup[T]) AwaitError(ctx Context) (int, status.Status) {
	return awaitError(ctx, g...)
}

// AwaitResults waits for the completion of all routines in the group, and returns the results.
// The method returns nil and the context status if the context is cancelled.
func (g RoutineGroup[T]) AwaitResults(ctx Context) ([]Result[T], status.Status) {
	return awaitResults(ctx, g...)
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
func (g RoutineGroup[T]) Results() []Result[T] {
	results := make([]Result[T], 0, len(g))
	for _, r := range g {
		value, st := r.Result()
		result := Result[T]{value, st}
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
	StopAll(g...)
}

// StopWait stops and awaits all routines in the group.
func (g RoutineGroup[T]) StopWait() {
	StopWaitAll(g...)
}
