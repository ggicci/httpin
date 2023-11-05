package internal

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var reUnixtime = regexp.MustCompile(`^\d+(\.\d{1,9})?$`)

func DecodeBool(value string) (bool, error) {
	return strconv.ParseBool(value)
}

func DecodeInt(value string) (int, error) {
	v, err := strconv.ParseInt(value, 10, 64)
	return int(v), err
}

func DecodeInt8(value string) (int8, error) {
	v, err := strconv.ParseInt(value, 10, 8)
	return int8(v), err
}

func DecodeInt16(value string) (int16, error) {
	v, err := strconv.ParseInt(value, 10, 16)
	return int16(v), err
}

func DecodeInt32(value string) (int32, error) {
	v, err := strconv.ParseInt(value, 10, 32)
	return int32(v), err
}

func DecodeInt64(value string) (int64, error) {
	v, err := strconv.ParseInt(value, 10, 64)
	return int64(v), err
}

func DecodeUint(value string) (uint, error) {
	v, err := strconv.ParseUint(value, 10, 64)
	return uint(v), err
}

func DecodeUint8(value string) (uint8, error) {
	v, err := strconv.ParseUint(value, 10, 8)
	return uint8(v), err
}

func DecodeUint16(value string) (uint16, error) {
	v, err := strconv.ParseUint(value, 10, 16)
	return uint16(v), err
}

func DecodeUint32(value string) (uint32, error) {
	v, err := strconv.ParseUint(value, 10, 32)
	return uint32(v), err
}

func DecodeUint64(value string) (uint64, error) {
	v, err := strconv.ParseUint(value, 10, 64)
	return uint64(v), err
}

func DecodeFloat32(value string) (float32, error) {
	v, err := strconv.ParseFloat(value, 32)
	return float32(v), err
}

func DecodeFloat64(value string) (float64, error) {
	v, err := strconv.ParseFloat(value, 64)
	return float64(v), err
}

func DecodeComplex64(value string) (complex64, error) {
	v, err := strconv.ParseComplex(value, 64)
	return complex64(v), err
}

func DecodeComplex128(value string) (complex128, error) {
	v, err := strconv.ParseComplex(value, 128)
	return complex128(v), err
}

func DecodeString(value string) (string, error) {
	return value, nil
}

// DecodeTime parses data bytes as time.Time in UTC timezone.
// Supported formats of the data bytes are:
// 1. RFC3339Nano string, e.g. "2006-01-02T15:04:05-07:00".
// 2. Date string, e.g. "2006-01-02".
// 3. Unix timestamp, e.g. "1136239445", "1136239445.8", "1136239445.812738".
func DecodeTime(value string) (time.Time, error) {
	// Try parsing value as RFC3339 format.
	if t, err := time.ParseInLocation(time.RFC3339Nano, value, time.UTC); err == nil {
		return t.UTC(), nil
	}

	// Try parsing value as date format.
	if t, err := time.ParseInLocation("2006-01-02", value, time.UTC); err == nil {
		return t.UTC(), nil
	}

	// Try parsing value as timestamp, both integer and float formats supported.
	// e.g. "1618974933", "1618974933.284368".
	if reUnixtime.MatchString(value) {
		return DecodeUnixtime(value)
	}

	return time.Time{}, errors.New("invalid time value")
}

// value must be valid unix timestamp, matches reUnixtime.
func DecodeUnixtime(value string) (time.Time, error) {
	parts := strings.Split(value, ".")
	// Note: errors are ignored, since we already validated the value.
	sec, _ := strconv.ParseInt(parts[0], 10, 64)
	var nsec int64
	if len(parts) == 2 {
		nsec, _ = strconv.ParseInt(nanoSecondPrecision(parts[1]), 10, 64)
	}
	return time.Unix(sec, nsec).UTC(), nil
}

func nanoSecondPrecision(value string) string {
	return value + strings.Repeat("0", 9-len(value))
}
