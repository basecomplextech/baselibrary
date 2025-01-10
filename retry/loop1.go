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
	// LoopFunc1 is a loop function with one argument.
	LoopFunc1[A any] func(ctx async.Context, arg A, success *bool) status.Status

	// Loop1Call retries a single argument function in a loop.
	Loop1Call[A any] struct {
		call
		fn LoopFunc1[A]
	}
)

// RetryLoop1 returns a loop call.
func RetryLoop1[A any](fn LoopFunc1[A]) Loop1Call[A] {
	return Loop1Call[A]{
		call: newCall(),
		fn:   fn,
	}
}

var _ builder[Loop1Call[any]] = (*Loop1Call[any])(nil)

// Error sets the error message.
func (c Loop1Call[A]) Error(message string) Loop1Call[A] {
	c.opts.Error = message
	return c
}

// ErrorFunc sets the error function.
func (c Loop1Call[A]) ErrorFunc(fn ErrorFunc) Loop1Call[A] {
	c.opts.ErrorLogger = errorLoggerFunc(fn)
	return c
}

// ErrorLogger sets the error logger.
func (c Loop1Call[A]) ErrorLogger(logger ErrorLogger) Loop1Call[A] {
	c.opts.ErrorLogger = logger
	return c
}

// Logger sets the default logger.
func (c Loop1Call[A]) Logger(logger logging.Logger) Loop1Call[A] {
	c.opts.Logger = logger
	return c
}

// MinDelay sets the min delay.
func (c Loop1Call[A]) MinDelay(minDelay time.Duration) Loop1Call[A] {
	c.opts.MinDelay = minDelay
	return c
}

// MaxDelay sets the max delay.
func (c Loop1Call[A]) MaxDelay(maxDelay time.Duration) Loop1Call[A] {
	c.opts.MaxDelay = maxDelay
	return c
}

// MaxRetries sets the max retries.
func (c Loop1Call[A]) MaxRetries(maxRetries int) Loop1Call[A] {
	c.opts.MaxRetries = maxRetries
	return c
}

// Options overrides all options.
func (c Loop1Call[A]) Options(opts Options) Loop1Call[A] {
	c.opts = opts
	return c
}

// Run retries the function in a loop.
func (c Loop1Call[A]) Run(ctx async.Context, arg A) status.Status {
	success := new(bool)

	for attempt := 0; ; attempt++ {
		// Restart on success
		if *success {
			attempt = 0
			*success = false
		}

		// Call function
		st := c.run(ctx, arg, success)
		if st.Code == status.CodeCancelled {
			return st
		}

		// Log error
		if !st.OK() {
			c.logError(st, attempt)
		}

		// Sleep before retry
		if st := c.sleep(ctx, attempt); !st.OK() {
			return st
		}
	}
}

// private

func (c Loop1Call[A]) run(ctx async.Context, arg A, success *bool) (st status.Status) {
	defer func() {
		if e := recover(); e != nil {
			st = status.Recover(e)
		}
	}()

	return c.fn(ctx, arg, success)
}
