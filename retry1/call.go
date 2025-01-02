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

// builder provides chained methods for building a call.
type builder[C any] interface {
	// Error sets the error message.
	Error(message string) C

	// ErrorLogger sets the error logger.
	ErrorLogger(logger ErrorLogger) C

	// Logger sets the default logger.
	Logger(logger logging.Logger) C

	// MaxDelay sets the max delay.
	MaxDelay(maxDelay time.Duration) C

	// MaxRetries sets the max retries.
	MaxRetries(maxRetries int) C

	// Options overrides all options.
	Options(opts Options) C
}

// call

type call struct {
	opts Options
}

func newCall() call {
	return call{opts: Default()}
}

func (c call) logError(err status.Status, attempt int) {
	msg := c.opts.Error

	// Use error logger if set
	if logger := c.opts.ErrorLogger; logger != nil {
		logger.RetryError(msg, err, attempt)
		return
	}

	// Skip logging if no logger
	logger := c.opts.Logger
	if logger == nil {
		return
	}

	// Log error
	if attempt == 0 || attempt%10 == 0 {
		logger.ErrorStatus(msg, err)
	} else {
		logger.DebugStatus(msg, err)
	}
}

func (c call) sleep(ctx async.Context, attempt int) status.Status {
	// Check max retries
	if c.opts.MaxRetries != 0 {
		if attempt >= c.opts.MaxRetries {
			return status.Errorf("failed to retry function, max retries reached")
		}
	}

	// Sleep before retry
	delay := delay(attempt, c.opts.MinDelay, c.opts.MaxDelay)
	timer := time.NewTimer(delay)
	select {
	case <-ctx.Wait():
		timer.Stop()
		return ctx.Status()
	case <-timer.C:
		return status.OK
	}
}
