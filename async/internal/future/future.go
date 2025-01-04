// Copyright 2021 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package future

import (
	"github.com/basecomplextech/baselibrary/status"
)

// Future represents a result available in the future.
type Future[T any] interface {
	// Done returns true if the future is complete.
	Done() bool

	// Wait returns a channel which is closed when the future is complete.
	Wait() <-chan struct{}

	// Result returns a value and a status.
	Result() (T, status.Status)

	// Status returns a status or none.
	Status() status.Status
}

// FutureDyn is a future interface without generics, i.e. Future[?].
type FutureDyn interface {
	// Done returns true if the future is complete.
	Done() bool

	// Wait returns a channel which is closed when the future is complete.
	Wait() <-chan struct{}

	// Status returns a status or none.
	Status() status.Status
}
