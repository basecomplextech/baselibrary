package ref

import (
	_ "unsafe"
)

//go:linkname runtimeFastRand runtime.cheaprand
func runtimeFastRand() uint32

// cheaprandn is like cheaprand() % n but faster.
//
//go:linkname runtimeFastRandn runtime.cheaprandn
func runtimeFastRandn(n uint32) uint32
