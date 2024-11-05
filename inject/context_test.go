// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package inject

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContext_Get__should_build_object_and_set_pointer(t *testing.T) {
	type object struct {
		a int
		b string
		c bool
	}
	var obj object

	Test(t,
		1,
		"hello",
		func() bool { return true },
		func(a int, b string, c bool) object {
			return object{a, b, c}
		},
	).Get(&obj)

	assert.Equal(t, 1, obj.a)
	assert.Equal(t, "hello", obj.b)
	assert.Equal(t, true, obj.c)
}

func TestContext_Get__should_panic_on_cycle(t *testing.T) {
	assert.PanicsWithValue(t, "cycle detected: int -> int32 -> string -> int", func() {
		var obj int

		New(
			func(b int32) int { return int(b) },
			func(c string) int32 { return 123 },
			func(a int) string { return "123" },
		).Get(&obj)
	})
}

func TestContext_Add__should_panic_on_duplicate_provider(t *testing.T) {
	assert.PanicsWithValue(t, "duplicate provider: int", func() {
		New().
			Add(1).
			Add(1)
	})
}
