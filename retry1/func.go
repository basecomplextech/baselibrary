// Copyright 2025 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package retry1

import (
	"time"

	"github.com/basecomplextech/baselibrary/async"
	"github.com/basecomplextech/baselibrary/logging"
	"github.com/basecomplextech/baselibrary/status"
)

type (
	// Func is a function that returns the result.
	Func[T any] func(ctx async.Context) (T, status.Status)

	// Func1 is a function that accepts one argument and returns the result.
	Func1[T any, A any] func(ctx async.Context, arg A) (T, status.Status)
)

type (
	// FuncCall retries a function and returns the result.
	FuncCall[T any] struct {
		call
		fn Func[T]
	}

	// Func1Call retries a function and returns the result.
	Func1Call[T any, A any] struct {
		call
		fn Func1[T, A]
	}
)

// Retry

// RetryFunc returns a function call.
func RetryFunc[T any](fn Func[T]) FuncCall[T] {
	return FuncCall[T]{
		call: newCall(),
		fn:   fn,
	}
}

// RetryFunc1 returns a function call.
func RetryFunc1[T any, A any](fn Func1[T, A]) Func1Call[T, A] {
	return Func1Call[T, A]{
		call: newCall(),
		fn:   fn,
	}
}

// Run

// Run retries the function and returns the result.
func (c FuncCall[T]) Run(ctx async.Context) (T, status.Status) {
	for attempt := 0; ; attempt++ {
		// Call function
		result, st := c.run(ctx)
		switch st.Code {
		case status.CodeOK, status.CodeCancelled:
			return result, st
		}

		// Log error
		c.logError(st, attempt)

		// Sleep
		if st := c.sleep(ctx, attempt); !st.OK() {
			return result, st
		}
	}
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

		// Log error
		c.logError(st, attempt)

		// Sleep
		if st := c.sleep(ctx, attempt); !st.OK() {
			return result, st
		}
	}
}

// FuncCall

var _ builder[FuncCall[any]] = (*FuncCall[any])(nil)

// Error sets the error message.
func (c FuncCall[T]) Error(message string) FuncCall[T] {
	c.opts.Error = message
	return c
}

// ErrorLogger sets the error logger.
func (c FuncCall[T]) ErrorLogger(logger ErrorLogger) FuncCall[T] {
	c.opts.ErrorLogger = logger
	return c
}

// Logger sets the default logger.
func (c FuncCall[T]) Logger(logger logging.Logger) FuncCall[T] {
	c.opts.Logger = logger
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

// Func1Call

var _ builder[Func1Call[any, any]] = (*Func1Call[any, any])(nil)

// Error sets the error message.
func (c Func1Call[T, A]) Error(message string) Func1Call[T, A] {
	c.opts.Error = message
	return c
}

// ErrorLogger sets the error logger.
func (c Func1Call[T, A]) ErrorLogger(logger ErrorLogger) Func1Call[T, A] {
	c.opts.ErrorLogger = logger
	return c
}

// Logger sets the default logger.
func (c Func1Call[T, A]) Logger(logger logging.Logger) Func1Call[T, A] {
	c.opts.Logger = logger
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

// private

func (c FuncCall[T]) run(ctx async.Context) (_ T, st status.Status) {
	defer func() {
		if e := recover(); e != nil {
			st = status.Recover(e)
		}
	}()

	return c.fn(ctx)
}

func (c Func1Call[T, A]) run(ctx async.Context, arg A) (_ T, st status.Status) {
	defer func() {
		if e := recover(); e != nil {
			st = status.Recover(e)
		}
	}()

	return c.fn(ctx, arg)
}
