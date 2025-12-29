// Copyright 2025 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package retry

import "github.com/basecomplextech/baselibrary/status"

// ErrorHandler handles retry errors.
type ErrorHandler interface {
	// RetryError handles an error, attempt is zero-based.
	RetryError(msg string, err status.Status, attempt int) status.Status
}

// Func

// ErrorFunc handles retry errors, attempt is zero-based.
type ErrorFunc func(msg string, err status.Status, attempt int) status.Status

// RetryError handles an error, attempt is zero-based.
func (f ErrorFunc) RetryError(msg string, err status.Status, attempt int) status.Status {
	return f(msg, err, attempt)
}

// Sample

// Sample10 returns true for every tenth attempt.
func Sample10(attempt int) bool {
	return attempt == 0 || attempt%10 == 0
}
