package status

import "fmt"

// Cancelled

// Cancelledf formats a message and returns a cancelled status.
func Cancelledf(format string, a ...interface{}) Status {
	text := fmt.Sprintf(format, a...)
	return Status{Code: CodeCancelled, Text: text}
}

// Timeout

// Timeoutf returns a timeout status and formats its message.
func Timeoutf(format string, a ...interface{}) Status {
	text := fmt.Sprintf(format, a...)
	return Status{Code: CodeTimeout, Text: text}
}

// Unavailable

// Unavailable returns an unavailable status.
func Unavailable(text string) Status {
	return Status{Code: CodeUnavailable, Text: text}
}

// Unavailablef returns an unavailable status and formats its message.
func Unavailablef(format string, a ...interface{}) Status {
	text := fmt.Sprintf(format, a...)
	return Status{Code: CodeUnavailable, Text: text}
}
