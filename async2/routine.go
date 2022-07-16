package async2

import "github.com/epochtimeout/baselibrary/status"

// Routine is an async routine with a future result which wraps a goroutine.
type Routine[T any] interface {
	Result[T]

	// Stop requests the routine to stop and returns the wait channel.
	Stop() <-chan struct{}

	// Kill requests the routine to die and returns the wait channel.
	Kill() <-chan struct{}
}

// Run runs a function in an async routine.
func Run(fn func(stop <-chan struct{}) status.Status) Routine[struct{}] {
	return nil
}

// Call calls a function in an async routine.
func Call[T any](fn func(stop <-chan struct{}) (T, status.Status)) Routine[T] {
	return nil
}

// StopAll stops all routines.
func StopAll[R Routine[T], T any](routines ...R) {
	for _, r := range routines {
		r.Stop()
	}
}

// internal
