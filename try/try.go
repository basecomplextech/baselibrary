package try

// Run runs a function, recovers on a panic.
func Run(fn func() error) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = Recover(e)
		}
	}()

	return fn()
}

// Call calls a function and returns its result, recovers on a panic.
func Call[T any](fn func() (T, error)) (result T, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = Recover(e)
		}
	}()

	return fn()
}
