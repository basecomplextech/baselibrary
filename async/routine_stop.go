package async

// StopAll stops all routines, but does not await their stop.
func StopAll[R RoutineDyn](routines ...R) {
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
func StopWait[R RoutineDyn](r R) {
	<-r.Stop()
}

// StopWaitAll stops and awaits all routines.
//
// Usually used with defer:
//
//	w0 := runWorker()
//	w1 := runWorker()
//	defer StopWaitAll(w0, w1)
func StopWaitAll[R RoutineDyn](routines ...R) {
	for _, r := range routines {
		r.Stop()
	}

	for _, r := range routines {
		<-r.Wait()
	}
}
