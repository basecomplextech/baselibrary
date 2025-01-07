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
	// Proc is a procedure without arguments and results.
	Proc func(ctx async.Context) status.Status

	// Proc1 is a procedure which accepts one argument.
	Proc1[A any] func(ctx async.Context, arg A) status.Status
)

type (
	// Call retries a procedure, i.e. a function without results.
	Call struct {
		call
		proc Proc
	}

	// Call1 retries a single argument procedure.
	Call1[A any] struct {
		call
		proc Proc1[A]
	}
)

// Retry

// Retry returns a new procedure call.
func Retry(fn Proc) Call {
	return Call{proc: fn}
}

// Retry1 returns a new procedure call.
func Retry1[A any](fn Proc1[A]) Call1[A] {
	return Call1[A]{proc: fn}
}

// Run

// Run retries the procedure.
func (c Call) Run(ctx async.Context) status.Status {
	for attempt := 0; ; attempt++ {
		// Call function
		st := c.run(ctx)
		switch st.Code {
		case status.CodeOK, status.CodeCancelled:
			return st
		}

		// Check max retries
		if c.opts.MaxRetries != 0 {
			if attempt >= c.opts.MaxRetries {
				return st
			}
		}

		// Log error
		c.logError(st, attempt)

		// Sleep
		if st := c.sleep(ctx, attempt); !st.OK() {
			return st
		}
	}
}

// Run retries the procedure.
func (c Call1[A]) Run(ctx async.Context, arg A) status.Status {
	for attempt := 0; ; attempt++ {
		// Call function
		st := c.run(ctx, arg)
		switch st.Code {
		case status.CodeOK, status.CodeCancelled:
			return st
		}

		// Check max retries
		if c.opts.MaxRetries != 0 {
			if attempt >= c.opts.MaxRetries {
				return st
			}
		}

		// Log error
		c.logError(st, attempt)

		// Sleep
		if st := c.sleep(ctx, attempt); !st.OK() {
			return st
		}
	}
}

// Call

var _ builder[Call] = (*Call)(nil)

// Error sets the error message.
func (c Call) Error(message string) Call {
	c.opts.Error = message
	return c
}

// ErrorLogger sets the error logger.
func (c Call) ErrorLogger(logger ErrorLogger) Call {
	c.opts.ErrorLogger = logger
	return c
}

// Logger sets the default logger.
func (c Call) Logger(logger logging.Logger) Call {
	c.opts.Logger = logger
	return c
}

// MaxDelay sets the max delay.
func (c Call) MaxDelay(maxDelay time.Duration) Call {
	c.opts.MaxDelay = maxDelay
	return c
}

// MaxRetries sets the max retries.
func (c Call) MaxRetries(maxRetries int) Call {
	c.opts.MaxRetries = maxRetries
	return c
}

// Options overrides all options.
func (c Call) Options(opts Options) Call {
	c.opts = opts
	return c
}

// Call1

var _ builder[Call1[any]] = (*Call1[any])(nil)

// Error sets the error message.
func (c Call1[A]) Error(message string) Call1[A] {
	c.opts.Error = message
	return c
}

// ErrorLogger sets the error logger.
func (c Call1[A]) ErrorLogger(logger ErrorLogger) Call1[A] {
	c.opts.ErrorLogger = logger
	return c
}

// Logger sets the default logger.
func (c Call1[A]) Logger(logger logging.Logger) Call1[A] {
	c.opts.Logger = logger
	return c
}

// MaxDelay sets the max delay.
func (c Call1[A]) MaxDelay(maxDelay time.Duration) Call1[A] {
	c.opts.MaxDelay = maxDelay
	return c
}

// MaxRetries sets the max retries.
func (c Call1[A]) MaxRetries(maxRetries int) Call1[A] {
	c.opts.MaxRetries = maxRetries
	return c
}

// Options overrides all options.
func (c Call1[A]) Options(opts Options) Call1[A] {
	c.opts = opts
	return c
}

// private

func (c Call) run(ctx async.Context) (st status.Status) {
	defer func() {
		if e := recover(); e != nil {
			st = status.Recover(e)
		}
	}()

	return c.proc(ctx)
}

func (c Call1[A]) run(ctx async.Context, arg A) (st status.Status) {
	defer func() {
		if e := recover(); e != nil {
			st = status.Recover(e)
		}
	}()

	return c.proc(ctx, arg)
}
