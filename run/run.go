package run

// Call calls a function in another thread and returns its result.
func Call(fn func(stop <-chan struct{}) (interface{}, error)) Thread {
	th := newThread()

	go func() {
		defer func() {
			th.catch(recover())
		}()

		result, err := fn(th.stopCh)
		th.complete(result, err)
	}()

	return th
}

// Run runs a function in another thread.
func Run(fn func(stop <-chan struct{}) error) Thread {
	th := newThread()

	go func() {
		defer func() {
			th.catch(recover())
		}()

		err := fn(th.stopCh)
		th.complete(nil, err)
	}()

	return th
}

// StopAndWait stops a thread and waits for its result.
func StopAndWait(th Thread) error {
	f := th.Stop()
	<-f.Done()
	return f.Err()
}
