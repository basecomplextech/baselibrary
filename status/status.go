// Copyright 2022 Ivan Korobkov. All rights reserved.

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

// Cancelled returns true if the status code is Cancelled.
func (s Status) Cancelled() bool {
	return s.Code == CodeCancelled
}

// String returns "code: msg".
func (s Status) String() string {
	code := string(s.Code)
	if s.Code == CodeNone {
		code = "none"
	}
	if len(s.Message) == 0 {
		return code
	}
	return fmt.Sprintf("%s: %s", code, s.Message)
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
	if err == nil {
		return s
	}

	s1 := s
	s1.Error = err
	s1 = s1.WrapText(err.Error())
	return s1
}

// WrapText returns a status clone with a new msg and an appended previous msg.
func (s Status) WrapText(msg string) Status {
	s1 := s
	if len(s.Message) == 0 {
		s1.Message = msg
	} else {
		s1.Message = msg + ": " + s.Message
	}
	return s1
}

// WrapTextf returns a status clone with a new msg and an appended previous msg.
func (s Status) WrapTextf(format string, a ...any) Status {
	s1 := s
	msg := fmt.Sprintf(format, a...)
	if len(s.Message) == 0 {
		s1.Message = msg
	} else {
		s1.Message = msg + ": " + s.Message
	}
	return s1
}
