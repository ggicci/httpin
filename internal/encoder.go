package internal

import (
	"strconv"
	"time"
)

func EncodeBool(value bool) (string, error) {
	return strconv.FormatBool(value), nil
}

func EncodeInt(value int) (string, error) {
	return strconv.FormatInt(int64(value), 10), nil
}

func EncodeInt8(value int8) (string, error) {
	return strconv.FormatInt(int64(value), 10), nil
}

func EncodeInt16(value int16) (string, error) {
	return strconv.FormatInt(int64(value), 10), nil
}

func EncodeInt32(value int32) (string, error) {
	return strconv.FormatInt(int64(value), 10), nil
}

func EncodeInt64(value int64) (string, error) {
	return strconv.FormatInt(value, 10), nil
}

func EncodeUint(value uint) (string, error) {
	return strconv.FormatUint(uint64(value), 10), nil
}

func EncodeUint8(value uint8) (string, error) {
	return strconv.FormatUint(uint64(value), 10), nil
}

func EncodeUint16(value uint16) (string, error) {
	return strconv.FormatUint(uint64(value), 10), nil
}

func EncodeUint32(value uint32) (string, error) {
	return strconv.FormatUint(uint64(value), 10), nil
}

func EncodeUint64(value uint64) (string, error) {
	return strconv.FormatUint(value, 10), nil
}

func EncodeFloat32(value float32) (string, error) {
	return strconv.FormatFloat(float64(value), 'f', -1, 32), nil
}

func EncodeFloat64(value float64) (string, error) {
	return strconv.FormatFloat(value, 'f', -1, 64), nil
}

func EncodeComplex64(value complex64) (string, error) {
	return strconv.FormatComplex(complex128(value), 'f', -1, 64), nil
}

func EncodeComplex128(value complex128) (string, error) {
	return strconv.FormatComplex(value, 'f', -1, 128), nil
}

func EncodeString(value string) (string, error) {
	return value, nil
}

// encodeTime encodes a time.Time value to a string in RFC3339Nano format, in UTC timezone.
func EncodeTime(value time.Time) (string, error) {
	return value.UTC().Format(time.RFC3339Nano), nil
}
