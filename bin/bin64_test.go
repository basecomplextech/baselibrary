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

func TestPattern64__should_match_string(t *testing.T) {
	s0 := (Bin64{}).String()
	s1 := Random64().String()
	s2 := " 341a7d60bc5893a6 "
	s3 := "341a7d60-bc5893a6"

	m0 := Pattern64.MatchString(s0)
	m1 := Pattern64.MatchString(s1)
	m2 := Pattern64.MatchString(s2)
	m3 := Pattern64.MatchString(s3)

	assert.True(t, m0)
	assert.True(t, m1)
	assert.False(t, m2)
	assert.False(t, m3)
}
