// Copyright 2026 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package retry

import (
	"github.com/basecomplextech/baselibrary/async"
	"github.com/basecomplextech/baselibrary/status"
)

// VoidRecoverer calls a function once, recovers on panics, logs errors and returns a status.
type VoidRecoverer struct {
	retrier
	fn VoidFunc
}

// RecoverVoid returns a void function recoverer.
//
// Example:
//
//	fn := func(ctx async.Context) status.Status {
//	    // ...
//	}
//
//	st := retry.RecoverVoid(fn).
//		Error("operation failed").
//		ErrorHandler(myErrorHandler).
//		Run(ctx)
func RecoverVoid(fn VoidFunc) VoidRecoverer {
	return VoidRecoverer{
		retrier: newRetrier(),
		fn:      fn,
	}
}

// Run calls the function, recovers on panics, logs errors and returns a status.
func (r VoidRecoverer) Run(ctx async.Context) status.Status {
	// Call function
	st := r.run(ctx)
	switch st.Code {
	case status.CodeOK, status.CodeCancelled:
		return st
	}

	// Handle error
	_ = r.handleError(st, 0)
	return st
}

// Error sets the error message.
func (r VoidRecoverer) Error(message string) VoidRecoverer {
	r.opts.Error = message
	return r
}

// ErrorFunc sets the error handler.
func (r VoidRecoverer) ErrorFunc(fn ErrorFunc) VoidRecoverer {
	r.opts.ErrorHandler = fn
	return r
}

// ErrorHandler sets the error handler.
func (r VoidRecoverer) ErrorHandler(handler ErrorHandler) VoidRecoverer {
	r.opts.ErrorHandler = handler
	return r
}

// private

func (r VoidRecoverer) run(ctx async.Context) (st status.Status) {
	defer func() {
		if e := recover(); e != nil {
			st = status.Recover(e)
		}
	}()

	return r.fn(ctx)
}
