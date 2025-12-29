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

// VoidFunc is a function without arguments and results.
type VoidFunc func(ctx async.Context) status.Status

// VoidRetrier retries a void function.
type VoidRetrier struct {
	retrier
	fn VoidFunc
}

// RetryVoid returns a void function retrier.
//
// Example:
//
//	fn := func(ctx async.Context) status.Status {
//	    // ...
//	}
//
//	st := retry.RetryVoid(fn).
//		MaxRetries(5).
//		MinDelay(time.Second).
//		MaxDelay(10 * time.Second).
//		Error("operation failed").
//		ErrorHandler(myErrorHandler).
//		Run(ctx)
func RetryVoid(fn VoidFunc) VoidRetrier {
	return VoidRetrier{
		retrier: newRetrier(),
		fn:      fn,
	}
}

var _ builder[VoidRetrier] = (*VoidRetrier)(nil)

// Run retries the procedure.
func (r VoidRetrier) Run(ctx async.Context) status.Status {
	for attempt := 0; ; attempt++ {
		// VoidRetrier function
		st := r.run(ctx)
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
func (r VoidRetrier) Error(message string) VoidRetrier {
	r.opts.Error = message
	return r
}

// ErrorFunc sets the error handler.
func (r VoidRetrier) ErrorFunc(fn ErrorFunc) VoidRetrier {
	r.opts.ErrorHandler = fn
	return r
}

// ErrorHandler sets the error handler.
func (r VoidRetrier) ErrorHandler(handler ErrorHandler) VoidRetrier {
	r.opts.ErrorHandler = handler
	return r
}

// Logger sets the default logger.
func (r VoidRetrier) Logger(logger logging.Logger) VoidRetrier {
	r.opts.Logger = logger
	return r
}

// MinDelay sets the min delay.
func (r VoidRetrier) MinDelay(minDelay time.Duration) VoidRetrier {
	r.opts.MinDelay = minDelay
	return r
}

// MaxDelay sets the max delay.
func (r VoidRetrier) MaxDelay(maxDelay time.Duration) VoidRetrier {
	r.opts.MaxDelay = maxDelay
	return r
}

// MaxRetries sets the max retries.
func (r VoidRetrier) MaxRetries(maxRetries int) VoidRetrier {
	r.opts.MaxRetries = maxRetries
	return r
}

// Options overrides all options.
func (r VoidRetrier) Options(opts Options) VoidRetrier {
	r.opts = opts
	return r
}

// private

func (r VoidRetrier) run(ctx async.Context) (st status.Status) {
	defer func() {
		if e := recover(); e != nil {
			st = status.Recover(e)
		}
	}()

	return r.fn(ctx)
}
