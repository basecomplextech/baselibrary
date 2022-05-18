package status

import "fmt"

var (
	OK      = New(CodeOK, "ok")
	Stopped = New(CodeStopped, "stopped")
	Timeout = New(CodeTimeout, "timeout")
)

type Status struct {
	Code  Code
	Text  string
	Error error
}

// New returns a new status.
func New(code Code, text string) Status {
	return Status{Code: code, Text: text}
}

// Newf returns a new status and formats its message.
func Newf(code Code, format string, a ...interface{}) Status {
	text := fmt.Sprintf(format, a...)
	return Status{Code: code, Text: text}
}

// OK returns true if the status code is OK.
func (s Status) OK() bool {
	return s.Code == CodeOK
}

// IsError returns true if the status code is error.
func (s Status) IsError() bool {
	return s.Code == CodeError
}

// String returns "code: text".
func (s Status) String() string {
	return fmt.Sprintf("%s: %s", s.Code, s.Text)
}

// Wrap wraps an existing status into a new message.
func (s Status) Wrap(text string) Status {
	text1 := text + ": " + s.Text
	return Status{Code: s.Code, Text: text1, Error: s.Error}
}

// Wrapf wraps an existing status into a new message.
func (s Status) Wrapf(format string, a ...interface{}) Status {
	text := fmt.Sprintf(format, a...)
	return s.Wrap(text)
}

// Utility constructors

// Error returns an error status.
func Error(text string) Status {
	return Status{Code: CodeError, Text: text}
}

// Errorf returns an error status and formats its message.
func Errorf(format string, a ...interface{}) Status {
	text := fmt.Sprintf(format, a...)
	return Status{Code: CodeError, Text: text}
}

// NotFound returns a not found status.
func NotFound(text string) Status {
	return Status{Code: CodeNotFound, Text: text}
}

// NotFoundf returns a not found status and formats its message.
func NotFoundf(format string, a ...interface{}) Status {
	text := fmt.Sprintf(format, a...)
	return Status{Code: CodeNotFound, Text: text}
}

// Stoppedf returns a stopped status and formats its message.
func Stoppedf(format string, a ...interface{}) Status {
	text := fmt.Sprintf(format, a...)
	return Status{Code: CodeStopped, Text: text}
}

// Timeoutf returns a timeout status and formats its message.
func Timeoutf(format string, a ...interface{}) Status {
	text := fmt.Sprintf(format, a...)
	return Status{Code: CodeTimeout, Text: text}
}

// Unavailable returns an unavailable status.
func Unavailable(text string) Status {
	return Status{Code: CodeUnavailable, Text: text}
}

// Unavailablef returns an unavailable status and formats its message.
func Unavailablef(format string, a ...interface{}) Status {
	text := fmt.Sprintf(format, a...)
	return Status{Code: CodeUnavailable, Text: text}
}

// WrapError wraps an error into a status or returns ok if err is nil.
func WrapError(err error) Status {
	if err == nil {
		return OK
	}

	text := err.Error()
	return Status{Code: CodeError, Text: text, Error: err}
}

// WrapErrorf wraps an error into a status or returns ok if err is nil.
func WrapErrorf(err error, format string, a ...interface{}) Status {
	if err == nil {
		return OK
	}

	text := fmt.Sprintf(format, a...)
	return Status{Code: CodeError, Text: text, Error: err}
}
