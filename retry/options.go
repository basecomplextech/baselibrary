// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package retry

import (
	"time"
)

type Options struct {
	// MaxRetries specifies the maximum number of retries, zero means no limit.
	MaxRetries int

	// MinTimeout specifies the minimum timeout for retries, 0 means the default.
	MinTimeout time.Duration

	// MaxTimeout specifies the maximum timeout for retries, 0 means the default.
	MaxTimeout time.Duration

	// ErrorLogger specifies an error logger.
	ErrorLogger ErrorLogger
}
