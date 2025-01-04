// Copyright 2025 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package future

import (
	"reflect"

	"github.com/basecomplextech/baselibrary/async/internal/context"
	"github.com/basecomplextech/baselibrary/status"
)

// AwaitAnyDyn awaits completion of any future, and returns its result.
// The method returns -1 and the context status if the context is cancelled.
func AwaitAnyDyn[F FutureDyn](ctx context.Context, futures ...F) (int, status.Status) {
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
