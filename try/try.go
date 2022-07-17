package try

import "github.com/epochtimeout/baselibrary/errors2"

// Run runs a function, recovers on a panic.
func Run(fn func() error) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors2.Recover(e)
		}
	}()

	return fn()
}

// Execute executes a function and returns its result, recovers on a panic.
func Execute[T any](fn func() (T, error)) (result T, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors2.Recover(e)
		}
	}()
	return fn()
}
