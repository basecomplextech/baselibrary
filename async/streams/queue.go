package streams

import "github.com/basecomplextech/baselibrary/async"

var (
	_ async.Queue[any] = (*queue[any])(nil)
	_ Listener[any]    = (*queue[any])(nil)
)

type queue[T any] struct {
	async.Queue[T]
	unsub func()
}

func newQueue[T any]() *queue[T] {
	q := async.NewQueue[T]()
	return &queue[T]{Queue: q}
}

// OnStreamMessage is called when a new message is available on the stream.
func (q *queue[T]) OnStreamMessage(msg T) {
	q.Queue.Push(msg)
}

// Free unsubscribes from the stream.
func (q *queue[T]) Free() {
	q.unsub()
	q.Queue.Free()
}
