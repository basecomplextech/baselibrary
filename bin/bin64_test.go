// Copyright 2023 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package bin

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTimeRandom64(t *testing.T) {
	b0 := TimeRandom64().String()
	b1 := TimeRandom64().String()
	assert.NotEqual(t, b0, b1)
}

func TestParseString64(t *testing.T) {
	b0 := Random64()
	s := b0.String()

	b1, err := ParseString64(s)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, b0, b1)
}

func TestRegexp64__should_match_string(t *testing.T) {
	s0 := (Bin64{}).String()
	s1 := Random64().String()
	s2 := " 341a7d60bc5893a6 "
	s3 := "341a7d60-bc5893a6"

	m0 := Regexp64.MatchString(s0)
	m1 := Regexp64.MatchString(s1)
	m2 := Regexp64.MatchString(s2)
	m3 := Regexp64.MatchString(s3)

	assert.True(t, m0)
	assert.True(t, m1)
	assert.False(t, m2)
	assert.False(t, m3)
}
