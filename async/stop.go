// Copyright 2025 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package async

import "github.com/basecomplextech/baselibrary/async/internal"

// Stopper stops a routine and awaits its stop.
type Stopper = internal.Stopper

// StopAll stops all routines, but does not await their stop.
func StopAll[R Stopper](routines ...R) {
	internal.StopAll(routines...)
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
	internal.StopWaitAll(routines...)
}
