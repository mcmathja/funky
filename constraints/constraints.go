// constraints provides a set of useful type sets over core types.
// Some of these duplicate those specified in the standard library's
// "experimental" package in order to avoid external dependencies.
package constraints

type Complex interface {
	~complex64 | ~complex128
}

type Float interface {
	~float32 | ~float64
}

type Integer interface {
	Signed | Unsigned
}

type Numeric interface {
	Integer | Float | Complex
}

type Ordered interface {
	Integer | Float | ~string
}

type Real interface {
	Integer | Float
}

type Signed interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

type Unsigned interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}
