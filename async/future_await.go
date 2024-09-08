// Copyright 2024 Ivan Korobkov. All rights reserved.

package async

import (
	"reflect"

	"github.com/basecomplextech/baselibrary/status"
)

// AwaitAll awaits all futures completion in a group, and returns the results.
//
// The method returns nil and the context status if the context is canceled.
func AwaitAll[F Future[T], T any](ctx Context, futures ...F) ([]Result[T], status.Status) {
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

// AwaitAny awaits the first future completion in a group, and returns the result.
//
// The method returns -1 and the context status if the context is canceled.
func AwaitAny[F Future[T], T any](ctx Context, futures ...F) (int, T, status.Status) {
	var zero T

	switch len(futures) {
	case 0:
		return -1, zero, status.OK
	case 1:
		f := futures[0]
		select {
		case <-f.Wait():
			result, st := f.Result()
			return 0, result, st
		case <-ctx.Wait():
			return -1, zero, ctx.Status()
		}
	}

	cases := make([]reflect.SelectCase, 0, len(futures)+1)

	// Add context channel
	{
		wait := ctx.Wait()
		case_ := reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(wait),
		}
		cases = append(cases, case_)
	}

	// Add future channels
	for _, f := range futures {
		wait := f.Wait()
		case_ := reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(wait),
		}
		cases = append(cases, case_)
	}

	// Select first ready
	j, _, _ := reflect.Select(cases)
	if j == 0 {
		return -1, zero, ctx.Status()
	}

	// Return future
	f := futures[j-1]
	result, st := f.Result()
	return j - 1, result, st
}

// AwaitError awaits and returns the first error in a group, or OK if all futures are successful.
//
// The method returns -1 and the context status if the context is canceled.
func AwaitError[F Future[T], T any](ctx Context, futures ...F) (int, status.Status) {
	switch len(futures) {
	case 0:
		return -1, status.OK
	case 1:
		f := futures[0]
		select {
		case <-f.Wait():
			_, st := f.Result()
			return 0, st
		case <-ctx.Wait():
			return -1, ctx.Status()
		}
	}

	// Track initial indexes
	type indexed struct {
		future Future[T]
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

		// Add context channel
		{
			wait := ctx.Wait()
			case_ := reflect.SelectCase{
				Dir:  reflect.SelectRecv,
				Chan: reflect.ValueOf(wait),
			}
			cases = append(cases, case_)
		}

		// Add future channels
		for _, f := range current {
			wait := f.future.Wait()
			case_ := reflect.SelectCase{
				Dir:  reflect.SelectRecv,
				Chan: reflect.ValueOf(wait),
			}
			cases = append(cases, case_)
		}

		// Select first ready
		j, _, _ := reflect.Select(cases)
		if j == 0 {
			return -1, ctx.Status()
		}

		// Check future result
		cur := current[j-1]
		_, st := cur.future.Result()
		if !st.OK() {
			return cur.index, st
		}

		// Remove future from current
		current = append(current[:j-1], current[j:]...)
	}

	return -1, status.OK
}
