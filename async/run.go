package async

// Run runs a function in a goroutine, recovers panics, returns a promise.
func Run(fn func(stop <-chan struct{}) (Status, interface{}, error)) Promise {
	p := newPromise()

	go func() {
		defer p.Exit()
		defer func() {
			e := recover()
			if e == nil {
				return
			}

			err := Recover(e)
			p.Fail(err)
		}()

		stop := p.Stop()
		status, result, err := fn(stop)
		switch status {
		case StatusOK:
			p.OK(result)
		case StatusError:
			p.Fail(err)
		}
	}()

	return p
}
