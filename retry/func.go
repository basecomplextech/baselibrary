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

// Func is a function without arguments but with the result.
type Func[T any] func(ctx async.Context) (T, status.Status)

// FuncRetrier retries a function and returns the result.
type FuncRetrier[T any] struct {
	retrier
	fn Func[T]
}

// Retry returns a function retrier.
//
// Example:
//
//	fn := func(ctx async.Context) (Result, status.Status) {
//	    // ...
//	}
//
//	result, st := retry.Retry(fn).
//		MaxRetries(5).
//		MinDelay(time.Second).
//		MaxDelay(10 * time.Second).
//		Error("operation failed").
//		ErrorHandler(myErrorHandler).
//		Run(ctx)
func Retry[T any](fn Func[T]) FuncRetrier[T] {
	return FuncRetrier[T]{
		retrier: newRetrier(),
		fn:      fn,
	}
}

var _ builder[FuncRetrier[any]] = (*FuncRetrier[any])(nil)

// Run retries the function and returns the result.
func (r FuncRetrier[T]) Run(ctx async.Context) (T, status.Status) {
	for attempt := 0; ; attempt++ {
		// Call function
		result, st := r.run(ctx)
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
func (r FuncRetrier[T]) Error(message string) FuncRetrier[T] {
	r.opts.Error = message
	return r
}

// ErrorFunc sets the error handler.
func (r FuncRetrier[T]) ErrorFunc(fn ErrorFunc) FuncRetrier[T] {
	r.opts.ErrorHandler = fn
	return r
}

// ErrorHandler sets the error handler.
func (r FuncRetrier[T]) ErrorHandler(handler ErrorHandler) FuncRetrier[T] {
	r.opts.ErrorHandler = handler
	return r
}

// Logger sets the default logger.
func (r FuncRetrier[T]) Logger(logger logging.Logger) FuncRetrier[T] {
	r.opts.Logger = logger
	return r
}

// MinDelay sets the min delay.
func (r FuncRetrier[T]) MinDelay(minDelay time.Duration) FuncRetrier[T] {
	r.opts.MinDelay = minDelay
	return r
}

// MaxDelay sets the max delay.
func (r FuncRetrier[T]) MaxDelay(maxDelay time.Duration) FuncRetrier[T] {
	r.opts.MaxDelay = maxDelay
	return r
}

// MaxRetries sets the max retries.
func (r FuncRetrier[T]) MaxRetries(maxRetries int) FuncRetrier[T] {
	r.opts.MaxRetries = maxRetries
	return r
}

// Options overrides all options.
func (r FuncRetrier[T]) Options(opts Options) FuncRetrier[T] {
	r.opts = opts
	return r
}

// private

func (r FuncRetrier[T]) run(ctx async.Context) (_ T, st status.Status) {
	defer func() {
		if e := recover(); e != nil {
			st = status.Recover(e)
		}
	}()

	return r.fn(ctx)
}
