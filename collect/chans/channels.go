// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package chans

// Closed returns a closed void channel.
func Closed() chan struct{} {
	return closed
}

// private

var closed = func() chan struct{} {
	ch := make(chan struct{})
	close(ch)
	return ch
}()
