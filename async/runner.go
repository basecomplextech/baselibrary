// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package async

// Runner is an interface that provides a run method.
type Runner interface {
	// Run runs the function.
	Run()
}

// RunnerFunc returns a new runner from a function.
func RunnerFunc(fn func()) Runner {
	return runnerFunc(fn)
}

// private

type runnerFunc func()

func (r runnerFunc) Run() {
	r()
}
