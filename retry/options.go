// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package retry

import (
	"time"

	"github.com/basecomplextech/baselibrary/logging"
	"github.com/basecomplextech/baselibrary/status"
)

type Options struct {
	// MaxRetries specifies the maximum number of retries, zero means no limit.
	MaxRetries int

	// MinTimeout specifies the minimum timeout for retries, 0 means the default.
	MinTimeout time.Duration

	// MaxTimeout specifies the maximum timeout for retries, 0 means the default.
	MaxTimeout time.Duration

	// Logger specifies a logger to log errors if ErrorLogger is not set.
	Logger logging.Logger

	// ErrorLogger specifies an error logger.
	ErrorLogger ErrorLogger
}

type ErrorLogger interface {
	// LogError logs an error.
	LogError(attempt int, st status.Status)
}

// ErrorLoggerFunc returns an error logger from a function.
func ErrorLoggerFunc(fn func(attempt int, st status.Status)) ErrorLogger {
	return errorLoggerFunc(fn)
}

// private

type errorLoggerFunc func(attempt int, st status.Status)

func (f errorLoggerFunc) LogError(attempt int, st status.Status) {
	f(attempt, st)
}
