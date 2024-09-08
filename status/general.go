// Copyright 2022 Ivan Korobkov. All rights reserved.

package status

import "fmt"

var (
	OK   = New(CodeOK, "")
	None = New(CodeNone, "")
)

// OKf formats a message and returns an ok status.
func OKf(format string, a ...any) Status {
	msg := fmt.Sprintf(format, a...)
	return Status{Code: CodeOK, Message: msg}
}

// None

// Nonef formats a message and returns a none status.
func Nonef(format string, a ...any) Status {
	msg := fmt.Sprintf(format, a...)
	return Status{Code: CodeNone, Message: msg}
}

// Test

// Test returns a test status.
func Test(message string) Status {
	return Status{Code: CodeTest, Message: message}
}

// Testf formats a message and returns a test status.
func Testf(format string, a ...any) Status {
	msg := fmt.Sprintf(format, a...)
	return Status{Code: CodeTest, Message: msg}
}
