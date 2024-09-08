// Copyright 2023 Ivan Korobkov. All rights reserved.

package arena

import "github.com/basecomplextech/baselibrary/alloc/internal/heap"

func Test() Arena {
	h := heap.New()
	return newArena(h)
}
