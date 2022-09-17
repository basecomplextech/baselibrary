package bin

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatch256__should_match_byte_string(t *testing.T) {
	s0 := (Bin256{}).String()
	s1 := Random256().String()
	s2 := " 341a7d60bc5893a64bda3de06721534c341a7d60bc5893a64bda3de06721534c "
	s3 := "341a7d60bc5893a64bda3de06721534c-341a7d60bc5893a64bda3de06721534c"

	m0 := Match256([]byte(s0))
	m1 := Match256([]byte(s1))
	m2 := Match256([]byte(s2))
	m3 := Match256([]byte(s3))

	assert.True(t, m0)
	assert.True(t, m1)
	assert.False(t, m2)
	assert.False(t, m3)
}

func TestMatchString256__should_match_string(t *testing.T) {
	s0 := (Bin256{}).String()
	s1 := Random256().String()
	s2 := " 341a7d60bc5893a64bda3de06721534c341a7d60bc5893a64bda3de06721534c "
	s3 := "341a7d60bc5893a64bda3de06721534c-341a7d60bc5893a64bda3de06721534c"

	m0 := MatchString256(s0)
	m1 := MatchString256(s1)
	m2 := MatchString256(s2)
	m3 := MatchString256(s3)

	assert.True(t, m0)
	assert.True(t, m1)
	assert.False(t, m2)
	assert.False(t, m3)
}
