package exc

type Code int

const (
	CodeUndefined Code = 0

	// Internal error is a general internal error.
	CodeInternal Code = 1

	// Illegal argument means that a passed argument is illegal or invalid.
	// For example, it can indicate an invalid start/end indexes or an invalid
	// name.
	CodeIllegalArg Code = 2

	// Illegal operation means that a function or method cannot be executed.
	// For example, it can indicate an operation on a closed file or on a stopped
	// thread.
	CodeIllegalOp Code = 3

	// Invalid state means that the system state is illegal or invalid.
	// For example, it can indicate that a file is not found or that a config is
	// invalid.
	CodeInvalidState Code = 4

	// Cancelled operation error means that an operation has been cancelled.
	CodeCancelled Code = 5

	// Data corrupted specifies any data corruption error.
	CodeCorrupted Code = 6

	// Unavailable means that a service or object is not available for usage yet.
	// The operation can be retried later.
	CodeUnavailable Code = 7

	// Not found specifies that a requested object does not exist.
	CodeNotFound Code = 8

	// Test is a test error.
	CodeTest Code = 1000
)

func (c Code) String() string {
	switch c {
	case CodeInternal:
		return "internal"
	case CodeIllegalArg:
		return "illegal_arg"
	case CodeIllegalOp:
		return "illegal_op"
	case CodeInvalidState:
		return "invalid_state"
	case CodeCancelled:
		return "cancelled"
	case CodeCorrupted:
		return "corrupted"
	case CodeUnavailable:
		return "unavailable"
	case CodeNotFound:
		return "not_found"
	case CodeTest:
		return "test"
	}
	return ""
}
