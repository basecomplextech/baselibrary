package status

type Code string

const (
	// CodeUndefined is an empty status code.
	CodeUndefined Code = ""

	// CodeOK indicates that an operation completed successfully.
	CodeOK Code = "ok"

	// CodeError is a generic error status code.
	CodeError Code = "error"

	// CodeStopped indicates that an operation was stopped or cancelled on a request.
	CodeStopped Code = "stopped"

	// CodeTimeout indicates that an operation timed out.
	CodeTimeout Code = "timeout"

	// CodeTerminal indicates that an operation or a state is terminal and cannot be continued.
	CodeTerminal Code = "terminal"

	// CodeUnavailable indicates that a service is temporarily unavailable, the operation can be retried.
	CodeUnavailable Code = "unavailable"

	// CodeNotFound indicates that an object is not found.
	CodeNotFound Code = "not_found"

	// CodeCorrupted indicates any data corruption.
	CodeCorrupted Code = "corrupted"

	// CodeTest is a status code for testing.
	CodeTest Code = "test"
)
