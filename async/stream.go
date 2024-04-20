package async

// Stream is an asynchronous stream of messages.
type Stream[T any] interface {
	// Filter returns a new stream that only contains elements that satisfy the predicate.
	Filter(fn func(T) bool) Stream[T]

	// Listen adds a listener to the stream, and returns an unsubscribe function.
	Listen(ln StreamListener[T]) (unsub func())

	// Subscribe returns a queue subscribed to the stream, free the queue to unsubscribe.
	Subscribe() Queue[T]
}

// StreamListener receives messages from a stream.
type StreamListener[T any] interface {
	// OnStreamMessage is called when a new message is available on the stream.
	OnStreamMessage(T)
}

// StreamSource is a stream which can be written to.
type StreamSource[T any] interface {
	Stream[T]

	// Send sends a message to the stream.
	Send(T)
}

// NewStreamSource returns a new stream source.
func NewStreamSource[T any]() StreamSource[T] {
	return newStreamSource[T]()
}
