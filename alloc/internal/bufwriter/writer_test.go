// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package bufwriter

import (
	"crypto/rand"
	mrand "math/rand/v2"
	"testing"

	"github.com/basecomplextech/baselibrary/buffer"
	"github.com/stretchr/testify/require"
)

func TestWriter__should_buffer_large_writes(t *testing.T) {
	data := make([]byte, 1024*1024)

	_, err := rand.Read(data)
	if err != nil {
		t.Fatal(err)
	}

	dst := buffer.NewSize(len(data))
	w := New(dst)

	b := data
	for len(b) > 0 {
		m := mrand.IntN(32 * 1024)
		m = min(m, len(b))
		p := b[:m]

		n, err := w.Write(p)
		if err != nil {
			t.Fatal(err)
		}

		require.Equal(t, m, n)
		b = b[m:]
	}

	err = w.Flush()
	if err != nil {
		t.Fatal(err)
	}

	data1 := dst.Bytes()
	require.Equal(t, data, data1)
}

func TestWriter__should_buffer_small_writes(t *testing.T) {
	data := make([]byte, 1024*1024)

	_, err := rand.Read(data)
	if err != nil {
		t.Fatal(err)
	}

	dst := buffer.NewSize(len(data))
	w := New(dst)

	b := data
	for len(b) > 0 {
		m := mrand.IntN(1024)
		m = min(m, len(b))
		p := b[:m]

		n, err := w.Write(p)
		if err != nil {
			t.Fatal(err)
		}

		require.Equal(t, m, n)
		b = b[m:]
	}

	err = w.Flush()
	if err != nil {
		t.Fatal(err)
	}

	data1 := dst.Bytes()
	require.Equal(t, data, data1)
}
