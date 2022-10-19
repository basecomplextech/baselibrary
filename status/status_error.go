package status

import (
	"fmt"

	"github.com/complex1tech/baselibrary/errors2"
)

// Error returns an error status.
func Error(text string) Status {
	return Status{Code: CodeError, Text: text}
}

// Errorf formats a message and returns an error status.
func Errorf(format string, a ...interface{}) Status {
	text := fmt.Sprintf(format, a...)
	return Status{Code: CodeError, Text: text}
}

// Corruption

// Corruption returns a corrupted status.
func Corruption(text string) Status {
	return Status{Code: CodeCorruption, Text: text}
}

// Corruptionf formats a message and returns a corrupted status.
func Corruptionf(format string, a ...interface{}) Status {
	text := fmt.Sprintf(format, a...)
	return Status{Code: CodeCorruption, Text: text}
}

// IOError

// IOError returns an io error status.
func IOError(text string) Status {
	return Status{Code: CodeIOError, Text: text}
}

// IOErrorf formats a message and returns an io error status.
func IOErrorf(format string, a ...interface{}) Status {
	text := fmt.Sprintf(format, a...)
	return Status{Code: CodeIOError, Text: text}
}

// NotFound

// NotFound returns a not found status.
func NotFound(text string) Status {
	return Status{Code: CodeNotFound, Text: text}
}

// NotFoundf formats a message and returns a not found status.
func NotFoundf(format string, a ...interface{}) Status {
	text := fmt.Sprintf(format, a...)
	return Status{Code: CodeNotFound, Text: text}
}

// WrapError

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
func WrapErrorf(err error, format string, a ...interface{}) Status {
	if err == nil {
		return OK
	}

	text := fmt.Sprintf(format, a...)
	text += ": " + err.Error()
	return Status{Code: CodeError, Text: text, Error: err}
}

// Recover

// Recover recovers from a panic and returns an error status.
func Recover(e interface{}) Status {
	err := errors2.Recover(e)
	return WrapError(err)
}
