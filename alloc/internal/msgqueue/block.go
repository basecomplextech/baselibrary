package msgqueue

import (
	"encoding/binary"

	"github.com/basecomplextech/baselibrary/alloc/internal/heap"
)

type block struct {
	*heap.Block

	read    int  // read index
	started bool // used in next() to return the first message
}

// rem returns the remaining free space in bytes.
func (b *block) rem() int {
	return b.Rem()
}

// unread returns the number of unread bytes.
func (b *block) unread() int {
	len := b.Len()
	return len - b.read
}

// current returns the current message or false.
func (b *block) current() ([]byte, bool) {
	unread := b.unread()
	if unread == 0 {
		return nil, false
	}

	p := b.Bytes()
	p = p[b.read:]

	size := binary.BigEndian.Uint32(p)
	msg := p[4 : 4+size]

	if !b.started {
		b.started = true
	}
	return msg, true
}

// next moves to the next message, returns false if no more unread messages.
func (b *block) next() ([]byte, bool) {
	if !b.started {
		return b.current()
	}

	unread := b.unread()
	if unread == 0 {
		return nil, false
	}

	p := b.Bytes()
	p = p[b.read:]

	size := binary.BigEndian.Uint32(p)
	b.read += int(4 + size)

	return b.current()
}

// write writes a message to the block or returns panics if not enough space.
func (b *block) write(msg []byte) {
	rem := b.Rem()
	size := len(msg)

	n := 4 + size
	if rem < n {
		panic("no more space in block")
	}

	p := b.Grow(n)
	binary.BigEndian.PutUint32(p, uint32(size))
	copy(p[4:], msg)
}

// reset resets the block read and write indexes.
func (b *block) reset() {
	b.Reset()

	b.read = 0
	b.started = false
}
