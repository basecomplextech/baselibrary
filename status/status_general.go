package status

import "fmt"

// OKf formats a message and returns an ok status.
func OKf(format string, a ...interface{}) Status {
	text := fmt.Sprintf(format, a...)
	return Status{Code: CodeOK, Text: text}
}

// Terminal

// Terminal returns a terminal status.
func Terminal(text string) Status {
	return Status{Code: CodeTerminal, Text: text}
}

// Terminalf formats a message and returns a terminal status.
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

// Test

// Test returns a test status.
func Test(message string) Status {
	return Status{Code: CodeTest, Text: message}
}

// Testf formats a message and returns a test status.
func Testf(format string, a ...interface{}) Status {
	text := fmt.Sprintf(format, a...)
	return Status{Code: CodeTest, Text: text}
}
