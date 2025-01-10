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
	// LoopFunc is a loop function.
	LoopFunc func(ctx async.Context, success *bool) status.Status

	// LoopCall retries a loop function.
	LoopCall struct {
		call
		fn LoopFunc
	}
)

// RetryLoop returns a loop call.
func RetryLoop(fn LoopFunc) LoopCall {
	return LoopCall{
		call: newCall(),
		fn:   fn,
	}
}

var _ builder[LoopCall] = (*LoopCall)(nil)

// Run retries the function in a loop.
func (c LoopCall) Run(ctx async.Context) status.Status {
	success := new(bool)

	for attempt := 0; ; attempt++ {
		// Restart on success
		if *success {
			attempt = 0
			*success = false
		}

		// Call function
		st := c.run(ctx, success)
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

// Error sets the error message.
func (c LoopCall) Error(message string) LoopCall {
	c.opts.Error = message
	return c
}

// ErrorFunc sets the error function.
func (c LoopCall) ErrorFunc(fn ErrorFunc) LoopCall {
	c.opts.ErrorLogger = errorLoggerFunc(fn)
	return c
}

// ErrorLogger sets the error logger.
func (c LoopCall) ErrorLogger(logger ErrorLogger) LoopCall {
	c.opts.ErrorLogger = logger
	return c
}

// Logger sets the default logger.
func (c LoopCall) Logger(logger logging.Logger) LoopCall {
	c.opts.Logger = logger
	return c
}

// MinDelay sets the min delay.
func (c LoopCall) MinDelay(minDelay time.Duration) LoopCall {
	c.opts.MinDelay = minDelay
	return c
}

// MaxDelay sets the max delay.
func (c LoopCall) MaxDelay(maxDelay time.Duration) LoopCall {
	c.opts.MaxDelay = maxDelay
	return c
}

// MaxRetries sets the max retries.
func (c LoopCall) MaxRetries(maxRetries int) LoopCall {
	c.opts.MaxRetries = maxRetries
	return c
}

// Options overrides all options.
func (c LoopCall) Options(opts Options) LoopCall {
	c.opts = opts
	return c
}

// private

func (c LoopCall) run(ctx async.Context, success *bool) (st status.Status) {
	defer func() {
		if e := recover(); e != nil {
			st = status.Recover(e)
		}
	}()

	return c.fn(ctx, success)
}
