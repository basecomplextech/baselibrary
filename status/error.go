package status

import (
	"fmt"

	"github.com/basecomplextech/baselibrary/panics"
)

// Error

// Error returns an internal error status.
func Error(msg string) Status {
	return Status{
		Code:    CodeError,
		Message: msg,
	}
}

// Errorf formats and returns an internal error status.
func Errorf(format string, a ...any) Status {
	msg := format
	if len(a) > 0 {
		msg = fmt.Sprintf(format, a...)
	}

	return Status{
		Code:    CodeError,
		Message: msg,
	}
}

// WrapError returns an internal error status.
func WrapError(err error) Status {
	msg := "Internal error"
	if err != nil {
		msg = err.Error()
	}

	return Status{
		Code:    CodeError,
		Message: msg,
		Error:   err,
	}
}

// WrapErrorf formats and returns an internal error status.
func WrapErrorf(err error, format string, a ...any) Status {
	msg := format
	if len(a) > 0 {
		msg = fmt.Sprintf(format, a...)
	}
	if err != nil {
		msg += ": " + err.Error()
	}

	return Status{
		Code:    CodeError,
		Message: msg,
		Error:   err,
	}
}

// ExternalError

// ExternalError returns an external error status.
func ExternalError(msg string) Status {
	return Status{
		Code:    CodeExternalError,
		Message: msg,
	}
}

// ExternalErrorf formats and returns an external error status.
func ExternalErrorf(format string, a ...any) Status {
	msg := format
	if len(a) > 0 {
		msg = fmt.Sprintf(format, a...)
	}

	return Status{
		Code:    CodeExternalError,
		Message: msg,
	}
}

// WrapExternalError returns an external error status.
func WrapExternalError(err error) Status {
	msg := "External error"
	if err != nil {
		msg = err.Error()
	}

	return Status{
		Code:    CodeExternalError,
		Message: msg,
		Error:   err,
	}
}

// WrapExternalErrorf formats and returns an external error status.
func WrapExternalErrorf(err error, format string, a ...any) Status {
	msg := format
	if len(a) > 0 {
		msg = fmt.Sprintf(format, a...)
	}
	if err != nil {
		msg += ": " + err.Error()
	}

	return Status{
		Code:    CodeExternalError,
		Message: msg,
		Error:   err,
	}
}

// Corrupted

// Corrupted returns a data corruption error status.
func Corrupted(msg string) Status {
	return Status{
		Code:    CodeCorrupted,
		Message: msg,
	}
}

// Corruptedf formats and returns a data corruption error status.
func Corruptedf(format string, a ...any) Status {
	msg := format
	if len(a) > 0 {
		msg = fmt.Sprintf(format, a...)
	}

	return Status{
		Code:    CodeCorrupted,
		Message: msg,
	}
}

// WrapCorrupted returns a data corruption error status.
func WrapCorrupted(err error) Status {
	msg := "Data corrupted"
	if err != nil {
		msg = err.Error()
	}

	return Status{
		Code:    CodeCorrupted,
		Message: msg,
		Error:   err,
	}
}

// WrapCorruptedf formats and returns a data corruption error status.
func WrapCorruptedf(err error, format string, a ...any) Status {
	msg := format
	if len(a) > 0 {
		msg = fmt.Sprintf(format, a...)
	}
	if err != nil {
		msg += ": " + err.Error()
	}

	return Status{
		Code:    CodeCorrupted,
		Message: msg,
		Error:   err,
	}
}

// Fatal

// Fatal returns a fatal error status.
func Fatal(msg string) Status {
	return Status{
		Code:    CodeFatal,
		Message: msg,
	}
}

// Fatalf formats and returns a fatal error status.
func Fatalf(format string, a ...any) Status {
	msg := format
	if len(a) > 0 {
		msg = fmt.Sprintf(format, a...)
	}

	return Status{
		Code:    CodeFatal,
		Message: msg,
	}
}

// WrapFatal returns a fatal error status.
func WrapFatal(err error) Status {
	msg := "Fatal error"
	if err != nil {
		msg = err.Error()
	}

	return Status{
		Code:    CodeFatal,
		Message: msg,
		Error:   err,
	}
}

// WrapFatalf formats and returns a fatal error status.
func WrapFatalf(err error, format string, a ...any) Status {
	msg := format
	if len(a) > 0 {
		msg = fmt.Sprintf(format, a...)
	}
	if err != nil {
		msg += ": " + err.Error()
	}

	return Status{
		Code:    CodeFatal,
		Message: msg,
		Error:   err,
	}
}

// Recover

// Recover recovers from a panic and returns an internal error status.
func Recover(e any) Status {
	err := panics.Recover(e)
	return WrapError(err)
}

// RecoverStack recovers from a panic and returns an internal error status and a stack trace.
func RecoverStack(e any) (Status, []byte) {
	err, stack := panics.RecoverStack(e)
	return WrapError(err), stack
}
