package status

import "fmt"

// NotFound

// NotFound returns a not found status.
func NotFound(text string) Status {
	return Status{Code: CodeNotFound, Text: text}
}

// NotFoundf formats a message and returns a not found status.
func NotFoundf(format string, a ...any) Status {
	text := fmt.Sprintf(format, a...)
	return Status{Code: CodeNotFound, Text: text}
}

// Forbidden

// Forbidden returns a forbidden status.
func Forbidden(text string) Status {
	return Status{Code: CodeForbidden, Text: text}
}

// Forbiddenf formats a message and returns a forbidden status.
func Forbiddenf(format string, a ...any) Status {
	text := fmt.Sprintf(format, a...)
	return Status{Code: CodeForbidden, Text: text}
}

// Unauthorized

// Unauthorized returns an unauthorized status.
func Unauthorized(text string) Status {
	return Status{Code: CodeUnauthorized, Text: text}
}

// Unauthorizedf formats a message and returns an unauthorized status.
func Unauthorizedf(format string, a ...any) Status {
	text := fmt.Sprintf(format, a...)
	return Status{Code: CodeUnauthorized, Text: text}
}
