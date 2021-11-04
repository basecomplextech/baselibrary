package tests

import (
	"crypto/rand"
	"io/ioutil"
	"os"
)

type T interface {
	Fatal(args ...interface{})
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
