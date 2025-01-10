// Copyright 2025 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package retry

import (
	"time"

	"github.com/basecomplextech/baselibrary/logging"
)

// builder provides chained methods for building a call.
type builder[C any] interface {
	// Error sets the error message.
	Error(message string) C

	// ErrorFunc sets the error function.
	ErrorFunc(fn ErrorFunc) C

	// ErrorLogger sets the error logger.
	ErrorLogger(logger ErrorLogger) C

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
