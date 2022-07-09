package async

import "github.com/epochtimeout/baselibrary/status"

// Future returns an async result in the future.
type Future[T any] interface {
	// Wait returns a channel which is closed on a future completion.
	Wait() <-chan struct{}

	// Result returns a future result and an error.
	Result() (T, status.Status)

	// Status returns a future status when completed or an empty status.
	Status() status.Status
}

// Await awaits a future and returns its result.
func Await[T any](f Future[T]) (T, status.Status) {
	<-f.Wait()
	return f.Result()
}
