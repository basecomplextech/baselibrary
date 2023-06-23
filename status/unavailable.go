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
	text := fmt.Sprintf(format, a...)
	return Status{Code: CodeClosed, Text: text}
}

// Cancelled

// Cancelledf formats a message and returns a cancelled status.
func Cancelledf(format string, a ...any) Status {
	text := fmt.Sprintf(format, a...)
	return Status{Code: CodeCancelled, Text: text}
}

// Rollback

// Rollback returns a rollback status.
func Rollback(text string) Status {
	return Status{Code: CodeRollback, Text: text}
}

// Rollbackf formats a message and returns a rollback status.
func Rollbackf(format string, a ...any) Status {
	text := fmt.Sprintf(format, a...)
	return Status{Code: CodeRollback, Text: text}
}

// Timeout

// Timeoutf returns a timeout status and formats its message.
func Timeoutf(format string, a ...any) Status {
	text := fmt.Sprintf(format, a...)
	return Status{Code: CodeTimeout, Text: text}
}

// Unavailable

// Unavailable returns an unavailable status.
func Unavailable(text string) Status {
	return Status{Code: CodeUnavailable, Text: text}
}

// Unavailablef returns an unavailable status and formats its message.
func Unavailablef(format string, a ...any) Status {
	text := fmt.Sprintf(format, a...)
	return Status{Code: CodeUnavailable, Text: text}
}

// Unsupported

// Unsupported returns an unsupported status.
func Unsupported(text string) Status {
	return Status{Code: CodeUnsupported, Text: text}
}

// Unsupportedf returns an unsupported status and formats its message.
func Unsupportedf(format string, a ...any) Status {
	text := fmt.Sprintf(format, a...)
	return Status{Code: CodeUnsupported, Text: text}
}
