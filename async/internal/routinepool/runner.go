// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package routinepool

type Runner interface {
	// Run runs the function.
	Run()
}

// private

var _ Runner = (runnerFunc)(nil)

type runnerFunc func()

func (f runnerFunc) Run() {
	f()
}
