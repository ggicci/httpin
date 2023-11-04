package httpin

import (
	"encoding/base64"
	"strconv"
	"time"
)

func encodeBool(value bool) (string, error) {
	return strconv.FormatBool(value), nil
}

func encodeInt(value int) (string, error) {
	return strconv.FormatInt(int64(value), 10), nil
}

func encodeInt8(value int8) (string, error) {
	return strconv.FormatInt(int64(value), 10), nil
}

func encodeInt16(value int16) (string, error) {
	return strconv.FormatInt(int64(value), 10), nil
}

func encodeInt32(value int32) (string, error) {
	return strconv.FormatInt(int64(value), 10), nil
}

func encodeInt64(value int64) (string, error) {
	return strconv.FormatInt(value, 10), nil
}

func encodeUint(value uint) (string, error) {
	return strconv.FormatUint(uint64(value), 10), nil
}

func encodeUint8(value uint8) (string, error) {
	return strconv.FormatUint(uint64(value), 10), nil
}

func encodeByteSlice(bytes []byte) (string, error) {
	// NOTE: we're using base64.StdEncoding here, not base64.URLEncoding.
	return base64.StdEncoding.EncodeToString(bytes), nil
}

func encodeUint16(value uint16) (string, error) {
	return strconv.FormatUint(uint64(value), 10), nil
}

func encodeUint32(value uint32) (string, error) {
	return strconv.FormatUint(uint64(value), 10), nil
}

func encodeUint64(value uint64) (string, error) {
	return strconv.FormatUint(value, 10), nil
}

func encodeFloat32(value float32) (string, error) {
	return strconv.FormatFloat(float64(value), 'f', -1, 32), nil
}

func encodeFloat64(value float64) (string, error) {
	return strconv.FormatFloat(value, 'f', -1, 64), nil
}

func encodeComplex64(value complex64) (string, error) {
	return strconv.FormatComplex(complex128(value), 'f', -1, 64), nil
}

func encodeComplex128(value complex128) (string, error) {
	return strconv.FormatComplex(value, 'f', -1, 128), nil
}

func encodeString(value string) (string, error) {
	return value, nil
}

// encodeTime encodes a time.Time value to a string in RFC3339Nano format, in UTC timezone.
func encodeTime(value time.Time) (string, error) {
	return value.UTC().Format(time.RFC3339Nano), nil
}
