package basic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatchBin128__should_match_byte_string(t *testing.T) {
	s0 := (Bin128{}).String()
	s1 := RandomBin128().String()
	s2 := " 341a7d60bc5893a64bda3de06721534c "
	s3 := "341a7d60bc5893a6-4bda3de06721534c"

	m0 := MatchBin128([]byte(s0))
	m1 := MatchBin128([]byte(s1))
	m2 := MatchBin128([]byte(s2))
	m3 := MatchBin128([]byte(s3))

	assert.True(t, m0)
	assert.True(t, m1)
	assert.False(t, m2)
	assert.False(t, m3)
}

func TestMatchString128__should_match_string(t *testing.T) {
	s0 := (Bin128{}).String()
	s1 := RandomBin128().String()
	s2 := " 341a7d60bc5893a64bda3de06721534c "
	s3 := "341a7d60bc5893a6-4bda3de06721534c"

	m0 := MatchString128(s0)
	m1 := MatchString128(s1)
	m2 := MatchString128(s2)
	m3 := MatchString128(s3)

	assert.True(t, m0)
	assert.True(t, m1)
	assert.False(t, m2)
	assert.False(t, m3)
}
