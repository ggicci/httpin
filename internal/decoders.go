package internal

import (
	"fmt"
	"mime/multipart"
	"reflect"
	"strconv"
	"time"
)

type ValueTypeDecoder interface {
	Decode(value string) (interface{}, error)
}

type FileTypeDecoder interface {
	Decode(file *multipart.FileHeader) (interface{}, error)
}

type ValueTypeDecoderFunc func(string) (interface{}, error)
type FileTypeDecoderFunc func(*multipart.FileHeader) (interface{}, error)

// Decode calls fn(data).
func (fn ValueTypeDecoderFunc) Decode(value string) (interface{}, error) {
	return fn(value)
}

func (fn FileTypeDecoderFunc) Decode(file *multipart.FileHeader) (interface{}, error) {
	return fn(file)
}

var builtinDecoders = map[reflect.Type]interface{}{
	reflect.TypeOf(true):               ValueTypeDecoderFunc(DecodeBool),
	reflect.TypeOf(int(0)):             ValueTypeDecoderFunc(DecodeInt),
	reflect.TypeOf(int8(0)):            ValueTypeDecoderFunc(DecodeInt8),
	reflect.TypeOf(int16(0)):           ValueTypeDecoderFunc(DecodeInt16),
	reflect.TypeOf(int32(0)):           ValueTypeDecoderFunc(DecodeInt32),
	reflect.TypeOf(int64(0)):           ValueTypeDecoderFunc(DecodeInt64),
	reflect.TypeOf(uint(0)):            ValueTypeDecoderFunc(DecodeUint),
	reflect.TypeOf(uint8(0)):           ValueTypeDecoderFunc(DecodeUint8),
	reflect.TypeOf(uint16(0)):          ValueTypeDecoderFunc(DecodeUint16),
	reflect.TypeOf(uint32(0)):          ValueTypeDecoderFunc(DecodeUint32),
	reflect.TypeOf(uint64(0)):          ValueTypeDecoderFunc(DecodeUint64),
	reflect.TypeOf(float32(0.0)):       ValueTypeDecoderFunc(DecodeFloat32),
	reflect.TypeOf(float64(0.0)):       ValueTypeDecoderFunc(DecodeFloat64),
	reflect.TypeOf(complex64(0 + 1i)):  ValueTypeDecoderFunc(DecodeComplex64),
	reflect.TypeOf(complex128(0 + 1i)): ValueTypeDecoderFunc(DecodeComplex128),
	reflect.TypeOf(string("0")):        ValueTypeDecoderFunc(DecodeString),
	reflect.TypeOf(time.Now()):         ValueTypeDecoderFunc(DecodeTime),
}

func DecoderOf(t reflect.Type) interface{} {
	return builtinDecoders[t]
}

func DecodeBool(value string) (interface{}, error) {
	return strconv.ParseBool(value)
}

func DecodeInt(value string) (interface{}, error) {
	v, err := strconv.ParseInt(value, 10, 64)
	return int(v), err
}

func DecodeInt8(value string) (interface{}, error) {
	v, err := strconv.ParseInt(value, 10, 8)
	return int8(v), err
}

func DecodeInt16(value string) (interface{}, error) {
	v, err := strconv.ParseInt(value, 10, 16)
	return int16(v), err
}

func DecodeInt32(value string) (interface{}, error) {
	v, err := strconv.ParseInt(value, 10, 32)
	return int32(v), err
}

func DecodeInt64(value string) (interface{}, error) {
	v, err := strconv.ParseInt(value, 10, 64)
	return int64(v), err
}

func DecodeUint(value string) (interface{}, error) {
	v, err := strconv.ParseUint(value, 10, 64)
	return uint(v), err
}

func DecodeUint8(value string) (interface{}, error) {
	v, err := strconv.ParseUint(value, 10, 8)
	return uint8(v), err
}

func DecodeUint16(value string) (interface{}, error) {
	v, err := strconv.ParseUint(value, 10, 16)
	return uint16(v), err
}

func DecodeUint32(value string) (interface{}, error) {
	v, err := strconv.ParseUint(value, 10, 32)
	return uint32(v), err
}

func DecodeUint64(value string) (interface{}, error) {
	v, err := strconv.ParseUint(value, 10, 64)
	return uint64(v), err
}

func DecodeFloat32(value string) (interface{}, error) {
	v, err := strconv.ParseFloat(value, 32)
	return float32(v), err
}

func DecodeFloat64(value string) (interface{}, error) {
	v, err := strconv.ParseFloat(value, 64)
	return float64(v), err
}

func DecodeComplex64(value string) (interface{}, error) {
	v, err := strconv.ParseComplex(value, 64)
	return complex64(v), err
}

func DecodeComplex128(value string) (interface{}, error) {
	v, err := strconv.ParseComplex(value, 128)
	return complex128(v), err
}

func DecodeString(value string) (interface{}, error) {
	return value, nil
}

// DecodeTime parses data bytes as time.Time in UTC timezone.
// Supported formats of the data bytes are:
// 1. RFC3339Nano string, e.g. "2006-01-02T15:04:05-07:00"
// 2. Unix timestamp, e.g. "1136239445"
func DecodeTime(value string) (interface{}, error) {
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
