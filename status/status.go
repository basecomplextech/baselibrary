package status

import "fmt"

// Status represents an operation status.
type Status struct {
	Code    Code
	Message string
	Error   error
}

// New returns a new status.
func New(code Code, msg string) Status {
	return Status{Code: code, Message: msg}
}

// Newf returns a new status and formats its message.
func Newf(code Code, format string, a ...any) Status {
	msg := fmt.Sprintf(format, a...)
	return Status{Code: code, Message: msg}
}

// OK returns true if the status code is OK.
func (s Status) OK() bool {
	return s.Code == CodeOK
}

// String returns "code: msg".
func (s Status) String() string {
	if len(s.Message) == 0 {
		return string(s.Code)
	}
	return fmt.Sprintf("%s: %s", s.Code, s.Message)
}

// With

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

// WrapText returns a status clone with a new msg and an appended previous msg.
func (s Status) WrapText(msg string) Status {
	s1 := s
	s1.Message = msg + ": " + s.Message
	return s1
}

// WrapTextf returns a status clone with a new msg and an appended previous msg.
func (s Status) WrapTextf(format string, a ...any) Status {
	msg := fmt.Sprintf(format, a...)
	s1 := s
	s1.Message = msg + ": " + s.Message
	return s1
}
