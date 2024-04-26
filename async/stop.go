package async

// Stopper requests a service or routine to stop.
type Stopper interface {
	// Stop requests a service or routine to stop and returns a wait channel.
	Stop() <-chan struct{}
}

// StopWaiter requests a service or routine to stop and awaits its stop.
type StopWaiter interface {
	Stopper
	Waiter
}

// Utility

// StopAll stops all routines.
func StopAll[W Stopper](w ...W) {
	for _, w := range w {
		w.Stop()
	}
}

// StopWait stops a routine and awaits it stop.
//
// Usually used with defer:
//
//	worker := runWorker()
//	defer StopWait(worker)
func StopWait(w Stopper) {
	<-w.Stop()
}

// StopWaitAll cancels and awaits all operations.
func StopWaitAll[W StopWaiter](w ...W) {
	for _, w := range w {
		w.Stop()
	}

	for _, w := range w {
		<-w.Wait()
	}
}
