// Copyright 2026 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package bin

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseString192(t *testing.T) {
	b0 := Random192()
	s := b0.String()

	b1, err := ParseString192(s)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, b0, b1)
}

func TestRegexp192__should_match_byte_string(t *testing.T) {
	s0 := (Bin192{}).String()
	s1 := Random192().String()
	s2 := " 0000000000000000-0000000000000000-0000000000000000 "
	s3 := "0000000000000000-00000000000000000000000000000000"
	s4 := "000000000000000000000000000000000000000000000000"

	m0 := Regexp192.Match([]byte(s0))
	m1 := Regexp192.Match([]byte(s1))
	m2 := Regexp192.Match([]byte(s2))
	m3 := Regexp192.Match([]byte(s3))
	m4 := Regexp192.Match([]byte(s4))

	assert.True(t, m0)
	assert.True(t, m1)
	assert.False(t, m2)
	assert.False(t, m3)
	assert.False(t, m4)
}

func TestRegexp192__should_match_string(t *testing.T) {
	s0 := (Bin192{}).String()
	s1 := Random192().String()
	s2 := " 0000000000000000-0000000000000000-0000000000000000 "
	s3 := "0000000000000000-00000000000000000000000000000000"
	s4 := "000000000000000000000000000000000000000000000000"

	m0 := Regexp192.MatchString(s0)
	m1 := Regexp192.MatchString(s1)
	m2 := Regexp192.MatchString(s2)
	m3 := Regexp192.MatchString(s3)
	m4 := Regexp192.MatchString(s4)

	assert.True(t, m0)
	assert.True(t, m1)
	assert.False(t, m2)
	assert.False(t, m3)
	assert.False(t, m4)
}
