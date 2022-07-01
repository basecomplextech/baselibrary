package status

import "fmt"

var (
	OK        = New(CodeOK, "")
	End       = New(CodeEnd, "")
	Wait      = New(CodeWait, "")
	Cancelled = New(CodeCancelled, "")
	Timeout   = New(CodeTimeout, "")
)

// Status represents an operation status.
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

// String returns "code: text".
func (s Status) String() string {
	if len(s.Text) == 0 {
		return string(s.Code)
	}
	return fmt.Sprintf("%s: %s", s.Code, s.Text)
}

// Is methods

// IsError returns true if the status code is error.
func (s Status) IsError() bool {
	return s.Code == CodeError
}

// IsCancelled returns true if the status code is cancelled.
func (s Status) IsCancelled() bool {
	return s.Code == CodeCancelled
}

// IsTerminal returns true if the status code is terminal.
func (s Status) IsTerminal() bool {
	return s.Code == CodeTerminal
}

// With methods

// WithCode returns a status clone with a new code.
func (s Status) WithCode(code Code) Status {
	s1 := s
	s1.Code = code
	return s1
}

// WithError returns a status clone with a new error.
func (s Status) WithError(err error) Status {
	s1 := s
	s1.Error = err
	return s1
}

// WrapText returns a status clone with a new text and an appended previous text.
func (s Status) WrapText(text string) Status {
	s1 := s
	s1.Text = text + ": " + s.Text
	return s1
}

// WrapTextf returns a status clone with a new text and an appended previous text.
func (s Status) WrapTextf(format string, a ...interface{}) Status {
	text := fmt.Sprintf(format, a...)
	s1 := s
	s1.Text = text + ": " + s.Text
	return s1
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

// Corrupted returns a corrupted status.
func Corrupted(text string) Status {
	return Status{Code: CodeCorrupted, Text: text}
}

// Corruptedf returns a corrupted status and formats its message.
func Corruptedf(format string, a ...interface{}) Status {
	text := fmt.Sprintf(format, a...)
	return Status{Code: CodeCorrupted, Text: text}
}

// Stoppedf returns a stopped status and formats its message.
func Stoppedf(format string, a ...interface{}) Status {
	text := fmt.Sprintf(format, a...)
	return Status{Code: CodeCancelled, Text: text}
}

// Timeoutf returns a timeout status and formats its message.
func Timeoutf(format string, a ...interface{}) Status {
	text := fmt.Sprintf(format, a...)
	return Status{Code: CodeTimeout, Text: text}
}

// Terminal returns a terminal status.
func Terminal(text string) Status {
	return Status{Code: CodeTerminal, Text: text}
}

// Terminalf returns a terminal status and formats its message.
func Terminalf(format string, a ...interface{}) Status {
	text := fmt.Sprintf(format, a...)
	return Status{Code: CodeTerminal, Text: text}
}

// TerminalError wraps an error and returns a terminal status.
func TerminalError(err error) Status {
	text := "terminated"
	if err != nil {
		text = err.Error()
	}
	return Status{Code: CodeTerminal, Text: text, Error: err}
}

// TerminalErrorf wraps an error, formats a message and returns a terminal status.
func TerminalErrorf(err error, format string, a ...interface{}) Status {
	text := fmt.Sprintf(format, a...)
	if err != nil {
		text += ": " + err.Error()
	}
	return Status{Code: CodeTerminal, Text: text, Error: err}
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
