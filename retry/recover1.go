// Copyright 2026 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package retry

import (
	"github.com/basecomplextech/baselibrary/async"
	"github.com/basecomplextech/baselibrary/status"
)

// VoidRecoverer1 retries a function and returns the result.
type VoidRecoverer1[A any] struct {
	retrier
	fn VoidFunc1[A]
}

// RecoverVoid1 returns a function recoverer.
//
// Example:
//
//	fn := func(ctx async.Context) status.Status {
//	    // ...
//	}
//
//	st := retry.RecoverVoid1(fn).
//		Error("operation failed").
//		ErrorHandler(myErrorHandler).
//		Run(ctx, arg)
func RecoverVoid1[A any](fn VoidFunc1[A]) VoidRecoverer1[A] {
	return VoidRecoverer1[A]{
		retrier: newRetrier(),
		fn:      fn,
	}
}

// Run calls the function, recovers on panics, logs errors and returns a status.
func (r VoidRecoverer1[A]) Run(ctx async.Context, arg A) status.Status {
	// Call function
	st := r.run(ctx, arg)
	switch st.Code {
	case status.CodeOK, status.CodeCancelled:
		return st
	}

	// Handle error
	_ = r.handleError(st, 0)
	return st
}

// Error sets the error message.
func (r VoidRecoverer1[A]) Error(message string) VoidRecoverer1[A] {
	r.opts.Error = message
	return r
}

// ErrorFunc sets the error handler.
func (r VoidRecoverer1[A]) ErrorFunc(fn ErrorFunc) VoidRecoverer1[A] {
	r.opts.ErrorHandler = fn
	return r
}

// ErrorHandler sets the error handler.
func (r VoidRecoverer1[A]) ErrorHandler(handler ErrorHandler) VoidRecoverer1[A] {
	r.opts.ErrorHandler = handler
	return r
}

// private

func (r VoidRecoverer1[A]) run(ctx async.Context, arg A) (st status.Status) {
	defer func() {
		if e := recover(); e != nil {
			st = status.Recover(e)
		}
	}()

	return r.fn(ctx, arg)
}
