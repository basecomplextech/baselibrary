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

// LoopFunc1 is a function with one argument that is retried in a loop.
type LoopFunc1[A any] func(ctx async.Context, arg A, success *bool) status.Status

// Loop1Retrier retries a single argument function in a loop.
type Loop1Retrier[A any] struct {
	retrier
	fn LoopFunc1[A]
}

// RetryLoop1 returns a loop retrier.
//
// Example:
//
//	fn := func(ctx async.Context, arg ArgType, success *bool) status.Status {
//	    // ...
//	}
//
//	st := retry.RetryLoop1(fn).
//		MaxRetries(5).
//		MinDelay(time.Second).
//		MaxDelay(10 * time.Second).
//		Error("operation failed").
//		ErrorHandler(myErrorHandler).
//		Run(ctx, arg)
func RetryLoop1[A any](fn LoopFunc1[A]) Loop1Retrier[A] {
	return Loop1Retrier[A]{
		retrier: newRetrier(),
		fn:      fn,
	}
}

var _ builder[Loop1Retrier[any]] = (*Loop1Retrier[any])(nil)

// Run retries the function in a loop.
func (r Loop1Retrier[A]) Run(ctx async.Context, arg A) status.Status {
	success := new(bool)

	for attempt := 0; ; attempt++ {
		// Restart on success
		if *success {
			attempt = 0
			*success = false
		}

		// Call function
		st := r.run(ctx, arg, success)
		if st.Code == status.CodeCancelled {
			return st
		}

		// Handler error
		if !st.OK() {
			if st := r.handleError(st, attempt); !st.OK() {
				return st
			}
		}

		// Sleep before retry
		if st := r.sleep(ctx, attempt); !st.OK() {
			return st
		}
	}
}

// Error sets the error message.
func (r Loop1Retrier[A]) Error(message string) Loop1Retrier[A] {
	r.opts.Error = message
	return r
}

// ErrorFunc sets the error handler.
func (r Loop1Retrier[A]) ErrorFunc(fn ErrorFunc) Loop1Retrier[A] {
	r.opts.ErrorHandler = fn
	return r
}

// ErrorHandler sets the error handler.
func (r Loop1Retrier[A]) ErrorHandler(handler ErrorHandler) Loop1Retrier[A] {
	r.opts.ErrorHandler = handler
	return r
}

// Logger sets the default logger.
func (r Loop1Retrier[A]) Logger(logger logging.Logger) Loop1Retrier[A] {
	r.opts.Logger = logger
	return r
}

// MinDelay sets the min delay.
func (r Loop1Retrier[A]) MinDelay(minDelay time.Duration) Loop1Retrier[A] {
	r.opts.MinDelay = minDelay
	return r
}

// MaxDelay sets the max delay.
func (r Loop1Retrier[A]) MaxDelay(maxDelay time.Duration) Loop1Retrier[A] {
	r.opts.MaxDelay = maxDelay
	return r
}

// MaxRetries sets the max retries.
func (r Loop1Retrier[A]) MaxRetries(maxRetries int) Loop1Retrier[A] {
	r.opts.MaxRetries = maxRetries
	return r
}

// Options overrides all options.
func (r Loop1Retrier[A]) Options(opts Options) Loop1Retrier[A] {
	r.opts = opts
	return r
}

// private

func (r Loop1Retrier[A]) run(ctx async.Context, arg A, success *bool) (st status.Status) {
	defer func() {
		if e := recover(); e != nil {
			st = status.Recover(e)
		}
	}()

	return r.fn(ctx, arg, success)
}
