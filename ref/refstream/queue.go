package refstream

import (
	"github.com/basecomplextech/baselibrary/async"
	"github.com/basecomplextech/baselibrary/ref"
	"github.com/basecomplextech/baselibrary/ref/refqueue"
)

var _ async.StreamListener[ref.Ref] = (*queue[ref.Ref])(nil)

type queue[T ref.Ref] struct {
	refqueue.Queue[T]
	unsub func()
}

func newQueue[T ref.Ref]() *queue[T] {
	return &queue[T]{
		Queue: refqueue.New[T](),
	}
}

// OnStreamMessage is called when a new message is available on the stream.
func (q *queue[T]) OnStreamMessage(msg T) {
	q.Queue.Push(msg)
}
