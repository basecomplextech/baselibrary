package arena

import "github.com/complex1tech/baselibrary/alloc/internal/heap"

func Test() Arena {
	h := heap.New()
	return New(h)
}
