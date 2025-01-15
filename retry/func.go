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
	// Func is a function that returns the result.
	Func[T any] func(ctx async.Context) (T, status.Status)

	// FuncCall retries a function and returns the result.
	FuncCall[T any] struct {
		call
		fn Func[T]
	}
)

// Retry returns a function call.
func Retry[T any](fn Func[T]) FuncCall[T] {
	return FuncCall[T]{
		call: newCall(),
		fn:   fn,
	}
}

var _ builder[FuncCall[any]] = (*FuncCall[any])(nil)

// Error sets the error message.
func (c FuncCall[T]) Error(message string) FuncCall[T] {
	c.opts.Error = message
	return c
}

// ErrorFunc sets the error handler.
func (c FuncCall[T]) ErrorFunc(fn ErrorFunc) FuncCall[T] {
	c.opts.ErrorHandler = fn
	return c
}

// ErrorHandler sets the error handler.
func (c FuncCall[T]) ErrorHandler(handler ErrorHandler) FuncCall[T] {
	c.opts.ErrorHandler = handler
	return c
}

// Logger sets the default logger.
func (c FuncCall[T]) Logger(logger logging.Logger) FuncCall[T] {
	c.opts.Logger = logger
	return c
}

// MinDelay sets the min delay.
func (c FuncCall[T]) MinDelay(minDelay time.Duration) FuncCall[T] {
	c.opts.MinDelay = minDelay
	return c
}

// MaxDelay sets the max delay.
func (c FuncCall[T]) MaxDelay(maxDelay time.Duration) FuncCall[T] {
	c.opts.MaxDelay = maxDelay
	return c
}

// MaxRetries sets the max retries.
func (c FuncCall[T]) MaxRetries(maxRetries int) FuncCall[T] {
	c.opts.MaxRetries = maxRetries
	return c
}

// Options overrides all options.
func (c FuncCall[T]) Options(opts Options) FuncCall[T] {
	c.opts = opts
	return c
}

// Run retries the function and returns the result.
func (c FuncCall[T]) Run(ctx async.Context) (T, status.Status) {
	for attempt := 0; ; attempt++ {
		// Call function
		result, st := c.run(ctx)
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

func (c FuncCall[T]) run(ctx async.Context) (_ T, st status.Status) {
	defer func() {
		if e := recover(); e != nil {
			st = status.Recover(e)
		}
	}()

	return c.fn(ctx)
}
