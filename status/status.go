package status

import "fmt"

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
func Newf(code Code, format string, a ...any) Status {
	text := fmt.Sprintf(format, a...)
	return Status{Code: code, Text: text}
}

// OK returns true if the status code is OK.
func (s Status) OK() bool {
	return s.Code == CodeOK
}

// Fatal returns true if the status code is terminal.
func (s Status) Fatal() bool {
	return s.Code == CodeFatal
}

// String returns "code: text".
func (s Status) String() string {
	if len(s.Text) == 0 {
		return string(s.Code)
	}
	return fmt.Sprintf("%s: %s", s.Code, s.Text)
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
func (s Status) WrapTextf(format string, a ...any) Status {
	text := fmt.Sprintf(format, a...)
	s1 := s
	s1.Text = text + ": " + s.Text
	return s1
}
