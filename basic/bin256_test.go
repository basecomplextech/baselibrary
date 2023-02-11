package basic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseBin256String(t *testing.T) {
	u0 := RandomBin256()
	s := u0.String()

	u1, err := ParseBin256String(s)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, u0, u1)
}

func TestBin256Pattern__should_match_byte_string(t *testing.T) {
	s0 := (Bin256{}).String()
	s1 := RandomBin256().String()
	s2 := " 341a7d60bc5893a64bda3de06721534c341a7d60bc5893a64bda3de06721534c "
	s3 := "341a7d60bc5893a64bda3de06721534c-341a7d60bc5893a64bda3de06721534c"

	m0 := Bin256Pattern.Match([]byte(s0))
	m1 := Bin256Pattern.Match([]byte(s1))
	m2 := Bin256Pattern.Match([]byte(s2))
	m3 := Bin256Pattern.Match([]byte(s3))

	assert.True(t, m0)
	assert.True(t, m1)
	assert.False(t, m2)
	assert.False(t, m3)
}

func TestBin256Pattern__should_match_string(t *testing.T) {
	s0 := (Bin256{}).String()
	s1 := RandomBin256().String()
	s2 := " 341a7d60bc5893a64bda3de06721534c341a7d60bc5893a64bda3de06721534c "
	s3 := "341a7d60bc5893a64bda3de06721534c-341a7d60bc5893a64bda3de06721534c"

	m0 := Bin256Pattern.MatchString(s0)
	m1 := Bin256Pattern.MatchString(s1)
	m2 := Bin256Pattern.MatchString(s2)
	m3 := Bin256Pattern.MatchString(s3)

	assert.True(t, m0)
	assert.True(t, m1)
	assert.False(t, m2)
	assert.False(t, m3)
}
