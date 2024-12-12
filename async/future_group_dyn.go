// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package async

import (
	"reflect"

	"github.com/basecomplextech/baselibrary/status"
)

// FutureGroupDyn is a group of futures of different types.
// Use [FutureGroup] for a group of futures of the same type.
type FutureGroupDyn []FutureDyn

// Await waits for the completion of all futures in the group.
// The method returns nil and the context status if the context is cancelled.
func (g FutureGroupDyn) Await(ctx Context) status.Status {
	return awaitAll(ctx, g...)
}

// AwaitAny waits for the completion of any future in the group, and returns its status.
// The method returns -1 and the context status if the context is cancelled.
func (g FutureGroupDyn) AwaitAny(ctx Context) (int, status.Status) {
	return awaitAnyDyn(ctx, g...)
}

// AwaitError waits for any failure of the futures in the group, and returns the error.
// The method returns ok if all futures are successful.
func (g FutureGroupDyn) AwaitError(ctx Context) (int, status.Status) {
	return awaitError(ctx, g...)
}

// Statuses returns the statuses of all futures in the group.
func (g FutureGroupDyn) Statuses() []status.Status {
	result := make([]status.Status, 0, len(g))
	for _, f := range g {
		st := f.Status()
		result = append(result, st)
	}
	return result
}

// internal

// awaitAnyDyn awaits completion of any future, and returns its result.
// The method returns -1 and the context status if the context is cancelled.
func awaitAnyDyn[F FutureDyn](ctx Context, futures ...F) (int, status.Status) {
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
		return -1, ctx.Status()
	}

	// Return future result
	f := futures[j-1]
	st := f.Status()
	return j - 1, st
}
