// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package retry

import (
	"time"

	"github.com/basecomplextech/baselibrary/async"
	"github.com/basecomplextech/baselibrary/logging"
	"github.com/basecomplextech/baselibrary/status"
)

// Retrier retries functions on errors and panics, uses exponential backoff.
type Retrier interface {
	// Recover calls a function once, recovers on panics, logs errors.
	Recover(ctx async.Context, fn func(ctx async.Context) status.Status) status.Status

	// Retry calls a function, retries on errors and panics, uses exponential backoff.
	Retry(ctx async.Context, fn func(ctx async.Context) status.Status) status.Status

	// Loop calls a function in a loop, retries on errors and panics, uses exponential backoff.
	// The method ignores the max retries limit.
	Loop(ctx async.Context, fn func(ctx async.Context, success *bool) status.Status) status.Status
}

// NewRetrier returns a new retrier.
func NewRetrier(opts Options) Retrier {
	return newRetrier(opts)
}

// internal

var _ Retrier = (*retrier)(nil)

type retrier struct {
	opts Options
}

func newRetrier(opts Options) Retrier {
	return &retrier{opts: opts}
}

// Recover calls a function once, recovers on panics, logs errors.
func (r *retrier) Recover(ctx async.Context, fn func(ctx async.Context) status.Status) status.Status {
	// Call function
	st := r.run(ctx, fn)
	switch st.Code {
	case status.CodeOK, status.CodeCancelled:
		return st
	}

	// Log error
	r.logError(0, st)
	return st
}

// Retry calls a function, retries on errors and panics, uses exponential backoff.
func (r *retrier) Retry(ctx async.Context, fn func(ctx async.Context) status.Status) status.Status {
	for attempt := 0; ; attempt++ {
		// Call function
		st := r.run(ctx, fn)
		switch st.Code {
		case status.CodeOK, status.CodeCancelled:
			return st
		}

		// Log error
		r.logError(attempt, st)

		// Check max retries
		if r.opts.MaxRetries != 0 {
			if attempt >= r.opts.MaxRetries {
				return status.Errorf("failed to retry function, max retries reached")
			}
		}

		// Sleep before retry
		delay := timeout(attempt, r.opts.MinTimeout, r.opts.MaxTimeout)
		timer := time.NewTimer(delay)
		select {
		case <-ctx.Wait():
			timer.Stop()
			return ctx.Status()
		case <-timer.C:
		}
	}
}

// Loop calls a function in a loop, retries on errors and panics, uses exponential backoff.
// The method ignores the max retries limit.
func (r *retrier) Loop(ctx async.Context, fn func(ctx async.Context, success *bool) status.Status) status.Status {
	success := new(bool)

	for attempt := 0; ; attempt++ {
		// Restart on success
		if *success {
			attempt = 0
			*success = false
		}

		// Call function
		st := fn(ctx, success)
		if st.Code == status.CodeCancelled {
			return st
		}

		// Log error
		r.logError(attempt, st)

		// Sleep before retry
		delay := timeout(attempt, r.opts.MinTimeout, r.opts.MaxTimeout)
		timer := time.NewTimer(delay)
		select {
		case <-ctx.Wait():
			timer.Stop()
			return ctx.Status()
		case <-timer.C:
		}
	}
}

// private

func (r *retrier) run(ctx async.Context, fn func(ctx async.Context) status.Status) (st status.Status) {
	defer func() {
		if e := recover(); e != nil {
			st = status.Recover(e)
		}
	}()

	return fn(ctx)
}

func (r *retrier) logError(attempt int, st status.Status) {
	if r.opts.ErrorLogger != nil {
		r.opts.ErrorLogger.LogError(attempt, st)
		return
	}

	logger := logging.Stderr
	if r.opts.Logger != nil {
		logger = r.opts.Logger
	}

	if attempt == 0 {
		logger.ErrorStatus("Failed to execute function", st)
	} else {
		logger.DebugStatus("Failed to execute function", st, "attempt", attempt)
	}
}
