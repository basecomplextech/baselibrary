// Copyright 2025 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package retry

import "time"

const (
	MinDelay       = 25 * time.Millisecond
	MinDelayMedium = 250 * time.Millisecond
	MaxDelay       = 1 * time.Second
)

// Delay returns a delay for retrying, uses an exponential backoff.
func Delay(attempt int) time.Duration {
	return delay(attempt, 0, 0)
}

// DelayOpts returns a delay for retrying, uses an exponential backoff.
func DelayOpts(attempt int, minDelay, maxDelay time.Duration) time.Duration {
	return delay(attempt, minDelay, maxDelay)
}

// private

// delay returns a delay for retrying, uses an exponential backoff.
func delay(attempt int, minDelay, maxDelay time.Duration) time.Duration {
	if attempt == 0 {
		return 0
	}

	if minDelay == 0 {
		minDelay = MinDelay
	}
	if maxDelay == 0 {
		maxDelay = MaxDelay
	}

	multi := uint16(1<<attempt - 1)
	delay := minDelay * time.Duration(multi)
	return min(delay, maxDelay)
}
