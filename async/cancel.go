package async

// Canceller requests an operation to cancel.
type Canceller interface {
	// Cancel requests an operation to cancel and returns a wait channel.
	Cancel() <-chan struct{}
}

// Waiter awaits an operation completion.
type Waiter interface {
	// Wait returns a channel which is closed when the operation is complete.
	Wait() <-chan struct{}
}

// CancelWaiter requests an operation to cancel and awaits its completion.
type CancelWaiter interface {
	Canceller
	Waiter
}

// CancelWait cancels and awaits an operation.
func CancelWait(w CancelWaiter) {
	<-w.Cancel()
}
