package status

type Code string

// General class
const (
	// CodeNone indicates an undefined status code.
	CodeNone Code = ""

	// CodeOK indicates that an operation completed successfully.
	CodeOK Code = "ok"

	// CodeTest is a status code for testing.
	CodeTest Code = "test"
)

// Error class
const (
	// CodeError indicates an internal general purpose error.
	CodeError Code = "error"

	// CodeCorrupted indicates any data corruption or loss.
	CodeCorrupted Code = "corrupted"

	// CodeExternalError indicates an external error, i.e. an invalid argument, validation error, etc.
	CodeExternalError Code = "external_error"
)

// Invalid class
const (
	// CodeNotFound indicates that an object is not found.
	CodeNotFound Code = "not_found"

	// CodeForbidden indicates that an operation is forbidden.
	CodeForbidden Code = "forbidden"

	// CodeUnauthorized indicates that a user is not authorized.
	CodeUnauthorized Code = "unauthorized"
)

// Unavailable class
const (
	// CodeClosed indicates that an object is closed and cannot be used anymore.
	CodeClosed Code = "closed"

	// CodeCancelled indicates that an operation was cancelled or stopped on a request.
	CodeCancelled Code = "cancelled"

	// CodeRollback indicates that an operation was rolled backed.
	CodeRollback Code = "rollback" // TODO: Maybe remove

	// CodeTimeout indicates that an operation timed out.
	CodeTimeout Code = "timeout"

	// CodeUnavailable indicates that a service is temporarily unavailable, the operation can be retried.
	CodeUnavailable Code = "unavailable"

	// CodeUnsupported indicates that an operation is not supported or not implemented.
	CodeUnsupported Code = "unsupported"
)

// Iteration/streaming class
const (
	// CodeEnd indicates a file/channel/stream end.
	CodeEnd Code = "end"

	// CodeWait indicates that the caller should wait for the next events/messages/etc.
	CodeWait Code = "wait"
)
