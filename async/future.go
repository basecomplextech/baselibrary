package async

// Future returns an async result in the future.
type Future[T any] interface {
	// Err returns the future error or nil.
	Err() error

	// Wait returns a channel which is closed on a future completion.
	Wait() <-chan struct{}

	// Result returns a future result and an error.
	Result() (T, error)
}

// Await awaits a future and returns its result.
func Await[T any](f Future[T]) (T, error) {
	<-f.Wait()
	return f.Result()
}
