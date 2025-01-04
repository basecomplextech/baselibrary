// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package future

import (
	"github.com/basecomplextech/baselibrary/async/internal/context"
	"github.com/basecomplextech/baselibrary/status"
)

// FutureGroup is a group of futures of the same type.
// Use [FutureGroupDyn] for a group of futures of different types.
type FutureGroup[T any] []Future[T]

// Await waits for the completion of all futures in the group.
// The method returns the context status if the context is cancelled.
func (g FutureGroup[T]) Await(ctx context.Context) status.Status {
	return AwaitAll(ctx, g...)
}

// AwaitAny waits for the completion of any future in the group, and returns its result.
// The method returns -1 and the context status if the context is cancelled.
func (g FutureGroup[T]) AwaitAny(ctx context.Context) (T, int, status.Status) {
	return AwaitAny(ctx, g...)
}

// AwaitError waits for any failure of the futures in the group, and returns the error.
// The method returns ok if all futures are successful.
func (g FutureGroup[T]) AwaitError(ctx context.Context) (int, status.Status) {
	return AwaitError(ctx, g...)
}

// AwaitResults waits for the completion of all futures in the group, and returns the results.
// The method returns nil and the context status if the context is cancelled.
func (g FutureGroup[T]) AwaitResults(ctx context.Context) ([]Result[T], status.Status) {
	return AwaitResults(ctx, g...)
}

// Results returns the results of all futures in the group.
func (g FutureGroup[T]) Results() []Result[T] {
	results := make([]Result[T], 0, len(g))
	for _, f := range g {
		value, st := f.Result()
		result := Result[T]{value, st}
		results = append(results, result)
	}
	return results
}

// Statuses returns the statuses of all futures in the group.
func (g FutureGroup[T]) Statuses() []status.Status {
	result := make([]status.Status, 0, len(g))
	for _, f := range g {
		st := f.Status()
		result = append(result, st)
	}
	return result
}
