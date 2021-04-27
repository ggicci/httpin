package httpin

import "time"

type Bool struct {
	Value bool
	Valid bool // Valid is true if the corresponding key were found in the source
}

type Int struct {
	Value int
	Valid bool
}

type Int8 struct {
	Value int8
	Valid bool
}

type Int16 struct {
	Value int16
	Valid bool
}

type Int32 struct {
	Value int32
	Valid bool
}

type Int64 struct {
	Value int64
	Valid bool
}

type Uint struct {
	Value uint
	Valid bool
}

type Uint8 struct {
	Value uint8
	Valid bool
}

type Uint16 struct {
	Value uint16
	Valid bool
}

type Uint32 struct {
	Value uint32
	Valid bool
}

type Uint64 struct {
	Value uint64
	Valid bool
}

type Float32 struct {
	Value float32
	Valid bool
}

type Float64 struct {
	Value float64
	Valid bool
}

type Complex64 struct {
	Value complex64
	Valid bool
}

type Complex128 struct {
	Value complex128
	Valid bool
}

type String struct {
	Value string
	Valid bool
}

type Time struct {
	Value time.Time
	Valid bool
}

type BoolArray struct {
	Value []bool
	Valid bool
}

type IntArray struct {
	Value []int
	Valid bool
}

type Int8Array struct {
	Value []int8
	Valid bool
}

type Int16Array struct {
	Value []int16
	Valid bool
}

type Int32Array struct {
	Value []int32
	Valid bool
}

type Int64Array struct {
	Value []int64
	Valid bool
}

type UintArray struct {
	Value []uint
	Valid bool
}

type Uint8Array struct {
	Value []uint8
	Valid bool
}

type Uint16Array struct {
	Value []uint16
	Valid bool
}

type Uint32Array struct {
	Value []uint32
	Valid bool
}

type Uint64Array struct {
	Value []uint64
	Valid bool
}

type Float32Array struct {
	Value []float32
	Valid bool
}

type Float64Array struct {
	Value []float64
	Valid bool
}

type Complex64Array struct {
	Value []complex64
	Valid bool
}

type Complex128Array struct {
	Value []complex128
	Valid bool
}

type StringArray struct {
	Value []string
	Valid bool
}

type TimeArray struct {
	Value []time.Time
	Valid bool
}
