package httpin

import (
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/ggicci/httpin/patch"
	"github.com/stretchr/testify/assert"
)

func removeTypeDecoder[T any]() {
	var zero [0]T
	var fzero [0]patch.Field[T]

	delete(customDecoders, reflect.TypeOf(zero).Elem())
	delete(customDecoders, reflect.TypeOf(fzero).Elem())
}

// decodeCustomBool additionally parses "yes/no".
func decodeCustomBool(value string) (interface{}, error) {
	sdata := strings.ToLower(value)
	if sdata == "yes" {
		return true, nil
	}
	if sdata == "no" {
		return false, nil
	}
	return strconv.ParseBool(sdata)
}

func invalidDecoder(string) error {
	return nil
}

func TestRegisterTypeDecoder(t *testing.T) {
	assert.Panics(t, func() { RegisterTypeDecoder[bool](nil) })
	assert.Panics(t, func() { RegisterTypeDecoder[bool](invalidDecoder) })
	removeTypeDecoder[bool]() // remove the custom decoder

	// Register duplicate decoder should fail.
	assert.NotPanics(t, func() {
		RegisterTypeDecoder[bool](ValueTypeDecoderFunc(decodeCustomBool))
	})
	assert.Panics(t, func() {
		RegisterTypeDecoder[bool](ValueTypeDecoderFunc(decodeCustomBool))
	})
	removeTypeDecoder[bool]() // remove the custom decoder
}

func TestRegisterNamedDecoder(t *testing.T) {
	assert.Panics(t, func() { RegisterNamedDecoder("myBool", nil) })
	assert.Panics(t, func() { RegisterNamedDecoder("myBool", invalidDecoder) })
	removeTypeDecoder[bool]() // remove the custom decoder

	// Register duplicate decoder should fail.
	assert.NotPanics(t, func() {
		RegisterNamedDecoder("mybool", ValueTypeDecoderFunc(decodeCustomBool))
	})
	assert.Panics(t, func() {
		RegisterNamedDecoder("mybool", ValueTypeDecoderFunc(decodeCustomBool))
	})
	removeTypeDecoder[bool]() // remove the custom decoder
}

func TestReplaceTypeDecoder(t *testing.T) {
	assert.NotPanics(t, func() {
		ReplaceTypeDecoder[bool](ValueTypeDecoderFunc(decodeCustomBool))
	})

	assert.NotPanics(t, func() {
		ReplaceTypeDecoder[bool](ValueTypeDecoderFunc(decodeCustomBool))
	})
}

func TestReplaceNamedDecoder(t *testing.T) {
	assert.NotPanics(t, func() {
		ReplaceNamedDecoder("mybool", ValueTypeDecoderFunc(decodeCustomBool))
	})

	assert.NotPanics(t, func() {
		ReplaceNamedDecoder("mybool", ValueTypeDecoderFunc(decodeCustomBool))
	})
}

// Test that the builtin decoders are valid.
func success[T any](t *testing.T, expected T, got interface{}, err error) {
	assert.NoError(t, err)
	_, ok := got.(T)
	assert.True(t, ok)
	assert.Equal(t, expected, got)
}

func fail[T any](t *testing.T, expected T, got interface{}, err error) {
	assert.Error(t, err)
	_, ok := got.(T)
	assert.True(t, ok)
	assert.Equal(t, expected, got)
}

func TestDecoder_bool(t *testing.T) {
	v, err := decodeBool("true")
	success[bool](t, true, v, err)

	v, err = decodeBool("false")
	success[bool](t, false, v, err)

	v, err = decodeBool("1")
	success[bool](t, true, v, err)

	v, err = decodeBool("apple")
	fail[bool](t, false, v, err)
}

func TestDecoder_int(t *testing.T) {
	v, err := decodeInt("2045")
	success[int](t, 2045, v, err)

	v, err = decodeInt("apple")
	fail[int](t, 0, v, err)
}

func TestDecoder_int8(t *testing.T) {
	v, err := decodeInt8("127")
	success[int8](t, 127, v, err)

	v, err = decodeInt8("128")
	fail[int8](t, 127, v, err)

	v, err = decodeInt8("apple")
	fail[int8](t, 0, v, err)
}

func TestDecoder_int16(t *testing.T) {
	v, err := decodeInt16("32767")
	success[int16](t, 32767, v, err)

	v, err = decodeInt16("32768")
	fail[int16](t, 32767, v, err)

	v, err = decodeInt16("apple")
	fail[int16](t, 0, v, err)
}

func TestDecoder_int32(t *testing.T) {
	v, err := decodeInt32("2147483647")
	success[int32](t, 2147483647, v, err)

	v, err = decodeInt32("2147483648")
	fail[int32](t, 2147483647, v, err)

	v, err = decodeInt32("apple")
	fail[int32](t, 0, v, err)
}

func TestDecoder_int64(t *testing.T) {
	v, err := decodeInt64("9223372036854775807")
	success[int64](t, 9223372036854775807, v, err)

	v, err = decodeInt64("9223372036854775808")
	fail[int64](t, 9223372036854775807, v, err)

	v, err = decodeInt64("apple")
	fail[int64](t, 0, v, err)
}

func TestDecoder_uint(t *testing.T) {
	v, err := decodeUint("2045")
	success[uint](t, 2045, v, err)

	v, err = decodeUint("apple")
	fail[uint](t, 0, v, err)
}

func TestDecoder_uint8(t *testing.T) {
	v, err := decodeUint8("255")
	success[uint8](t, 255, v, err)

	v, err = decodeUint8("256")
	fail[uint8](t, 255, v, err)

	v, err = decodeUint8("apple")
	fail[uint8](t, 0, v, err)
}

func TestDecoder_uint16(t *testing.T) {
	v, err := decodeUint16("65535")
	success[uint16](t, 65535, v, err)

	v, err = decodeUint16("65536")
	fail[uint16](t, 65535, v, err)

	v, err = decodeUint16("apple")
	fail[uint16](t, 0, v, err)
}

func TestDecoder_uint32(t *testing.T) {
	v, err := decodeUint32("4294967295")
	success[uint32](t, 4294967295, v, err)

	v, err = decodeUint32("4294967296")
	fail[uint32](t, 4294967295, v, err)

	v, err = decodeUint32("apple")
	fail[uint32](t, 0, v, err)
}

func TestDecoder_uint64(t *testing.T) {
	v, err := decodeUint64("18446744073709551615")
	success[uint64](t, 18446744073709551615, v, err)

	v, err = decodeUint64("18446744073709551616")
	fail[uint64](t, 18446744073709551615, v, err)

	v, err = decodeUint64("apple")
	fail[uint64](t, 0, v, err)
}

func TestDecoder_float32(t *testing.T) {
	v, err := decodeFloat32("3.1415926")
	success[float32](t, 3.1415926, v, err)

	v, err = decodeFloat32("apple")
	fail[float32](t, 0, v, err)
}

func TestDecoder_float64(t *testing.T) {
	v, err := decodeFloat64("3.1415926")
	success[float64](t, 3.1415926, v, err)

	v, err = decodeFloat64("apple")
	fail[float64](t, 0, v, err)
}

func TestDecoder_complex64(t *testing.T) {
	v, err := decodeComplex64("1+4i")
	success[complex64](t, 1+4i, v, err)

	v, err = decodeComplex64("apple")
	fail[complex64](t, 0, v, err)
}

func TestDecoder_complex128(t *testing.T) {
	v, err := decodeComplex128("1+4i")
	success[complex128](t, 1+4i, v, err)

	v, err = decodeComplex128("apple")
	fail[complex128](t, 0, v, err)
}

func TestDecoder_string(t *testing.T) {
	v, err := decodeString("hello")
	success[string](t, "hello", v, err)
}

func TestDecoder_time(t *testing.T) {
	v, err := decodeTime("1991-11-10T08:00:00+08:00")
	assert.NoError(t, err)
	expected := time.Date(1991, 11, 10, 8, 0, 0, 0, time.FixedZone("Asia/Shanghai", +8*3600))
	assert.True(t, equalTime(expected, v.(time.Time)))

	v, err = decodeTime("678088800")
	expected = time.Date(1991, 6, 28, 6, 0, 0, 0, time.UTC)
	assert.NoError(t, err)
	assert.True(t, equalTime(expected, v.(time.Time)))

	v, err = decodeTime("678088800.123456")
	assert.NoError(t, err)
	expected = time.Date(1991, 6, 28, 6, 0, 0, 123456000, time.UTC)
	assert.True(t, equalTime(expected, v.(time.Time)))

	v, err = decodeTime("apple")
	assert.Error(t, err)
	assert.True(t, equalTime(time.Time{}, v.(time.Time)))
}

func equalFuncs(expected, actual interface{}) bool {
	return reflect.ValueOf(expected).Pointer() == reflect.ValueOf(actual).Pointer()
}

func equalTime(expected, actual time.Time) bool {
	return expected.UTC() == actual.UTC()
}
