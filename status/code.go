package status

type Code string

// General codes
const (
	// CodeUndefined is an empty status code.
	CodeUndefined Code = ""

	// CodeOK indicates that an operation completed successfully.
	CodeOK Code = "ok"

	// CodeNone indicates that no status code is available.
	CodeNone Code = "none"

	// CodeClosed indicates that an object is closed and cannot be used anymore.
	CodeClosed Code = "closed"

	// CodeTerminal indicates that an operation or a state is terminal and cannot be continued.
	CodeTerminal Code = "terminal"
)

// Error codes
const (
	// CodeError is a generic error status code.
	CodeError Code = "error"

	// CodeCorruption indicates any data corruption.
	CodeCorruption Code = "corruption"

	// CodeIOError indicates that an I/O error occurred, the operation can be retried later.
	CodeIOError Code = "io_error"

	// CodeNotFound indicates that an object is not found.
	CodeNotFound Code = "not_found"
)

// Unavailable codes
const (
	// CodeCancelled indicates that an operation was cancelled or stopped on a request.
	CodeCancelled Code = "cancelled"

	// CodeTimeout indicates that an operation timed out.
	CodeTimeout Code = "timeout"

	// CodeUnavailable indicates that a service is temporarily unavailable, the operation can be retried.
	CodeUnavailable Code = "unavailable"
)

// Iteration/streaming codes
const (
	// CodeStart indicates a file/channel/stream start.
	CodeStart Code = "start"

	// CodeEnd indicates a file/channel/stream end.
	CodeEnd Code = "end"

	// CodeWait indicates that the caller should wait for the next events/messages/etc.
	CodeWait Code = "wait"
)

// Test codes
const (
	// CodeTest is a status code for testing.
	CodeTest Code = "test"
)
