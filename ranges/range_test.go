package ranges

import (
	"testing"

	"github.com/epochtimeout/baselibrary/compare"
	"github.com/stretchr/testify/assert"
)

// Contains

func TestRange_Contains__should_return_whether_key_is_inside_this_range(t *testing.T) {
	r := Range[int]{
		Start: 5,
		End:   7,
	}

	cmp := compare.Int
	assert.False(t, r.Contains(0, cmp))
	assert.False(t, r.Contains(4, cmp))
	assert.True(t, r.Contains(5, cmp))
	assert.True(t, r.Contains(6, cmp))
	assert.True(t, r.Contains(7, cmp))
	assert.False(t, r.Contains(8, cmp))
	assert.False(t, r.Contains(10, cmp))
}

// Overlaps

func TestRange_Overlaps__should_return_whether_another_range_overlaps_this_range(t *testing.T) {
	r0 := Range[int]{
		Start: 5,
		End:   8,
	}

	cmp := compare.Int
	assert.False(t, r0.Overlaps(Range[int]{0, 4}, cmp))
	assert.True(t, r0.Overlaps(Range[int]{0, 5}, cmp))
	assert.True(t, r0.Overlaps(Range[int]{5, 5}, cmp))
	assert.True(t, r0.Overlaps(Range[int]{5, 8}, cmp))
	assert.True(t, r0.Overlaps(Range[int]{6, 7}, cmp))
	assert.True(t, r0.Overlaps(Range[int]{8, 10}, cmp))
	assert.False(t, r0.Overlaps(Range[int]{9, 10}, cmp))

	assert.True(t, r0.Overlaps(Range[int]{0, 10}, cmp))
	assert.True(t, r0.Overlaps(Range[int]{4, 9}, cmp))
}

// Inside

func TestRange_Inside(t *testing.T) {
	r0 := Range[int]{
		Start: 3,
		End:   7,
	}

	cmp := compare.Int
	assert.True(t, r0.Inside(Range[int]{3, 7}, cmp))
	assert.True(t, r0.Inside(Range[int]{2, 8}, cmp))
	assert.True(t, r0.Inside(Range[int]{0, 10}, cmp))

	assert.False(t, r0.Inside(Range[int]{0, 0}, cmp))
	assert.False(t, r0.Inside(Range[int]{4, 6}, cmp))
	assert.False(t, r0.Inside(Range[int]{2, 6}, cmp))
	assert.False(t, r0.Inside(Range[int]{4, 8}, cmp))
}

// Expand

func TestRange_Expand__should_expand_range(t *testing.T) {
	r := Range[int]{3, 7}
	cmp := func(a, b int) int { return a - b }

	assert.Equal(t, Range[int]{1, 7}, r.Expand(Range[int]{1, 2}, cmp))
	assert.Equal(t, Range[int]{3, 9}, r.Expand(Range[int]{3, 9}, cmp))
	assert.Equal(t, Range[int]{0, 10}, r.Expand(Range[int]{0, 10}, cmp))
}
