package status

import (
	"fmt"

	"github.com/complex1tech/baselibrary/panics"
)

// Error returns an error status.
func Error(text string) Status {
	return Status{Code: CodeError, Text: text}
}

// Errorf formats a message and returns an error status.
func Errorf(format string, a ...any) Status {
	text := fmt.Sprintf(format, a...)
	return Status{Code: CodeError, Text: text}
}

// WrapError wraps an error into an error status, or returns ok if the error is nil.
func WrapError(err error) Status {
	if err == nil {
		return OK
	}

	text := err.Error()
	return Status{Code: CodeError, Text: text, Error: err}
}

// WrapErrorf wraps an error, formats a message and returns an error status,
// or returns ok if the error is nil.
func WrapErrorf(err error, format string, a ...any) Status {
	if err == nil {
		return OK
	}

	text := fmt.Sprintf(format, a...)
	text += ": " + err.Error()
	return Status{Code: CodeError, Text: text, Error: err}
}

// Recover recovers from a panic and returns an error status.
func Recover(e any) Status {
	err := panics.Recover(e)
	return WrapError(err)
}

// RecoverStack recovers from a panic and returns an error status and a stack trace.
func RecoverStack(e any) (Status, []byte) {
	err, stack := panics.RecoverStack(e)
	return WrapError(err), stack
}

// IOError

// IOError returns an io error status.
func IOError(text string) Status {
	return Status{Code: CodeIOError, Text: text}
}

// IOErrorf formats a message and returns an io error status.
func IOErrorf(format string, a ...any) Status {
	text := fmt.Sprintf(format, a...)
	return Status{Code: CodeIOError, Text: text}
}

// Corrupted

// Corrupted returns a corrupted status.
func Corrupted(text string) Status {
	return Status{Code: CodeCorrupted, Text: text}
}

// Corruptedf formats a message and returns a corrupted status.
func Corruptedf(format string, a ...any) Status {
	text := fmt.Sprintf(format, a...)
	return Status{Code: CodeCorrupted, Text: text}
}

// Fatal

// Fatal returns a fatal status.
func Fatal(text string) Status {
	return Status{Code: CodeFatal, Text: text}
}

// Fatalf formats a message and returns a fatal status.
func Fatalf(format string, a ...any) Status {
	text := fmt.Sprintf(format, a...)
	return Status{Code: CodeFatal, Text: text}
}

// FatalError wraps an error and returns a fatal status.
func FatalError(err error) Status {
	text := "terminated"
	if err != nil {
		text = err.Error()
	}
	return Status{Code: CodeFatal, Text: text, Error: err}
}

// FatalErrorf wraps an error, formats a message and returns a fatal status.
func FatalErrorf(err error, format string, a ...any) Status {
	text := fmt.Sprintf(format, a...)
	if err != nil {
		text += ": " + err.Error()
	}
	return Status{Code: CodeFatal, Text: text, Error: err}
}
