// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package retry

import (
	"github.com/basecomplextech/baselibrary/logging"
	"github.com/basecomplextech/baselibrary/status"
)

type ErrorLogger interface {
	// RetryError is called on an error, the attempt is zero-based.
	RetryError(st status.Status, attempt int)
}

// ErrorLoggerFunc returns an error logger from a function.
func ErrorLoggerFunc(fn func(attempt int, st status.Status)) ErrorLogger {
	return errorLoggerFunc(fn)
}

// private

type errorLoggerFunc func(attempt int, st status.Status)

func (f errorLoggerFunc) RetryError(st status.Status, attempt int) {
	f(attempt, st)
}

var defaultErrorLogger = ErrorLoggerFunc(func(attempt int, st status.Status) {
	if attempt == 0 {
		logging.Stderr.ErrorStatus("Failed to execute function", st)
	} else {
		logging.Stderr.DebugStatus("Failed to execute function", st, "attempt", attempt)
	}
})
