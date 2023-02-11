package basic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatchBin256__should_match_byte_string(t *testing.T) {
	s0 := (Bin256{}).String()
	s1 := RandomBin256().String()
	s2 := " 341a7d60bc5893a64bda3de06721534c341a7d60bc5893a64bda3de06721534c "
	s3 := "341a7d60bc5893a64bda3de06721534c-341a7d60bc5893a64bda3de06721534c"

	m0 := MatchBin256([]byte(s0))
	m1 := MatchBin256([]byte(s1))
	m2 := MatchBin256([]byte(s2))
	m3 := MatchBin256([]byte(s3))

	assert.True(t, m0)
	assert.True(t, m1)
	assert.False(t, m2)
	assert.False(t, m3)
}

func TestMatchStringBin256__should_match_string(t *testing.T) {
	s0 := (Bin256{}).String()
	s1 := RandomBin256().String()
	s2 := " 341a7d60bc5893a64bda3de06721534c341a7d60bc5893a64bda3de06721534c "
	s3 := "341a7d60bc5893a64bda3de06721534c-341a7d60bc5893a64bda3de06721534c"

	m0 := MatchStringBin256(s0)
	m1 := MatchStringBin256(s1)
	m2 := MatchStringBin256(s2)
	m3 := MatchStringBin256(s3)

	assert.True(t, m0)
	assert.True(t, m1)
	assert.False(t, m2)
	assert.False(t, m3)
}
