package constraints

// Integer specifies all integers.
type Integer interface {
	Signed | Unsigned
}

// Signed specifies signed integers.
type Signed interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

// Unsigned specifies unsigned integers.
type Unsigned interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

// Float specifies all float-point types.
type Float interface {
	~float32 | ~float64
}

// Number specifies an integer or a float.
type Number interface {
	Integer | Float
}

// Complex specifies all complex types.
type Complex interface {
	~complex64 | ~complex128
}

// Ordered specifies types which support comparision.
type Ordered interface {
	Integer | Float | ~string
}
