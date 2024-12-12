// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package async

import (
	"reflect"

	"github.com/basecomplextech/baselibrary/status"
)

// FutureGroup is a group of futures of the same type.
// Use [FutureGroupDyn] for a group of futures of different types.
type FutureGroup[T any] []Future[T]

// Await waits for the completion of all futures in the group.
// The method returns the context status if the context is cancelled.
func (g FutureGroup[T]) Await(ctx Context) status.Status {
	return awaitAll(ctx, g...)
}

// AwaitAny waits for the completion of any future in the group, and returns its result.
// The method returns -1 and the context status if the context is cancelled.
func (g FutureGroup[T]) AwaitAny(ctx Context) (T, int, status.Status) {
	return awaitAny(ctx, g...)
}

// AwaitError waits for any failure of the futures in the group, and returns the error.
// The method returns ok if all futures are successful.
func (g FutureGroup[T]) AwaitError(ctx Context) (int, status.Status) {
	return awaitError(ctx, g...)
}

// AwaitResults waits for the completion of all futures in the group, and returns the results.
// The method returns nil and the context status if the context is cancelled.
func (g FutureGroup[T]) AwaitResults(ctx Context) ([]Result[T], status.Status) {
	return awaitResults(ctx, g...)
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

// internal

// awaitAll waits for the completion of all futures.
// The method the context status if the context is cancelled.
func awaitAll[F FutureDyn](ctx Context, futures ...F) status.Status {
	for _, f := range futures {
		select {
		case <-f.Wait():
		case <-ctx.Wait():
			return ctx.Status()
		}
	}
	return status.OK
}

// awaitResults waits for the completion of all futures, and returns the results.
// The method returns nil and the context status if the context is cancelled.
func awaitResults[F Future[T], T any](ctx Context, futures ...F) ([]Result[T], status.Status) {
	results := make([]Result[T], 0, len(futures))

	for _, f := range futures {
		select {
		case <-f.Wait():
			value, st := f.Result()
			result := Result[T]{value, st}
			results = append(results, result)

		case <-ctx.Wait():
			return nil, ctx.Status()
		}
	}

	return results, status.OK
}

// awaitAny waits for the completion of any future, and returns its result.
// The method returns -1 and the context status if the context is cancelled.
func awaitAny[F Future[T], T any](ctx Context, futures ...F) (T, int, status.Status) {
	var zero T

	// Special cases
	switch len(futures) {
	case 0:
		return zero, -1, status.OK
	case 1:
		f := futures[0]
		select {
		case <-f.Wait():
			result, st := f.Result()
			return result, 0, st
		case <-ctx.Wait():
			return zero, -1, ctx.Status()
		}
	}

	// Make select cases
	cases := make([]reflect.SelectCase, 0, len(futures)+1)

	// Add context case
	{
		wait := ctx.Wait()
		case_ := reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(wait),
		}
		cases = append(cases, case_)
	}

	// Add future cases
	for _, f := range futures {
		wait := f.Wait()
		case_ := reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(wait),
		}
		cases = append(cases, case_)
	}

	// Await any case
	j, _, _ := reflect.Select(cases)
	if j == 0 {
		return zero, -1, ctx.Status()
	}

	// Return future result
	f := futures[j-1]
	result, st := f.Result()
	return result, j - 1, st
}

// awaitError awaits failure of any future, and returns its error.
// The method returns -1 and the context status if the context is cancelled.
func awaitError[F FutureDyn](ctx Context, futures ...F) (int, status.Status) {
	// Special cases
	switch len(futures) {
	case 0:
		return -1, status.OK
	case 1:
		f := futures[0]
		select {
		case <-f.Wait():
			st := f.Status()
			return 0, st
		case <-ctx.Wait():
			return -1, ctx.Status()
		}
	}

	// Track initial indexes
	type indexed struct {
		future FutureDyn
		index  int
	}
	current := make([]indexed, len(futures))
	for i, f := range futures {
		current[i] = indexed{f, i}
	}

	// Await error or all completion
	cases := make([]reflect.SelectCase, 0, len(current)+1)
	for len(current) > 0 {
		clear(cases)
		cases = cases[:0]

		// Add context case
		{
			wait := ctx.Wait()
			case_ := reflect.SelectCase{
				Dir:  reflect.SelectRecv,
				Chan: reflect.ValueOf(wait),
			}
			cases = append(cases, case_)
		}

		// Add future cases
		for _, f := range current {
			wait := f.future.Wait()
			case_ := reflect.SelectCase{
				Dir:  reflect.SelectRecv,
				Chan: reflect.ValueOf(wait),
			}
			cases = append(cases, case_)
		}

		// Await any case
		j, _, _ := reflect.Select(cases)
		if j == 0 {
			return -1, ctx.Status()
		}

		// Check future result
		cur := current[j-1]
		st := cur.future.Status()
		if !st.OK() {
			return cur.index, st
		}

		// Remove future from current
		current = append(current[:j-1], current[j:]...)
	}

	return -1, status.OK
}
