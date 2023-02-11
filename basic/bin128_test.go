package basic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseBin128String(t *testing.T) {
	v0 := RandomBin128()
	s := v0.String()

	v1, err := ParseBin128String(s)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, v0, v1)
}

func TestBin128Pattern__should_match_byte_string(t *testing.T) {
	s0 := (Bin128{}).String()
	s1 := RandomBin128().String()
	s2 := " 341a7d60bc5893a64bda3de06721534c "
	s3 := "341a7d60bc5893a6-4bda3de06721534c"

	m0 := Bin128Pattern.Match([]byte(s0))
	m1 := Bin128Pattern.Match([]byte(s1))
	m2 := Bin128Pattern.Match([]byte(s2))
	m3 := Bin128Pattern.Match([]byte(s3))

	assert.True(t, m0)
	assert.True(t, m1)
	assert.False(t, m2)
	assert.False(t, m3)
}

func TestBin128Pattern__should_match_string(t *testing.T) {
	s0 := (Bin128{}).String()
	s1 := RandomBin128().String()
	s2 := " 341a7d60bc5893a64bda3de06721534c "
	s3 := "341a7d60bc5893a6-4bda3de06721534c"

	m0 := Bin128Pattern.MatchString(s0)
	m1 := Bin128Pattern.MatchString(s1)
	m2 := Bin128Pattern.MatchString(s2)
	m3 := Bin128Pattern.MatchString(s3)

	assert.True(t, m0)
	assert.True(t, m1)
	assert.False(t, m2)
	assert.False(t, m3)
}
