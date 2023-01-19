package status

import "fmt"

var (
	OK   = New(CodeOK, "")
	None = New(CodeNone, "")
)

// OKf formats a message and returns an ok status.
func OKf(format string, a ...any) Status {
	text := fmt.Sprintf(format, a...)
	return Status{Code: CodeOK, Text: text}
}

// None

// Nonef formats a message and returns a none status.
func Nonef(format string, a ...any) Status {
	text := fmt.Sprintf(format, a...)
	return Status{Code: CodeNone, Text: text}
}

// Test

// Test returns a test status.
func Test(message string) Status {
	return Status{Code: CodeTest, Text: message}
}

// Testf formats a message and returns a test status.
func Testf(format string, a ...any) Status {
	text := fmt.Sprintf(format, a...)
	return Status{Code: CodeTest, Text: text}
}
