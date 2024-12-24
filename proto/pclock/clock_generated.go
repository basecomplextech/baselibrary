package pclock

import (
	"github.com/basecomplextech/baselibrary/alloc"
	"github.com/basecomplextech/baselibrary/async"
	"github.com/basecomplextech/baselibrary/bin"
	"github.com/basecomplextech/baselibrary/buffer"
	"github.com/basecomplextech/baselibrary/pools"
	"github.com/basecomplextech/baselibrary/ref"
	"github.com/basecomplextech/baselibrary/status"
	"github.com/basecomplextech/spec"
	"github.com/basecomplextech/spec/proto/prpc"
	"github.com/basecomplextech/spec/rpc"
)

var (
	_ alloc.Buffer
	_ async.Context
	_ bin.Bin128
	_ buffer.Buffer
	_ spec.MessageTable
	_ pools.Pool[any]
	_ ref.Ref
	_ rpc.Client
	_ prpc.Request
	_ spec.Type
	_ status.Status
)

// HLTimestamp

type HLTimestamp struct {
	Wall  int64  `json:"Wall"`
	Logic uint32 `json:"Logic"`
}

func OpenHLTimestamp(b []byte) HLTimestamp {
	s, _, _ := DecodeHLTimestamp(b)
	return s
}

func DecodeHLTimestamp(b []byte) (s HLTimestamp, size int, err error) {
	size, err = s.Decode(b)
	return s, size, err
}

func EncodeHLTimestampTo(b buffer.Buffer, s HLTimestamp) (int, error) {
	return s.EncodeTo(b)
}

func (s *HLTimestamp) Decode(b []byte) (size int, err error) {
	dataSize, size, err := spec.DecodeStruct(b)
	if err != nil || size == 0 {
		return
	}

	b = b[len(b)-size:]
	n := size - dataSize
	off := len(b) - n

	// Decode in reverse order

	s.Logic, n, err = spec.DecodeUint32(b[:off])
	if err != nil {
		return
	}
	off -= n

	s.Wall, n, err = spec.DecodeInt64(b[:off])
	if err != nil {
		return
	}
	off -= n

	return size, err
}

func (s HLTimestamp) EncodeTo(b buffer.Buffer) (int, error) {
	var dataSize, n int
	var err error

	n, err = spec.EncodeInt64(b, s.Wall)
	if err != nil {
		return 0, err
	}
	dataSize += n

	n, err = spec.EncodeUint32(b, s.Logic)
	if err != nil {
		return 0, err
	}
	dataSize += n

	n, err = spec.EncodeStruct(b, dataSize)
	if err != nil {
		return 0, err
	}
	return dataSize + n, nil
}
