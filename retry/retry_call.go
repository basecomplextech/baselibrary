// Copyright 2025 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package retry

import (
	"time"

	"github.com/basecomplextech/baselibrary/async"
	"github.com/basecomplextech/baselibrary/status"
)

type call struct {
	opts Options
}

func newCall() call {
	return call{opts: Default()}
}

func (c call) handleError(err status.Status, attempt int) status.Status {
	msg := c.opts.Error

	// Maybe use error handler
	if handler := c.opts.ErrorHandler; handler != nil {
		return handler.RetryError(msg, err, attempt)
	}

	// Skip logging if no logger
	logger := c.opts.Logger
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

func (c call) sleep(ctx async.Context, attempt int) status.Status {
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
