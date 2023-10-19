package async

// Canceller requests an operation to cancel.
type Canceller interface {
	// Cancel requests an operation to cancel and returns a wait channel.
	Cancel() <-chan struct{}
}

// CancelWaiter requests an operation to cancel and awaits its completion.
type CancelWaiter interface {
	Canceller
	Waiter
}

// Utility

// CancelAll cancels all operations.
func CancelAll[W Canceller](w ...W) {
	for _, w := range w {
		w.Cancel()
	}
}

// CancelWait cancels and awaits an operation.
//
// Usually used with defer:
//
//	worker := runWorker()
//	defer CancelWait(worker)
func CancelWait(w Canceller) {
	<-w.Cancel()
}

// CancelWaitAll cancels and awaits all operations.
func CancelWaitAll[W CancelWaiter](w ...W) {
	for _, w := range w {
		w.Cancel()
	}

	for _, w := range w {
		<-w.Wait()
	}
}
