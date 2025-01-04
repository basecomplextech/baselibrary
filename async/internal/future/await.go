// Copyright 2025 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package future

import (
	"reflect"

	"github.com/basecomplextech/baselibrary/async/internal/context"
	"github.com/basecomplextech/baselibrary/status"
)

// AwaitAll waits for the completion of all futures.
// The method the context status if the context is cancelled.
func AwaitAll[F FutureDyn](ctx context.Context, futures ...F) status.Status {
	for _, f := range futures {
		select {
		case <-f.Wait():
		case <-ctx.Wait():
			return ctx.Status()
		}
	}
	return status.OK
}

// AwaitResults waits for the completion of all futures, and returns the results.
// The method returns nil and the context status if the context is cancelled.
func AwaitResults[F Future[T], T any](ctx context.Context, futures ...F) ([]Result[T], status.Status) {
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

// AwaitAny waits for the completion of any future, and returns its result.
// The method returns -1 and the context status if the context is cancelled.
func AwaitAny[F Future[T], T any](ctx context.Context, futures ...F) (T, int, status.Status) {
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

// AwaitError awaits failure of any future, and returns its error.
// The method returns -1 and the context status if the context is cancelled.
func AwaitError[F FutureDyn](ctx context.Context, futures ...F) (int, status.Status) {
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
