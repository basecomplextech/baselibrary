// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package retry

import (
	"github.com/basecomplextech/baselibrary/async"
	"github.com/basecomplextech/baselibrary/status"
)

// Retry calls a function, retries on errors and panics, uses exponential backoff.
func Retry(ctx async.Context, fn func(ctx async.Context) status.Status) status.Status {
	r := newRetrier(Options{})

	return r.Retry(ctx, func(ctx async.Context) status.Status {
		return fn(ctx)
	})
}

// Retry1 calls a function, retries on errors and panics, uses exponential backoff.
func Retry1[T any](ctx async.Context, fn func(ctx async.Context) (T, status.Status)) (
	result T, st status.Status) {
	r := newRetrier(Options{})

	st = r.Retry(ctx, func(ctx async.Context) status.Status {
		result, st = fn(ctx)
		return st
	})
	return result, st
}
