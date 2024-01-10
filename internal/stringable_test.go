package internal

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/ggicci/owl"
)

func TestNewStringable_string(t *testing.T) {
	var s string = "hello"
	rvString := reflect.ValueOf(s)
	assert.Panics(t, func() {
		NewStringable(rvString)
	})

	rvStringPointer := reflect.ValueOf(&s)
	sv, err := NewStringable(rvStringPointer)
	assert.NoError(t, err)
	got, err := sv.ToString()
	assert.NoError(t, err)
	assert.Equal(t, "hello", got)
	sv.FromString("world")
	assert.Equal(t, "world", s)
}

func TestNewStringable_bool(t *testing.T) {
	var b bool = true
	rvBool := reflect.ValueOf(b)
	assert.Panics(t, func() {
		NewStringable(rvBool)
	})

	rvBoolPointer := reflect.ValueOf(&b)
	sv, err := NewStringable(rvBoolPointer)
	assert.NoError(t, err)
	got, err := sv.ToString()
	assert.NoError(t, err)
	assert.Equal(t, "true", got)
	sv.FromString("false")
	assert.Equal(t, false, b)

	assert.Error(t, sv.FromString("hello"))
}

func TestNewStringable_int(t *testing.T) {
	testInteger[int](t, 2045, "hello")
}

func TestNewStringable_int8(t *testing.T) {
	testInteger[int8](t, int8(127), "128")
}

func TestNewStringable_int16(t *testing.T) {
	testInteger[int16](t, int16(32767), "32768")
}

func TestNewStringable_int32(t *testing.T) {
	testInteger[int32](t, int32(2147483647), "2147483648")
}

func TestNewStringable_int64(t *testing.T) {
	testInteger[int64](t, int64(9223372036854775807), "9223372036854775808")
}

func TestNewStringable_uint(t *testing.T) {
	testInteger[uint](t, uint(2045), "-1")
}

func TestNewStringable_uint8(t *testing.T) {
	testInteger[uint8](t, uint8(255), "256")
}

func TestNewStringable_uint16(t *testing.T) {
	testInteger[uint16](t, uint16(65535), "65536")
}

func TestNewStringable_uint32(t *testing.T) {
	testInteger[uint32](t, uint32(4294967295), "4294967296")
}

func TestNewStringable_uint64(t *testing.T) {
	testInteger[uint64](t, uint64(18446744073709551615), "18446744073709551616")
}

func TestNewStringable_float32(t *testing.T) {
	testInteger[float32](t, float32(3.1415926), "hello")
}

func TestNewStringable_float64(t *testing.T) {
	testInteger[float64](t, float64(3.14159265358979323846264338327950288419716939937510582097494459), "hello")
}

func TestNewStringable_complex64(t *testing.T) {
	testInteger[complex64](t, complex64(3.1415926+2.71828i), "hello")
}

func TestNewStringable_complex128(t *testing.T) {
	testInteger[complex128](t, complex128(3.14159265358979323846264338327950288419716939937510582097494459+2.71828182845904523536028747135266249775724709369995957496696763i), "hello")
}

func TestNewStringable_Time(t *testing.T) {
	var now = time.Now()
	rvTime := reflect.ValueOf(now)
	assert.Panics(t, func() {
		NewStringable(rvTime)
	})

	rvTimePointer := reflect.ValueOf(&now)
	sv, err := NewStringable(rvTimePointer)
	assert.NoError(t, err)

	// RFC3339Nano
	testTime(t, sv, "1991-11-10T08:00:00+08:00", time.Date(1991, 11, 10, 8, 0, 0, 0, time.FixedZone("Asia/Shanghai", +8*3600)), "1991-11-10T00:00:00Z")
	// Date string
	testTime(t, sv, "1991-11-10", time.Date(1991, 11, 10, 0, 0, 0, 0, time.UTC), "1991-11-10T00:00:00Z")

	// Unix timestamp
	testTime(t, sv, "678088800", time.Date(1991, 6, 28, 6, 0, 0, 0, time.UTC), "1991-06-28T06:00:00Z")

	// Unix timestamp fraction
	testTime(t, sv, "678088800.123456789", time.Date(1991, 6, 28, 6, 0, 0, 123456789, time.UTC), "1991-06-28T06:00:00.123456789Z")

	// Unsupported format
	assert.Error(t, sv.FromString("hello"))
}

func TestNewStringable_ByteSlice(t *testing.T) {
	var b []byte = []byte("hello")
	rvByteSlice := reflect.ValueOf(b)
	assert.NotPanics(t, func() {
		NewStringable(rvByteSlice)
	})

	rvByteSlicePointer := reflect.ValueOf(&b)
	sv, err := NewStringable(rvByteSlicePointer)
	assert.NoError(t, err)
	got, err := sv.ToString()
	assert.NoError(t, err)
	assert.Equal(t, "aGVsbG8=", got)

	sv.FromString("d29ybGQ=")
	assert.Equal(t, []byte("world"), b)

	assert.Error(t, sv.FromString("hello"))
}

func TestNewStringable_ErrUnsupportedType(t *testing.T) {
	type MyStruct struct{ Name string }
	var s MyStruct
	rvStruct := reflect.ValueOf(s)
	assert.Panics(t, func() {
		NewStringable(rvStruct)
	})
	rvStructPointer := reflect.ValueOf(&s)
	sv, err := NewStringable(rvStructPointer)
	assert.ErrorIs(t, err, owl.ErrUnsupportedType)
	assert.Nil(t, sv)
}

type Numeric interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64 | complex64 | complex128
}

func testInteger[T Numeric](t *testing.T, vSuccess T, invalidStr string) {
	rv := reflect.ValueOf(vSuccess)
	assert.Panics(t, func() {
		NewStringable(rv)
	})

	rvPointer := reflect.ValueOf(&vSuccess)
	sv, err := NewStringable(rvPointer)
	assert.NoError(t, err)
	got, err := sv.ToString()
	assert.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("%v", vSuccess), got)
	sv.FromString("2")
	assert.Equal(t, T(2), vSuccess)

	assert.Error(t, sv.FromString(invalidStr))
}

func testTime(t *testing.T, sv Stringable, fromStr string, expected time.Time, expectedToStr string) {
	assert.NoError(t, sv.FromString(fromStr))
	assert.True(t, equalTime(expected, time.Time(*sv.(*Time))))
	ts, err := sv.ToString()
	assert.NoError(t, err)
	assert.Equal(t, expectedToStr, ts)
}

func equalTime(expected, actual time.Time) bool {
	return expected.UTC() == actual.UTC()
}
