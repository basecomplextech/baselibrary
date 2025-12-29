// Copyright 2025 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package retry

import (
	"time"

	"github.com/basecomplextech/baselibrary/async"
	"github.com/basecomplextech/baselibrary/logging"
	"github.com/basecomplextech/baselibrary/status"
)

// Func1 is a function with one argument and the result.
type Func1[T any, A any] func(ctx async.Context, arg A) (T, status.Status)

// Func1Retrier retries a function and returns the result.
type Func1Retrier[T any, A any] struct {
	retrier
	fn Func1[T, A]
}

// Retry1 returns a function retrier.
//
// Example:
//
//	fn := func(ctx async.Context, arg ArgType) (Result, status.Status) {
//	    // ...
//	}
//
//	result, st := retry.Retry1(fn).
//		MaxRetries(5).
//		MinDelay(time.Second).
//		MaxDelay(10 * time.Second).
//		Error("operation failed").
//		ErrorHandler(myErrorHandler).
//		Run(ctx, arg)
func Retry1[T any, A any](fn Func1[T, A]) Func1Retrier[T, A] {
	return Func1Retrier[T, A]{
		retrier: newRetrier(),
		fn:      fn,
	}
}

var _ builder[Func1Retrier[any, any]] = (*Func1Retrier[any, any])(nil)

// Run retries the function and returns the result.
func (r Func1Retrier[T, A]) Run(ctx async.Context, arg A) (T, status.Status) {
	for attempt := 0; ; attempt++ {
		// Call function
		result, st := r.run(ctx, arg)
		switch st.Code {
		case status.CodeOK, status.CodeCancelled:
			return result, st
		}

		// Check max retries
		if r.opts.MaxRetries != 0 {
			if attempt >= r.opts.MaxRetries {
				return result, st
			}
		}

		// Handle error
		if st := r.handleError(st, attempt); !st.OK() {
			return result, st
		}

		// Sleep
		if st := r.sleep(ctx, attempt); !st.OK() {
			return result, st
		}
	}
}

// Error sets the error message.
func (r Func1Retrier[T, A]) Error(message string) Func1Retrier[T, A] {
	r.opts.Error = message
	return r
}

// ErrorFunc sets the error handler.
func (r Func1Retrier[T, A]) ErrorFunc(fn ErrorFunc) Func1Retrier[T, A] {
	r.opts.ErrorHandler = fn
	return r
}

// ErrorHandler sets the error handler.
func (r Func1Retrier[T, A]) ErrorHandler(handler ErrorHandler) Func1Retrier[T, A] {
	r.opts.ErrorHandler = handler
	return r
}

// Logger sets the default logger.
func (r Func1Retrier[T, A]) Logger(logger logging.Logger) Func1Retrier[T, A] {
	r.opts.Logger = logger
	return r
}

// MinDelay sets the min delay.
func (r Func1Retrier[T, A]) MinDelay(minDelay time.Duration) Func1Retrier[T, A] {
	r.opts.MinDelay = minDelay
	return r
}

// MaxDelay sets the max delay.
func (r Func1Retrier[T, A]) MaxDelay(maxDelay time.Duration) Func1Retrier[T, A] {
	r.opts.MaxDelay = maxDelay
	return r
}

// MaxRetries sets the max retries.
func (r Func1Retrier[T, A]) MaxRetries(maxRetries int) Func1Retrier[T, A] {
	r.opts.MaxRetries = maxRetries
	return r
}

// Options overrides all options.
func (r Func1Retrier[T, A]) Options(opts Options) Func1Retrier[T, A] {
	r.opts = opts
	return r
}

// private

func (r Func1Retrier[T, A]) run(ctx async.Context, arg A) (_ T, st status.Status) {
	defer func() {
		if e := recover(); e != nil {
			st = status.Recover(e)
		}
	}()

	return r.fn(ctx, arg)
}
