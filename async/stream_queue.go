package async

var (
	_ Queue[any]          = (*streamQueue[any])(nil)
	_ StreamListener[any] = (*streamQueue[any])(nil)
)

type streamQueue[T any] struct {
	Queue[T]
	unsub func()
}

func newStreamQueue[T any]() *streamQueue[T] {
	queue := newQueue[T]()
	return &streamQueue[T]{Queue: queue}
}

// OnStreamMessage is called when a new message is available on the stream.
func (q *streamQueue[T]) OnStreamMessage(msg T) {
	q.Queue.Push(msg)
}

// Free unsubscribes from the stream.
func (q *streamQueue[T]) Free() {
	q.unsub()
	q.Queue.Free()
}
