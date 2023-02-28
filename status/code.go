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
	// CodeError is a generic error status code.
	CodeError Code = "error"

	// CodeIOError indicates that an I/O error occurred, the operation can be retried later.
	CodeIOError Code = "io_error"

	// CodeCorrupted indicates any data corruption or loss.
	CodeCorrupted Code = "corruption"

	// CodeFatal indicates that a fatal error occurred, the operation cannot be retried.
	CodeFatal Code = "fatal"
)

// App/client class
const (
	// CodeInvalid indicates that a client request/operation/argument is invalid.
	CodeInvalid Code = "invalid"

	// CodeNotFound indicates that an object is not found.
	CodeNotFound Code = "notfound"

	// CodeForbidden indicates that an operation is forbidden.
	CodeForbidden Code = "forbidden"

	// CodeUnauthorized indicates that a user is not authorized.
	CodeUnauthorized Code = "unauthorized"
)

// Unavailable class
const (
	// CodeAborted indicates that an operation was aborted or rollbacked.
	CodeAborted Code = "aborted"

	// CodeClosed indicates that an object is closed and cannot be used anymore.
	CodeClosed Code = "closed"

	// CodeCancelled indicates that an operation was cancelled or stopped on a request.
	CodeCancelled Code = "cancelled"

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
