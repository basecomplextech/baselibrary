// Copyright 2022 Ivan Korobkov. All rights reserved.

package status

import "fmt"

var (
	Closed    = New(CodeClosed, "")
	Cancelled = New(CodeCancelled, "")
	Timeout   = New(CodeTimeout, "")
)

// Closed

// Closedf formats a message and returns a closed status.
func Closedf(format string, a ...any) Status {
	msg := fmt.Sprintf(format, a...)
	return Status{Code: CodeClosed, Message: msg}
}

// Cancelled

// Cancelledf formats a message and returns a cancelled status.
func Cancelledf(format string, a ...any) Status {
	msg := fmt.Sprintf(format, a...)
	return Status{Code: CodeCancelled, Message: msg}
}

// Redirect

// Redirect returns a redirect status.
func Redirect(msg string) Status {
	return Status{Code: CodeRedirect, Message: msg}
}

// Redirectf formats a message and returns a redirect status.
func Redirectf(format string, a ...any) Status {
	msg := fmt.Sprintf(format, a...)
	return Status{Code: CodeRedirect, Message: msg}
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

// Timeout

// Timeoutf returns a timeout status and formats its message.
func Timeoutf(format string, a ...any) Status {
	msg := fmt.Sprintf(format, a...)
	return Status{Code: CodeTimeout, Message: msg}
}

// Unavailable

// Unavailable returns an unavailable status.
func Unavailable(msg string) Status {
	return Status{Code: CodeUnavailable, Message: msg}
}

// Unavailablef returns an unavailable status and formats its message.
func Unavailablef(format string, a ...any) Status {
	msg := fmt.Sprintf(format, a...)
	return Status{Code: CodeUnavailable, Message: msg}
}

// Unsupported

// Unsupported returns an unsupported status.
func Unsupported(msg string) Status {
	return Status{Code: CodeUnsupported, Message: msg}
}

// Unsupportedf returns an unsupported status and formats its message.
func Unsupportedf(format string, a ...any) Status {
	msg := fmt.Sprintf(format, a...)
	return Status{Code: CodeUnsupported, Message: msg}
}
