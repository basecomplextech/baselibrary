// Copyright 2025 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package retry

import "github.com/basecomplextech/baselibrary/status"

// ErrorFunc logs a retry error, attempt is zero-based.
type ErrorFunc func(msg string, err status.Status, attempt int)

// ErrorLogger specifies an interface for logging retry errors.
type ErrorLogger interface {
	// RetryError is called on an error, attempt is zero-based.
	RetryError(msg string, err status.Status, attempt int)
}

// private

type errorLoggerFunc func(msg string, err status.Status, attempt int)

func (f errorLoggerFunc) RetryError(msg string, err status.Status, attempt int) {
	f(msg, err, attempt)
}
