package internal

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEncoder_bool(t *testing.T) {
	for _, tc := range []struct {
		value    bool
		expected string
	}{
		{true, "true"},
		{false, "false"},
	} {
		actual, err := EncodeBool(tc.value)
		assert.NoError(t, err)
		assert.Equal(t, tc.expected, actual)
	}
}

func TestEncoder_int(t *testing.T) {
	for _, tc := range []struct {
		value    int
		expected string
	}{
		{0, "0"},
		{2045, "2045"},
		{-2045, "-2045"},
	} {
		actual, err := EncodeInt(tc.value)
		assert.NoError(t, err)
		assert.Equal(t, tc.expected, actual)
	}
}

func TestEncoder_int8(t *testing.T) {
	for _, tc := range []struct {
		value    int8
		expected string
	}{
		{0, "0"},
		{127, "127"},
		{-127, "-127"},
	} {
		actual, err := EncodeInt8(tc.value)
		assert.NoError(t, err)
		assert.Equal(t, tc.expected, actual)
	}
}

func TestEncoder_int16(t *testing.T) {
	for _, tc := range []struct {
		value    int16
		expected string
	}{
		{0, "0"},
		{32767, "32767"},
		{-32767, "-32767"},
	} {
		actual, err := EncodeInt16(tc.value)
		assert.NoError(t, err)
		assert.Equal(t, tc.expected, actual)
	}
}

func TestEncoder_int32(t *testing.T) {
	for _, tc := range []struct {
		value    int32
		expected string
	}{
		{0, "0"},
		{2147483647, "2147483647"},
		{-2147483647, "-2147483647"},
	} {
		actual, err := EncodeInt32(tc.value)
		assert.NoError(t, err)
		assert.Equal(t, tc.expected, actual)
	}
}

func TestEncoder_int64(t *testing.T) {
	for _, tc := range []struct {
		value    int64
		expected string
	}{
		{0, "0"},
		{9223372036854775807, "9223372036854775807"},
		{-9223372036854775807, "-9223372036854775807"},
	} {
		actual, err := EncodeInt64(tc.value)
		assert.NoError(t, err)
		assert.Equal(t, tc.expected, actual)
	}
}

func TestEncoder_uint(t *testing.T) {
	for _, tc := range []struct {
		value    uint
		expected string
	}{
		{0, "0"},
		{2045, "2045"},
	} {
		actual, err := EncodeUint(tc.value)
		assert.NoError(t, err)
		assert.Equal(t, tc.expected, actual)
	}
}

func TestEncoder_uint8(t *testing.T) {
	for _, tc := range []struct {
		value    uint8
		expected string
	}{
		{0, "0"},
		{127, "127"},
	} {
		actual, err := EncodeUint8(tc.value)
		assert.NoError(t, err)
		assert.Equal(t, tc.expected, actual)
	}
}

func TestEncoder_uint16(t *testing.T) {
	for _, tc := range []struct {
		value    uint16
		expected string
	}{
		{0, "0"},
		{32767, "32767"},
	} {
		actual, err := EncodeUint16(tc.value)
		assert.NoError(t, err)
		assert.Equal(t, tc.expected, actual)
	}
}

func TestEncoder_uint32(t *testing.T) {
	for _, tc := range []struct {
		value    uint32
		expected string
	}{
		{0, "0"},
		{2147483647, "2147483647"},
	} {
		actual, err := EncodeUint32(tc.value)
		assert.NoError(t, err)
		assert.Equal(t, tc.expected, actual)
	}
}

func TestEncoder_uint64(t *testing.T) {
	for _, tc := range []struct {
		value    uint64
		expected string
	}{
		{0, "0"},
		{9223372036854775807, "9223372036854775807"},
	} {
		actual, err := EncodeUint64(tc.value)
		assert.NoError(t, err)
		assert.Equal(t, tc.expected, actual)
	}
}

func TestEncoder_float32(t *testing.T) {
	for _, tc := range []struct {
		value    float32
		expected string
	}{
		{0, "0"},
		{1.23456789, "1.2345679"},
		{-1.23456789, "-1.2345679"},
	} {
		actual, err := EncodeFloat32(tc.value)
		assert.NoError(t, err)
		assert.Equal(t, tc.expected, actual)
	}
}

func TestEncoder_float64(t *testing.T) {
	for _, tc := range []struct {
		value    float64
		expected string
	}{
		{0, "0"},
		{1.23456789, "1.23456789"},
		{-1.23456789, "-1.23456789"},
	} {
		actual, err := EncodeFloat64(tc.value)
		assert.NoError(t, err)
		assert.Equal(t, tc.expected, actual)
	}
}

func TestEncoder_complex64(t *testing.T) {
	for _, tc := range []struct {
		value    complex64
		expected string
	}{
		{0, "(0+0i)"},
		{1 + 4i, "(1+4i)"},
		{-1 - 4i, "(-1-4i)"},
	} {
		actual, err := EncodeComplex64(tc.value)
		assert.NoError(t, err)
		assert.Equal(t, tc.expected, actual)
	}
}

func TestEncoder_complex128(t *testing.T) {
	for _, tc := range []struct {
		value    complex128
		expected string
	}{
		{0, "(0+0i)"},
		{1 + 4i, "(1+4i)"},
		{-1 - 4i, "(-1-4i)"},
	} {
		actual, err := EncodeComplex128(tc.value)
		assert.NoError(t, err)
		assert.Equal(t, tc.expected, actual)
	}
}

func TestEncoder_string(t *testing.T) {
	for _, tc := range []struct {
		value    string
		expected string
	}{
		{"", ""},
		{"hello", "hello"},
	} {
		actual, err := EncodeString(tc.value)
		assert.NoError(t, err)
		assert.Equal(t, tc.expected, actual)
	}
}

func TestEncoder_time(t *testing.T) {
	testcases := []struct {
		value    time.Time
		expected string
	}{
		{time.Date(2006, 1, 2, 0, 0, 0, 0, time.UTC), "2006-01-02T00:00:00Z"},
		{time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC), "2006-01-02T15:04:05Z"},
		{time.Date(2006, 1, 2, 22, 4, 5, 999999999, time.UTC), "2006-01-02T22:04:05.999999999Z"},
		{time.Date(2006, 01, 02, 15, 4, 5, 0, time.FixedZone("UTC+7", 7*60*60)), "2006-01-02T08:04:05Z"},
	}

	for _, tc := range testcases {
		actual, err := EncodeTime(tc.value)
		assert.NoError(t, err)
		assert.Equal(t, tc.expected, actual)
	}
}
