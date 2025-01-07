// Copyright 2025 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package retry

import (
	"time"

	"github.com/basecomplextech/baselibrary/logging"
	"github.com/basecomplextech/baselibrary/status"
)

// Options specifies the options for a retrier.
type Options struct {
	// Error is the error message.
	Error string

	// ErrorLogger is a logger for retry errors.
	ErrorLogger ErrorLogger

	// Logger is the default logger if the error logger is not set.
	Logger logging.Logger

	// MinDelay is the min delay between retries.
	MinDelay time.Duration

	// MaxDelay is the max delay between retries.
	MaxDelay time.Duration

	// MaxRetries is the max retries, zero means unlimited.
	MaxRetries int
}

// Default returns the default options.
func Default() Options {
	return Options{
		Error:    "Failed to execute function",
		Logger:   logging.Stderr,
		MinDelay: 25 * time.Millisecond,
		MaxDelay: 1 * time.Second,
	}
}

// ErrorLogger

// ErrorLogger specifies an interface for logging retry errors.
type ErrorLogger interface {
	// RetryError is called on an error, the attempt is zero-based.
	RetryError(msg string, err status.Status, attempt int)
}

// RetryErrorFunc logs a retry error.
type RetryErrorFunc func(msg string, err status.Status, attempt int)

// private

type errorLoggerFunc func(msg string, err status.Status, attempt int)

func (f errorLoggerFunc) RetryError(msg string, err status.Status, attempt int) {
	f(msg, err, attempt)
}
