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
	// VoidFunc is a procedure without arguments and results.
	VoidFunc func(ctx async.Context) status.Status

	// VoidCall retries a procedure, i.e. a function without results.
	VoidCall struct {
		call
		fn VoidFunc
	}
)

// RetryVoid returns a new procedure call.
func RetryVoid(fn VoidFunc) VoidCall {
	return VoidCall{fn: fn}
}

var _ builder[VoidCall] = (*VoidCall)(nil)

// Error sets the error message.
func (c VoidCall) Error(message string) VoidCall {
	c.opts.Error = message
	return c
}

// ErrorFunc sets the error function.
func (c VoidCall) ErrorFunc(fn ErrorFunc) VoidCall {
	c.opts.ErrorLogger = errorLoggerFunc(fn)
	return c
}

// ErrorLogger sets the error logger.
func (c VoidCall) ErrorLogger(logger ErrorLogger) VoidCall {
	c.opts.ErrorLogger = logger
	return c
}

// Logger sets the default logger.
func (c VoidCall) Logger(logger logging.Logger) VoidCall {
	c.opts.Logger = logger
	return c
}

// MinDelay sets the min delay.
func (c VoidCall) MinDelay(minDelay time.Duration) VoidCall {
	c.opts.MinDelay = minDelay
	return c
}

// MaxDelay sets the max delay.
func (c VoidCall) MaxDelay(maxDelay time.Duration) VoidCall {
	c.opts.MaxDelay = maxDelay
	return c
}

// MaxRetries sets the max retries.
func (c VoidCall) MaxRetries(maxRetries int) VoidCall {
	c.opts.MaxRetries = maxRetries
	return c
}

// Options overrides all options.
func (c VoidCall) Options(opts Options) VoidCall {
	c.opts = opts
	return c
}

// Run retries the procedure.
func (c VoidCall) Run(ctx async.Context) status.Status {
	for attempt := 0; ; attempt++ {
		// VoidCall function
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

// private

func (c VoidCall) run(ctx async.Context) (st status.Status) {
	defer func() {
		if e := recover(); e != nil {
			st = status.Recover(e)
		}
	}()

	return c.fn(ctx)
}
