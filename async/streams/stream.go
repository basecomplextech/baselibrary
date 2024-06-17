package streams

import "github.com/basecomplextech/baselibrary/async"

// Stream is an asynchronous stream of messages.
type Stream[T any] interface {
	// Filter returns a new stream that only contains elements that satisfy the predicate.
	Filter(fn func(T) bool) Stream[T]

	// Listen adds a listener to the stream, and returns an unsubscribe function.
	Listen(ln Listener[T]) (unsub func())

	// Subscribe returns a queue subscribed to the stream, free the queue to unsubscribe.
	Subscribe() async.Queue[T]
}

// Listener receives messages from a stream.
type Listener[T any] interface {
	// OnStreamMessage is called when a new message is available on the stream.
	OnStreamMessage(T)
}
