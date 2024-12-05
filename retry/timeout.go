// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package retry

import "time"

// Timeout returns an exponential backoff timeout for retrying.
func Timeout(attempt int) time.Duration {
	return timeout(attempt, 0, 0)
}

// TimeoutOpts returns an exponential backoff timeout for retrying.
func TimeoutOpts(attempt int, minTimeout, maxTimeout time.Duration) time.Duration {
	return timeout(attempt, minTimeout, maxTimeout)
}

// private

const (
	defaultMinTimeout = time.Millisecond * 25
	defaultMaxTimeout = time.Second
)

// timeout returns an exponential backoff timeout for retrying.
func timeout(attempt int, minTimeout, maxTimeout time.Duration) time.Duration {
	if attempt == 0 {
		return 0
	}

	if minTimeout == 0 {
		minTimeout = defaultMinTimeout
	}
	if maxTimeout == 0 {
		maxTimeout = defaultMaxTimeout
	}

	multi := uint16(1<<attempt - 1)
	timeout := minTimeout * time.Duration(multi)
	return min(timeout, maxTimeout)
}
