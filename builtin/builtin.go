package builtin

type IntegerType interface {
	~uint8 | ~uint16 | ~uint32 | ~uint64 | ~int8 | ~int16 | ~int32 | ~int64 | ~int | ~uint | ~uintptr
}

type FloatType interface {
	~float32 | ~float64
}

type NumberType interface {
	IntegerType | FloatType
}

type Ord interface {
	NumberType | string
}

type Ordf[T any] func(a, b T) bool
type Compf[T any] func(a, b T) bool
type Eqf[T any] func(a T) bool

func Equal[T comparable](a, b T) bool { return a == b }
func Less[T Ord](a, b T) bool         { return a < b }
func LessEqual[T Ord](a, b T) bool    { return a <= b }
func Greater[T Ord](a, b T) bool      { return a > b }
func GreaterEqual[T Ord](a, b T) bool { return a >= b }
