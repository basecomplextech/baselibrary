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

// LoopFunc is a function without arguments that is retried in a loop.
type LoopFunc func(ctx async.Context, success *bool) status.Status

// LoopRetrier retries a loop function.
type LoopRetrier struct {
	retrier
	fn LoopFunc
}

// RetryLoop returns a loop retrier.
//
// Example:
//
//	fn := func(ctx async.Context, success *bool) status.Status {
//	    // ...
//	}
//
//	st := retry.RetryLoop(fn).
//		MaxRetries(5).
//		MinDelay(time.Second).
//		MaxDelay(10 * time.Second).
//		Error("operation failed").
//		ErrorHandler(myErrorHandler).
//		Run(ctx)
func RetryLoop(fn LoopFunc) LoopRetrier {
	return LoopRetrier{
		retrier: newRetrier(),
		fn:      fn,
	}
}

var _ builder[LoopRetrier] = (*LoopRetrier)(nil)

// Run retries the function in a loop.
func (r LoopRetrier) Run(ctx async.Context) status.Status {
	success := new(bool)

	for attempt := 0; ; attempt++ {
		// Restart on success
		if *success {
			attempt = 0
			*success = false
		}

		// Call function
		st := r.run(ctx, success)
		if st.Code == status.CodeCancelled {
			return st
		}

		// Handle error
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
func (r LoopRetrier) Error(message string) LoopRetrier {
	r.opts.Error = message
	return r
}

// ErrorFunc sets the error handler.
func (r LoopRetrier) ErrorFunc(fn ErrorFunc) LoopRetrier {
	r.opts.ErrorHandler = fn
	return r
}

// ErrorHandler sets the error handler.
func (r LoopRetrier) ErrorHandler(handler ErrorHandler) LoopRetrier {
	r.opts.ErrorHandler = handler
	return r
}

// Logger sets the default logger.
func (r LoopRetrier) Logger(logger logging.Logger) LoopRetrier {
	r.opts.Logger = logger
	return r
}

// MinDelay sets the min delay.
func (r LoopRetrier) MinDelay(minDelay time.Duration) LoopRetrier {
	r.opts.MinDelay = minDelay
	return r
}

// MaxDelay sets the max delay.
func (r LoopRetrier) MaxDelay(maxDelay time.Duration) LoopRetrier {
	r.opts.MaxDelay = maxDelay
	return r
}

// MaxRetries sets the max retries.
func (r LoopRetrier) MaxRetries(maxRetries int) LoopRetrier {
	r.opts.MaxRetries = maxRetries
	return r
}

// Options overrides all options.
func (r LoopRetrier) Options(opts Options) LoopRetrier {
	r.opts = opts
	return r
}

// private

func (r LoopRetrier) run(ctx async.Context, success *bool) (st status.Status) {
	defer func() {
		if e := recover(); e != nil {
			st = status.Recover(e)
		}
	}()

	return r.fn(ctx, success)
}
