package async

// Waiter waits for an operation to complete.
type Waiter interface {
	Wait() <-chan struct{}
}
