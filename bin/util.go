package bin

import "unsafe"

func unsafeByteString(s string) []byte {
	if s == "" {
		return nil
	}

	d := unsafe.StringData(s)
	return unsafe.Slice(d, len(s))
}
