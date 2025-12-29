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

// VoidFunc1 is a procedure which accepts one argument.
type VoidFunc1[A any] func(ctx async.Context, arg A) status.Status

// VoidRetrier1 retries a single argument procedure.
type VoidRetrier1[A any] struct {
	retrier
	fn VoidFunc1[A]
}

// RetryVoid1 returns a void function retrier.
//
// Example:
//
//	fn := func(ctx async.Context, arg ArgType) status.Status {
//	    // ...
//	}
//
//	st := retry.RetryVoid1(fn).
//		MaxRetries(5).
//		MinDelay(time.Second).
//		MaxDelay(10 * time.Second).
//		Error("operation failed").
//		ErrorHandler(myErrorHandler).
//		Run(ctx, arg)
func RetryVoid1[A any](fn VoidFunc1[A]) VoidRetrier1[A] {
	return VoidRetrier1[A]{
		retrier: newRetrier(),
		fn:      fn,
	}
}

var _ builder[VoidRetrier1[any]] = (*VoidRetrier1[any])(nil)

// Run retries the function.
func (r VoidRetrier1[A]) Run(ctx async.Context, arg A) status.Status {
	for attempt := 0; ; attempt++ {
		// Call function
		st := r.run(ctx, arg)
		switch st.Code {
		case status.CodeOK, status.CodeCancelled:
			return st
		}

		// Check max retries
		if r.opts.MaxRetries != 0 {
			if attempt >= r.opts.MaxRetries {
				return st
			}
		}

		// Handle error
		if st := r.handleError(st, attempt); !st.OK() {
			return st
		}

		// Sleep
		if st := r.sleep(ctx, attempt); !st.OK() {
			return st
		}
	}
}

// Error sets the error message.
func (r VoidRetrier1[A]) Error(message string) VoidRetrier1[A] {
	r.opts.Error = message
	return r
}

// ErrorFunc sets the error handler.
func (r VoidRetrier1[A]) ErrorFunc(fn ErrorFunc) VoidRetrier1[A] {
	r.opts.ErrorHandler = fn
	return r
}

// ErrorHandler sets the error handler.
func (r VoidRetrier1[A]) ErrorHandler(handler ErrorHandler) VoidRetrier1[A] {
	r.opts.ErrorHandler = handler
	return r
}

// Logger sets the default logger.
func (r VoidRetrier1[A]) Logger(logger logging.Logger) VoidRetrier1[A] {
	r.opts.Logger = logger
	return r
}

// MinDelay sets the min delay.
func (r VoidRetrier1[A]) MinDelay(minDelay time.Duration) VoidRetrier1[A] {
	r.opts.MinDelay = minDelay
	return r
}

// MaxDelay sets the max delay.
func (r VoidRetrier1[A]) MaxDelay(maxDelay time.Duration) VoidRetrier1[A] {
	r.opts.MaxDelay = maxDelay
	return r
}

// MaxRetries sets the max retries.
func (r VoidRetrier1[A]) MaxRetries(maxRetries int) VoidRetrier1[A] {
	r.opts.MaxRetries = maxRetries
	return r
}

// Options overrides all options.
func (r VoidRetrier1[A]) Options(opts Options) VoidRetrier1[A] {
	r.opts = opts
	return r
}

// private

func (r VoidRetrier1[A]) run(ctx async.Context, arg A) (st status.Status) {
	defer func() {
		if e := recover(); e != nil {
			st = status.Recover(e)
		}
	}()

	return r.fn(ctx, arg)
}
