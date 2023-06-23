package status

import "fmt"

// NotFound

// NotFound returns a not found status.
func NotFound(msg string) Status {
	return Status{Code: CodeNotFound, Message: msg}
}

// NotFoundf formats a message and returns a not found status.
func NotFoundf(format string, a ...any) Status {
	msg := fmt.Sprintf(format, a...)
	return Status{Code: CodeNotFound, Message: msg}
}

// Forbidden

// Forbidden returns a forbidden status.
func Forbidden(msg string) Status {
	return Status{Code: CodeForbidden, Message: msg}
}

// Forbiddenf formats a message and returns a forbidden status.
func Forbiddenf(format string, a ...any) Status {
	msg := fmt.Sprintf(format, a...)
	return Status{Code: CodeForbidden, Message: msg}
}

// Unauthorized

// Unauthorized returns an unauthorized status.
func Unauthorized(msg string) Status {
	return Status{Code: CodeUnauthorized, Message: msg}
}

// Unauthorizedf formats a message and returns an unauthorized status.
func Unauthorizedf(format string, a ...any) Status {
	msg := fmt.Sprintf(format, a...)
	return Status{Code: CodeUnauthorized, Message: msg}
}
