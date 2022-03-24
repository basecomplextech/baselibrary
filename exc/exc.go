package exc

import (
	"fmt"
)

var (
	ErrCancelled = Exc(CodeCancelled, "Cancelled")
	ErrInternal  = Exc(CodeInternal, "Internal error")
)

// Exception is a error with a code and a text message.
type Exception struct {
	Code Code
	Text string
}

// Exc returns a new exception.
func Exc(code Code, format string, a ...any) *Exception {
	var text string
	if len(a) == 0 {
		text = format
	} else {
		text = fmt.Sprintf(format, a...)
	}

	return &Exception{
		Code: code,
		Text: text,
	}
}

// InternalError returns a generic internal error.
func InternalError(format string, a ...any) *Exception {
	return Exc(CodeInternal, format, a...)
}

// IllegalArg returns a new illegal argument error.
func IllegalArg(format string, a ...any) *Exception {
	return Exc(CodeIllegalArg, format, a...)
}

// IllegalOp returns a new illegal operation error.
func IllegalOp(format string, a ...any) *Exception {
	return Exc(CodeIllegalOp, format, a...)
}

// InvalidState returns a new invalid state error.
func InvalidState(format string, a ...any) *Exception {
	return Exc(CodeInvalidState, format, a...)
}

// NotFound returns a new not found error.
func NotFound(format string, a ...any) *Exception {
	return Exc(CodeNotFound, format, a...)
}

// Corrupted returns a new data corruption error.
func Corrupted(format string, a ...any) *Exception {
	return Exc(CodeCorrupted, format, a...)
}

// Unavaible returns a new unavailable error.
func Unavaible(format string, a ...any) *Exception {
	return Exc(CodeUnavailable, format, a...)
}

// GetCode returns an error code or undefined.
func GetCode(err error) Code {
	if err == nil {
		return CodeUndefined
	}

	e, ok := err.(*Exception)
	if ok {
		return e.Code
	}
	return CodeUnavailable
}

// Error returns a "code: message" string.
func (e *Exception) Error() string {
	code := e.Code
	text := code.String()
	return fmt.Sprintf("%s %03d: %s", text, code, e.Text)
}
