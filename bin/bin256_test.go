package bin

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseString256(t *testing.T) {
	b0 := Random256()
	s := b0.String()

	b1, err := ParseString256(s)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, b0, b1)
}

func TestPattern256__should_match_byte_string(t *testing.T) {
	s0 := (Bin256{}).String()
	s1 := Random256().String()
	s2 := " 341a7d60bc5893a64bda3de06721534c341a7d60bc5893a64bda3de06721534c "
	s3 := "341a7d60bc5893a64bda3de06721534c-341a7d60bc5893a64bda3de06721534c"
	s4 := "341a7d60bc5893a64bda3de06721534c341a7d60bc5893a64bda3de06721534c"

	m0 := Pattern256.Match([]byte(s0))
	m1 := Pattern256.Match([]byte(s1))
	m2 := Pattern256.Match([]byte(s2))
	m3 := Pattern256.Match([]byte(s3))
	m4 := Pattern256.Match([]byte(s4))

	assert.True(t, m0)
	assert.True(t, m1)
	assert.False(t, m2)
	assert.False(t, m3)
	assert.False(t, m4)
}

func TestPattern256__should_match_string(t *testing.T) {
	s0 := (Bin256{}).String()
	s1 := Random256().String()
	s2 := " 341a7d60bc5893a64bda3de06721534c341a7d60bc5893a64bda3de06721534c "
	s3 := "341a7d60bc5893a64bda3de06721534c-341a7d60bc5893a64bda3de06721534c"
	s4 := "341a7d60bc5893a64bda3de06721534c341a7d60bc5893a64bda3de06721534c"

	m0 := Pattern256.MatchString(s0)
	m1 := Pattern256.MatchString(s1)
	m2 := Pattern256.MatchString(s2)
	m3 := Pattern256.MatchString(s3)
	m4 := Pattern256.MatchString(s4)

	assert.True(t, m0)
	assert.True(t, m1)
	assert.False(t, m2)
	assert.False(t, m3)
	assert.False(t, m4)
}
