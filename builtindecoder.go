package httpin

import (
	"fmt"
	"strconv"
	"time"
)

func decodeBool(value string) (interface{}, error) {
	return strconv.ParseBool(value)
}

func decodeInt(value string) (interface{}, error) {
	v, err := strconv.ParseInt(value, 10, 64)
	return int(v), err
}

func decodeInt8(value string) (interface{}, error) {
	v, err := strconv.ParseInt(value, 10, 8)
	return int8(v), err
}

func decodeInt16(value string) (interface{}, error) {
	v, err := strconv.ParseInt(value, 10, 16)
	return int16(v), err
}

func decodeInt32(value string) (interface{}, error) {
	v, err := strconv.ParseInt(value, 10, 32)
	return int32(v), err
}

func decodeInt64(value string) (interface{}, error) {
	v, err := strconv.ParseInt(value, 10, 64)
	return int64(v), err
}

func decodeUint(value string) (interface{}, error) {
	v, err := strconv.ParseUint(value, 10, 64)
	return uint(v), err
}

func decodeUint8(value string) (interface{}, error) {
	v, err := strconv.ParseUint(value, 10, 8)
	return uint8(v), err
}

func decodeUint16(value string) (interface{}, error) {
	v, err := strconv.ParseUint(value, 10, 16)
	return uint16(v), err
}

func decodeUint32(value string) (interface{}, error) {
	v, err := strconv.ParseUint(value, 10, 32)
	return uint32(v), err
}

func decodeUint64(value string) (interface{}, error) {
	v, err := strconv.ParseUint(value, 10, 64)
	return uint64(v), err
}

func decodeFloat32(value string) (interface{}, error) {
	v, err := strconv.ParseFloat(value, 32)
	return float32(v), err
}

func decodeFloat64(value string) (interface{}, error) {
	v, err := strconv.ParseFloat(value, 64)
	return float64(v), err
}

func decodeComplex64(value string) (interface{}, error) {
	v, err := strconv.ParseComplex(value, 64)
	return complex64(v), err
}

func decodeComplex128(value string) (interface{}, error) {
	v, err := strconv.ParseComplex(value, 128)
	return complex128(v), err
}

func decodeString(value string) (interface{}, error) {
	return value, nil
}

// DecodeTime parses data bytes as time.Time in UTC timezone.
// Supported formats of the data bytes are:
// 1. RFC3339Nano string, e.g. "2006-01-02T15:04:05-07:00"
// 2. Unix timestamp, e.g. "1136239445"
func decodeTime(value string) (interface{}, error) {
	// Try parsing value as RFC3339 format.
	if t, err := time.Parse(time.RFC3339Nano, value); err == nil {
		return t.UTC(), nil
	}

	// Try parsing value as timestamp, both integer and float formats supported.
	// e.g. "1618974933", "1618974933.284368".
	if timestamp, err := strconv.ParseInt(value, 10, 64); err == nil {
		return time.Unix(timestamp, 0).UTC(), nil
	}
	if timestamp, err := strconv.ParseFloat(value, 64); err == nil {
		return time.Unix(0, int64(timestamp*float64(time.Second))).UTC(), nil
	}

	return time.Time{}, fmt.Errorf("invalid time value")
}
