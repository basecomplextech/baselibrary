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
