package async

// Canceller tries to cancel an operation.
type Canceller interface {
	Cancel() <-chan struct{}
}

// CancelWaiter cancels an operation and awaits it.
type CancelWaiter interface {
	Canceller
	Waiter
}

// Cancel cancels all operations.
func Cancel(cc ...Canceller) {
	for _, c := range cc {
		c.Cancel()
	}
}

// CancelWait cancels all operations and awaits them.
func CancelWait(cc ...CancelWaiter) {
	for _, c := range cc {
		c.Cancel()
	}
	for _, c := range cc {
		c.Wait()
	}
}
