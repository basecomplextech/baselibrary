package async2

import "github.com/epochtimeout/baselibrary/status"

// Result is a generic result in the future.
type Result[T any] interface {
	// Wait awaits the result.
	Wait() <-chan struct{}

	// Result returns the result and its status or zer.
	Result() (T, status.Status)

	// Status returns the result status.
	Status() status.Status
}
