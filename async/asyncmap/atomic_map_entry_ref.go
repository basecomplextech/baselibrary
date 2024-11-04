// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package asyncmap

// packAtomicMapEntryRef packs an entry id and refcount into a single int64.
func packAtomicMapEntryRef(id int32, refcount int32) int64 {
	return int64(id)<<32 | int64(refcount)
}

// unpackAtomicMapEntryRef unpacks an entry id and refcount from a single int64.
func unpackAtomicMapEntryRef(r int64) (id int32, refcount int32) {
	return int32(r >> 32), int32(r & 0xffffffff)
}
