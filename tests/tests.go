// Copyright 2021 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package tests

import (
	"crypto/rand"
	"io/ioutil"
	"os"

	mrand "math/rand"
)

type T interface {
	Cleanup(f func())
	Error(...any)
	Fatal(...any)
	Fatalf(format string, args ...any)
	Helper()
}

func TestDir(t T) string {
	name, err := ioutil.TempDir("", "tests-")
	if err != nil {
		t.Fatal(err)
	}
	return name
}

func TestFile(t T) *os.File {
	file, err := ioutil.TempFile("", "tests-")
	if err != nil {
		t.Fatal(err)
	}
	return file
}

func TestBytes(t T, size int) []byte {
	b := make([]byte, size)
	if _, err := rand.Read(b); err != nil {
		t.Fatal(err)
	}
	return b
}

func TestBytesN(t T, size int, n int) [][]byte {
	result := make([][]byte, 0, n)
	for i := 0; i < n; i++ {
		data := TestBytes(t, size)
		result = append(result, data)
	}
	return result
}

func TestRandomSize(t T, maxSize int) []byte {
	size := mrand.Intn(maxSize)
	if size == 0 {
		size = maxSize
	}
	return TestBytes(t, size)
}

func TestRandomSizeN(t T, maxSize int, n int) [][]byte {
	result := make([][]byte, 0, n)
	for i := 0; i < n; i++ {
		data := TestRandomSize(t, maxSize)
		result = append(result, data)
	}
	return result
}
