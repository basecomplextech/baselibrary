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
	// VoidFunc1 is a procedure which accepts one argument.
	VoidFunc1[A any] func(ctx async.Context, arg A) status.Status

	// VoidCall1 retries a single argument procedure.
	VoidCall1[A any] struct {
		call
		proc VoidFunc1[A]
	}
)

// RetryVoid1 returns a new procedure call.
func RetryVoid1[A any](fn VoidFunc1[A]) VoidCall1[A] {
	return VoidCall1[A]{proc: fn}
}

var _ builder[VoidCall1[any]] = (*VoidCall1[any])(nil)

// Error sets the error message.
func (c VoidCall1[A]) Error(message string) VoidCall1[A] {
	c.opts.Error = message
	return c
}

// ErrorFunc sets the error function.
func (c VoidCall1[A]) ErrorFunc(fn ErrorFunc) VoidCall1[A] {
	c.opts.ErrorLogger = errorLoggerFunc(fn)
	return c
}

// ErrorLogger sets the error logger.
func (c VoidCall1[A]) ErrorLogger(logger ErrorLogger) VoidCall1[A] {
	c.opts.ErrorLogger = logger
	return c
}

// Logger sets the default logger.
func (c VoidCall1[A]) Logger(logger logging.Logger) VoidCall1[A] {
	c.opts.Logger = logger
	return c
}

// MinDelay sets the min delay.
func (c VoidCall1[A]) MinDelay(minDelay time.Duration) VoidCall1[A] {
	c.opts.MinDelay = minDelay
	return c
}

// MaxDelay sets the max delay.
func (c VoidCall1[A]) MaxDelay(maxDelay time.Duration) VoidCall1[A] {
	c.opts.MaxDelay = maxDelay
	return c
}

// MaxRetries sets the max retries.
func (c VoidCall1[A]) MaxRetries(maxRetries int) VoidCall1[A] {
	c.opts.MaxRetries = maxRetries
	return c
}

// Options overrides all options.
func (c VoidCall1[A]) Options(opts Options) VoidCall1[A] {
	c.opts = opts
	return c
}

// Run retries the procedure.
func (c VoidCall1[A]) Run(ctx async.Context, arg A) status.Status {
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

// private

func (c VoidCall1[A]) run(ctx async.Context, arg A) (st status.Status) {
	defer func() {
		if e := recover(); e != nil {
			st = status.Recover(e)
		}
	}()

	return c.proc(ctx, arg)
}
