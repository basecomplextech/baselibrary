// Copyright 2023 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package arena

import "github.com/basecomplextech/baselibrary/alloc/internal/heap"

func Test() Arena {
	h := heap.New()
	return newArena(h)
}

func TestMutex() Arena {
	h := heap.New()
	return newMutexArena(h)
}
