package bin128

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatch__should_match_byte_string(t *testing.T) {
	s0 := (B128{}).String()
	s1 := Random().String()
	s2 := " 341a7d60bc5893a64bda3de06721534c "
	s3 := "341a7d60bc5893a6-4bda3de06721534c"

	m0 := Match([]byte(s0))
	m1 := Match([]byte(s1))
	m2 := Match([]byte(s2))
	m3 := Match([]byte(s3))

	assert.True(t, m0)
	assert.True(t, m1)
	assert.False(t, m2)
	assert.False(t, m3)
}

func TestMatchString__should_match_string(t *testing.T) {
	s0 := (B128{}).String()
	s1 := Random().String()
	s2 := " 341a7d60bc5893a64bda3de06721534c "
	s3 := "341a7d60bc5893a6-4bda3de06721534c"

	m0 := MatchString(s0)
	m1 := MatchString(s1)
	m2 := MatchString(s2)
	m3 := MatchString(s3)

	assert.True(t, m0)
	assert.True(t, m1)
	assert.False(t, m2)
	assert.False(t, m3)
}
