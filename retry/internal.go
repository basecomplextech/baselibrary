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

// builder provides chained methods for building a call.
type builder[C any] interface {
	// Error sets the error message.
	Error(message string) C

	// ErrorFunc sets the error handler.
	ErrorFunc(fn ErrorFunc) C

	// ErrorHandler sets the error handler.
	ErrorHandler(handler ErrorHandler) C

	// Logger sets the default logger.
	Logger(logger logging.Logger) C

	// MinDelay sets the min delay.
	MinDelay(minDelay time.Duration) C

	// MaxDelay sets the max delay.
	MaxDelay(maxDelay time.Duration) C

	// MaxRetries sets the max retries.
	MaxRetries(maxRetries int) C

	// Options overrides all options.
	Options(opts Options) C
}

// retrier

// retrier provides common retry methods.
type retrier struct {
	opts Options
}

func newRetrier() retrier {
	return retrier{opts: Default()}
}

func (r retrier) handleError(err status.Status, attempt int) status.Status {
	msg := r.opts.Error

	// Maybe use error handler
	if handler := r.opts.ErrorHandler; handler != nil {
		return handler.RetryError(msg, err, attempt)
	}

	// Skip logging if no logger
	logger := r.opts.Logger
	if logger == nil {
		return status.OK
	}

	// Log error
	if attempt == 0 || attempt%10 == 0 {
		logger.ErrorStatus(msg, err)
	} else {
		logger.DebugStatus(msg, err)
	}
	return status.OK
}

func (r retrier) sleep(ctx async.Context, attempt int) status.Status {
	// Sleep before retry
	delay := delay(attempt, r.opts.MinDelay, r.opts.MaxDelay)
	timer := time.NewTimer(delay)
	select {
	case <-ctx.Wait():
		timer.Stop()
		return ctx.Status()
	case <-timer.C:
		return status.OK
	}
}
