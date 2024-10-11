// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package status

import "fmt"

// ConcurrencyError

// ConcurrencyError returns a concurrency error status.
func ConcurrencyError(msg string) Status {
	return Status{Code: CodeConcurrencyError, Message: msg}
}

// ConcurrencyErrorf formats a message and returns a concurrency error status.
func ConcurrencyErrorf(format string, a ...any) Status {
	msg := fmt.Sprintf(format, a...)
	return Status{Code: CodeConcurrencyError, Message: msg}
}

// Rollback

// Rollback returns a rollback status.
func Rollback(msg string) Status {
	return Status{Code: CodeRollback, Message: msg}
}

// Rollbackf formats a message and returns a rollback status.
func Rollbackf(format string, a ...any) Status {
	msg := fmt.Sprintf(format, a...)
	return Status{Code: CodeRollback, Message: msg}
}
