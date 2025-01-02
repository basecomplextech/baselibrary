// Copyright 2025 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package retry

import "time"

// Delay returns a delay for retrying, uses an exponential backoff.
func Delay(attempt int) time.Duration {
	return delay(attempt, 0, 0)
}

// DelayOpts returns a delay for retrying, uses an exponential backoff.
func DelayOpts(attempt int, minDelay, maxDelay time.Duration) time.Duration {
	return delay(attempt, minDelay, maxDelay)
}

// private

const (
	defaultMinDelay = time.Millisecond * 25
	defaultMaxDelay = time.Second
)

// delay returns a delay for retrying, uses an exponential backoff.
func delay(attempt int, minDelay, maxDelay time.Duration) time.Duration {
	if attempt == 0 {
		return 0
	}

	if minDelay == 0 {
		minDelay = defaultMinDelay
	}
	if maxDelay == 0 {
		maxDelay = defaultMaxDelay
	}

	multi := uint16(1<<attempt - 1)
	delay := minDelay * time.Duration(multi)
	return min(delay, maxDelay)
}
