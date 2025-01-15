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

type (
	// Func1 is a function that accepts one argument and returns the result.
	Func1[T any, A any] func(ctx async.Context, arg A) (T, status.Status)

	// Func1Call retries a function and returns the result.
	Func1Call[T any, A any] struct {
		call
		fn Func1[T, A]
	}
)

// Retry1 returns a function call.
func Retry1[T any, A any](fn Func1[T, A]) Func1Call[T, A] {
	return Func1Call[T, A]{
		call: newCall(),
		fn:   fn,
	}
}

var _ builder[Func1Call[any, any]] = (*Func1Call[any, any])(nil)

// Error sets the error message.
func (c Func1Call[T, A]) Error(message string) Func1Call[T, A] {
	c.opts.Error = message
	return c
}

// ErrorFunc sets the error handler.
func (c Func1Call[T, A]) ErrorFunc(fn ErrorFunc) Func1Call[T, A] {
	c.opts.ErrorHandler = fn
	return c
}

// ErrorHandler sets the error handler.
func (c Func1Call[T, A]) ErrorHandler(handler ErrorHandler) Func1Call[T, A] {
	c.opts.ErrorHandler = handler
	return c
}

// Logger sets the default logger.
func (c Func1Call[T, A]) Logger(logger logging.Logger) Func1Call[T, A] {
	c.opts.Logger = logger
	return c
}

// MinDelay sets the min delay.
func (c Func1Call[T, A]) MinDelay(minDelay time.Duration) Func1Call[T, A] {
	c.opts.MinDelay = minDelay
	return c
}

// MaxDelay sets the max delay.
func (c Func1Call[T, A]) MaxDelay(maxDelay time.Duration) Func1Call[T, A] {
	c.opts.MaxDelay = maxDelay
	return c
}

// MaxRetries sets the max retries.
func (c Func1Call[T, A]) MaxRetries(maxRetries int) Func1Call[T, A] {
	c.opts.MaxRetries = maxRetries
	return c
}

// Options overrides all options.
func (c Func1Call[T, A]) Options(opts Options) Func1Call[T, A] {
	c.opts = opts
	return c
}

// Run retries the function and returns the result.
func (c Func1Call[T, A]) Run(ctx async.Context, arg A) (T, status.Status) {
	for attempt := 0; ; attempt++ {
		// Call function
		result, st := c.run(ctx, arg)
		switch st.Code {
		case status.CodeOK, status.CodeCancelled:
			return result, st
		}

		// Check max retries
		if c.opts.MaxRetries != 0 {
			if attempt >= c.opts.MaxRetries {
				return result, st
			}
		}

		// Handle error
		if st := c.handleError(st, attempt); !st.OK() {
			return result, st
		}

		// Sleep
		if st := c.sleep(ctx, attempt); !st.OK() {
			return result, st
		}
	}
}

// private

func (c Func1Call[T, A]) run(ctx async.Context, arg A) (_ T, st status.Status) {
	defer func() {
		if e := recover(); e != nil {
			st = status.Recover(e)
		}
	}()

	return c.fn(ctx, arg)
}
