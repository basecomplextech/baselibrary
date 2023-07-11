package ranges

import (
	"testing"

	"github.com/complex1tech/baselibrary/compare"
	"github.com/stretchr/testify/assert"
)

// ExpandBytes

func TestExpandBytes__should_expand_binary_range_skip_nil_values(t *testing.T) {
	r := Range[[]byte]{[]byte{3}, []byte{7}}
	cmp := compare.Bytes

	assert.Equal(t,
		Range[[]byte]{[]byte{3}, []byte{7}},
		ExpandBytes(r, Range[[]byte]{nil, nil}, cmp))
	assert.Equal(t,
		Range[[]byte]{[]byte{3}, []byte{9}},
		ExpandBytes(r, Range[[]byte]{nil, []byte{9}}, cmp))
	assert.Equal(t,
		Range[[]byte]{[]byte{0}, []byte{7}},
		ExpandBytes(r, Range[[]byte]{[]byte{0}, nil}, cmp))

	assert.Equal(t,
		Range[[]byte]{[]byte{1}, []byte{7}},
		ExpandBytes(r, Range[[]byte]{[]byte{1}, []byte{2}}, cmp))
	assert.Equal(t,
		Range[[]byte]{[]byte{1}, []byte{7}},
		ExpandBytes(r, Range[[]byte]{[]byte{1}, []byte{2}}, cmp))
	assert.Equal(t,
		Range[[]byte]{[]byte{1}, []byte{7}},
		ExpandBytes(r, Range[[]byte]{[]byte{1}, []byte{2}}, cmp))
	assert.Equal(t,
		Range[[]byte]{[]byte{3}, []byte{9}},
		ExpandBytes(r, Range[[]byte]{[]byte{3}, []byte{9}}, cmp))
	assert.Equal(t,
		Range[[]byte]{[]byte{0}, []byte{10}},
		ExpandBytes(r, Range[[]byte]{[]byte{0}, []byte{10}}, cmp))

	assert.Equal(t,
		Range[[]byte]{[]byte{0}, []byte{10}},
		ExpandBytes(Range[[]byte]{nil, nil}, Range[[]byte]{[]byte{0}, []byte{10}}, cmp))
	assert.Equal(t,
		Range[[]byte]{nil, []byte{10}},
		ExpandBytes(Range[[]byte]{nil, nil}, Range[[]byte]{nil, []byte{10}}, cmp))
	assert.Equal(t,
		Range[[]byte]{[]byte{0}, nil},
		ExpandBytes(Range[[]byte]{nil, nil}, Range[[]byte]{[]byte{0}, nil}, cmp))
}
