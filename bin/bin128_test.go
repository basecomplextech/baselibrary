// Copyright 2023 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package bin

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTimeRandom128(t *testing.T) {
	b0 := TimeRandom128().String()
	b1 := TimeRandom128().String()
	assert.NotEqual(t, b0, b1)
}

func TestParseString128(t *testing.T) {
	b0 := Random128()
	s := b0.String()

	b1, err := ParseString128(s)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, b0, b1)
}

func TestRegexp128__should_match_byte_string(t *testing.T) {
	s0 := (Bin128{}).String()
	s1 := Random128().String()
	s2 := " 341a7d60bc5893a64bda3de06721534c "
	s3 := "341a7d60bc5893a64bda3de06721534c"

	m0 := Regexp128.Match([]byte(s0))
	m1 := Regexp128.Match([]byte(s1))
	m2 := Regexp128.Match([]byte(s2))
	m3 := Regexp128.Match([]byte(s3))

	assert.True(t, m0)
	assert.True(t, m1)
	assert.False(t, m2)
	assert.False(t, m3)
}

func TestRegexp128__should_match_string(t *testing.T) {
	s0 := (Bin128{}).String()
	s1 := Random128().String()
	s2 := " 341a7d60bc5893a64bda3de06721534c "
	s3 := "341a7d60bc5893a64bda3de06721534c"

	m0 := Regexp128.MatchString(s0)
	m1 := Regexp128.MatchString(s1)
	m2 := Regexp128.MatchString(s2)
	m3 := Regexp128.MatchString(s3)

	assert.True(t, m0)
	assert.True(t, m1)
	assert.False(t, m2)
	assert.False(t, m3)
}
