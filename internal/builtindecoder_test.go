package internal

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDecoder_bool(t *testing.T) {
	v, err := DecodeBool("true")
	success[bool](t, true, v, err)

	v, err = DecodeBool("false")
	success[bool](t, false, v, err)

	v, err = DecodeBool("1")
	success[bool](t, true, v, err)

	v, err = DecodeBool("apple")
	fail[bool](t, false, v, err)
}

func TestDecoder_int(t *testing.T) {
	v, err := DecodeInt("2045")
	success[int](t, 2045, v, err)

	v, err = DecodeInt("apple")
	fail[int](t, 0, v, err)
}

func TestDecoder_int8(t *testing.T) {
	v, err := DecodeInt8("127")
	success[int8](t, 127, v, err)

	v, err = DecodeInt8("128")
	fail[int8](t, 127, v, err)

	v, err = DecodeInt8("apple")
	fail[int8](t, 0, v, err)
}

func TestDecoder_int16(t *testing.T) {
	v, err := DecodeInt16("32767")
	success[int16](t, 32767, v, err)

	v, err = DecodeInt16("32768")
	fail[int16](t, 32767, v, err)

	v, err = DecodeInt16("apple")
	fail[int16](t, 0, v, err)
}

func TestDecoder_int32(t *testing.T) {
	v, err := DecodeInt32("2147483647")
	success[int32](t, 2147483647, v, err)

	v, err = DecodeInt32("2147483648")
	fail[int32](t, 2147483647, v, err)

	v, err = DecodeInt32("apple")
	fail[int32](t, 0, v, err)
}

func TestDecoder_int64(t *testing.T) {
	v, err := DecodeInt64("9223372036854775807")
	success[int64](t, 9223372036854775807, v, err)

	v, err = DecodeInt64("9223372036854775808")
	fail[int64](t, 9223372036854775807, v, err)

	v, err = DecodeInt64("apple")
	fail[int64](t, 0, v, err)
}

func TestDecoder_uint(t *testing.T) {
	v, err := DecodeUint("2045")
	success[uint](t, 2045, v, err)

	v, err = DecodeUint("apple")
	fail[uint](t, 0, v, err)
}

func TestDecoder_uint8(t *testing.T) {
	v, err := DecodeUint8("255")
	success[uint8](t, 255, v, err)

	v, err = DecodeUint8("256")
	fail[uint8](t, 255, v, err)

	v, err = DecodeUint8("apple")
	fail[uint8](t, 0, v, err)
}

func TestDecoder_uint16(t *testing.T) {
	v, err := DecodeUint16("65535")
	success[uint16](t, 65535, v, err)

	v, err = DecodeUint16("65536")
	fail[uint16](t, 65535, v, err)

	v, err = DecodeUint16("apple")
	fail[uint16](t, 0, v, err)
}

func TestDecoder_uint32(t *testing.T) {
	v, err := DecodeUint32("4294967295")
	success[uint32](t, 4294967295, v, err)

	v, err = DecodeUint32("4294967296")
	fail[uint32](t, 4294967295, v, err)

	v, err = DecodeUint32("apple")
	fail[uint32](t, 0, v, err)
}

func TestDecoder_uint64(t *testing.T) {
	v, err := DecodeUint64("18446744073709551615")
	success[uint64](t, 18446744073709551615, v, err)

	v, err = DecodeUint64("18446744073709551616")
	fail[uint64](t, 18446744073709551615, v, err)

	v, err = DecodeUint64("apple")
	fail[uint64](t, 0, v, err)
}

func TestDecoder_float32(t *testing.T) {
	v, err := DecodeFloat32("3.1415926")
	success[float32](t, 3.1415926, v, err)

	v, err = DecodeFloat32("apple")
	fail[float32](t, 0, v, err)
}

func TestDecoder_float64(t *testing.T) {
	v, err := DecodeFloat64("3.1415926")
	success[float64](t, 3.1415926, v, err)

	v, err = DecodeFloat64("apple")
	fail[float64](t, 0, v, err)
}

func TestDecoder_complex64(t *testing.T) {
	v, err := DecodeComplex64("1+4i")
	success[complex64](t, 1+4i, v, err)

	v, err = DecodeComplex64("apple")
	fail[complex64](t, 0, v, err)
}

func TestDecoder_complex128(t *testing.T) {
	v, err := DecodeComplex128("1+4i")
	success[complex128](t, 1+4i, v, err)

	v, err = DecodeComplex128("apple")
	fail[complex128](t, 0, v, err)
}

func TestDecoder_string(t *testing.T) {
	v, err := DecodeString("hello")
	success[string](t, "hello", v, err)
}

func TestDecoder_time(t *testing.T) {
	testcases := []struct {
		value    string
		expected time.Time
	}{
		{"2006-01-02", time.Date(2006, 1, 2, 0, 0, 0, 0, time.UTC)},
		{"2006-01-02T15:04:05Z", time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC)},
		{"2006-01-02T15:04:05-07:00", time.Date(2006, 1, 2, 22, 4, 5, 0, time.UTC)},
		{"1991-11-10T08:00:00+08:00", time.Date(1991, 11, 10, 8, 0, 0, 0, time.FixedZone("Asia/Shanghai", +8*3600))},
		{"2006-01-02T15:04:05.999999999-07:00", time.Date(2006, 1, 2, 15, 4, 5, 999999999, time.FixedZone("UTC-7", -7*3600))},
		{"0", time.Unix(0, 0).UTC()},
		{"1136239445", time.Unix(1136239445, 0).UTC()},
		{"1136239445.0", time.Unix(1136239445, 0).UTC()},
		{"1136239445.8", time.Unix(1136239445, 800000000).UTC()},
		{"1136239445.812738", time.Unix(1136239445, 812738000).UTC()},
		{"1136239445.123456789", time.Unix(1136239445, 123456789).UTC()},
	}

	for _, tc := range testcases {
		actual, err := DecodeTime(tc.value)
		assert.IsType(t, time.Time{}, actual)
		assert.NoError(t, err)
		_, zoneOffset := actual.Zone()
		assert.Equal(t, 0, zoneOffset) // always UTC
		equalTime(tc.expected, actual)
	}

	errorCases := []string{
		"2006-01-02T15:04:05",   // missing timezone
		"1136239445.",           // no fractional part
		"1136239445.1234567890", // nanosecond precision is limited to 9 digits
		"Jan 2 15:04:05 2006",   // not supported format
		"apple",                 // not supported format
	}
	for _, tc := range errorCases {
		actual, err := DecodeTime(tc)
		assert.Equal(t, time.Time{}, actual)
		assert.ErrorContains(t, err, "invalid time value")
	}
}

func equalTime(expected, actual time.Time) bool {
	return expected.UTC() == actual.UTC()
}

func success[T any](t *testing.T, expected T, got any, err error) {
	assert.NoError(t, err)
	_, ok := got.(T)
	assert.True(t, ok)
	assert.Equal(t, expected, got)
}

func fail[T any](t *testing.T, expected T, got any, err error) {
	assert.Error(t, err)
	_, ok := got.(T)
	assert.True(t, ok)
	assert.Equal(t, expected, got)
}
