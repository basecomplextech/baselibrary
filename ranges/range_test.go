package ranges

import (
	"bytes"
	"testing"

	"github.com/zeebo/assert"
)

// Contains

func TestRange_Contains__should_return_whether_key_is_inside_this_range(t *testing.T) {
	r := Range[int]{
		Start: 5,
		End:   7,
	}

	compare := func(a, b int) int {
		return a - b
	}

	assert.False(t, r.Contains(0, compare))
	assert.False(t, r.Contains(4, compare))
	assert.True(t, r.Contains(5, compare))
	assert.True(t, r.Contains(6, compare))
	assert.True(t, r.Contains(7, compare))
	assert.False(t, r.Contains(8, compare))
	assert.False(t, r.Contains(10, compare))
}

// Overlaps

func TestRange_Overlaps__should_return_whether_another_range_overlaps_this_range(t *testing.T) {
	r0 := Range[int]{
		Start: 5,
		End:   8,
	}

	compare := func(a, b int) int {
		return a - b
	}

	assert.False(t, r0.Overlaps(Range[int]{0, 4}, compare))
	assert.True(t, r0.Overlaps(Range[int]{0, 5}, compare))
	assert.True(t, r0.Overlaps(Range[int]{5, 5}, compare))
	assert.True(t, r0.Overlaps(Range[int]{5, 8}, compare))
	assert.True(t, r0.Overlaps(Range[int]{6, 7}, compare))
	assert.True(t, r0.Overlaps(Range[int]{8, 10}, compare))
	assert.False(t, r0.Overlaps(Range[int]{9, 10}, compare))

	assert.True(t, r0.Overlaps(Range[int]{0, 10}, compare))
	assert.True(t, r0.Overlaps(Range[int]{4, 9}, compare))
}

// Inside

func TestRange_Inside(t *testing.T) {
	r0 := Range[int]{
		Start: 3,
		End:   7,
	}

	compare := func(a, b int) int {
		return a - b
	}

	assert.True(t, r0.Inside(Range[int]{3, 7}, compare))
	assert.True(t, r0.Inside(Range[int]{2, 8}, compare))
	assert.True(t, r0.Inside(Range[int]{0, 10}, compare))

	assert.False(t, r0.Inside(Range[int]{0, 0}, compare))
	assert.False(t, r0.Inside(Range[int]{4, 6}, compare))
	assert.False(t, r0.Inside(Range[int]{2, 6}, compare))
	assert.False(t, r0.Inside(Range[int]{4, 8}, compare))
}

// Expand

func TestRange_Expand__should_expand_range(t *testing.T) {
	r := Range[int]{3, 7}
	cmp := func(a, b int) int { return a - b }

	assert.Equal(t, Range[int]{1, 7}, r.Expand(Range[int]{1, 2}, cmp))
	assert.Equal(t, Range[int]{3, 9}, r.Expand(Range[int]{3, 9}, cmp))
	assert.Equal(t, Range[int]{0, 10}, r.Expand(Range[int]{0, 10}, cmp))
}

// ExpandBinary

func TestExpandBinary__should_expand_binary_range_skip_nil_values(t *testing.T) {
	r := Range[[]byte]{[]byte{3}, []byte{7}}
	cmp := bytes.Compare

	assert.Equal(t,
		Range[[]byte]{[]byte{3}, []byte{7}},
		ExpandBinary(r, Range[[]byte]{nil, nil}, cmp))
	assert.Equal(t,
		Range[[]byte]{[]byte{3}, []byte{9}},
		ExpandBinary(r, Range[[]byte]{nil, []byte{9}}, cmp))
	assert.Equal(t,
		Range[[]byte]{[]byte{0}, []byte{7}},
		ExpandBinary(r, Range[[]byte]{[]byte{0}, nil}, cmp))

	assert.Equal(t,
		Range[[]byte]{[]byte{1}, []byte{7}},
		ExpandBinary(r, Range[[]byte]{[]byte{1}, []byte{2}}, cmp))
	assert.Equal(t,
		Range[[]byte]{[]byte{1}, []byte{7}},
		ExpandBinary(r, Range[[]byte]{[]byte{1}, []byte{2}}, cmp))
	assert.Equal(t,
		Range[[]byte]{[]byte{1}, []byte{7}},
		ExpandBinary(r, Range[[]byte]{[]byte{1}, []byte{2}}, cmp))
	assert.Equal(t,
		Range[[]byte]{[]byte{3}, []byte{9}},
		ExpandBinary(r, Range[[]byte]{[]byte{3}, []byte{9}}, cmp))
	assert.Equal(t,
		Range[[]byte]{[]byte{0}, []byte{10}},
		ExpandBinary(r, Range[[]byte]{[]byte{0}, []byte{10}}, cmp))

	assert.Equal(t,
		Range[[]byte]{[]byte{0}, []byte{10}},
		ExpandBinary(Range[[]byte]{nil, nil}, Range[[]byte]{[]byte{0}, []byte{10}}, cmp))
	assert.Equal(t,
		Range[[]byte]{nil, []byte{10}},
		ExpandBinary(Range[[]byte]{nil, nil}, Range[[]byte]{nil, []byte{10}}, cmp))
	assert.Equal(t,
		Range[[]byte]{[]byte{0}, nil},
		ExpandBinary(Range[[]byte]{nil, nil}, Range[[]byte]{[]byte{0}, nil}, cmp))
}
