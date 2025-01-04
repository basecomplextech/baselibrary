// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package internal

// Stopper stops a routine and awaits its stop.
type Stopper interface {
	// Wait returns a channel which is closed when the future is complete.
	Wait() <-chan struct{}

	// Stop requests the routine to stop and returns a wait channel.
	Stop() <-chan struct{}
}

// StopAll stops all routines, but does not await their stop.
func StopAll[R Stopper](routines ...R) {
	for _, r := range routines {
		r.Stop()
	}
}

// StopWait stops a routine and awaits it stop.
//
// Usually used with defer:
//
//	w := runWorker()
//	defer StopWait(w)
func StopWait[R Stopper](r R) {
	<-r.Stop()
}

// StopWaitAll stops and awaits all routines.
//
// Usually used with defer:
//
//	w0 := runWorker()
//	w1 := runWorker()
//	defer StopWaitAll(w0, w1)
func StopWaitAll[R Stopper](routines ...R) {
	for _, r := range routines {
		r.Stop()
	}

	for _, r := range routines {
		<-r.Wait()
	}
}
