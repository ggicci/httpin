package internal

import (
	"fmt"
	"reflect"
	"strconv"
	"time"
)

type Decoder interface {
	Decode([]byte) (interface{}, error)
}

type DecoderFunc func([]byte) (interface{}, error)

// Decode calls fn(data).
func (fn DecoderFunc) Decode(data []byte) (interface{}, error) {
	return fn(data)
}

var builtinDecoders = map[reflect.Type]Decoder{
	reflect.TypeOf(true):               DecoderFunc(DecodeBool),
	reflect.TypeOf(int(0)):             DecoderFunc(DecodeInt),
	reflect.TypeOf(int8(0)):            DecoderFunc(DecodeInt8),
	reflect.TypeOf(int16(0)):           DecoderFunc(DecodeInt16),
	reflect.TypeOf(int32(0)):           DecoderFunc(DecodeInt32),
	reflect.TypeOf(int64(0)):           DecoderFunc(DecodeInt64),
	reflect.TypeOf(uint(0)):            DecoderFunc(DecodeUint),
	reflect.TypeOf(uint8(0)):           DecoderFunc(DecodeUint8),
	reflect.TypeOf(uint16(0)):          DecoderFunc(DecodeUint16),
	reflect.TypeOf(uint32(0)):          DecoderFunc(DecodeUint32),
	reflect.TypeOf(uint64(0)):          DecoderFunc(DecodeUint64),
	reflect.TypeOf(float32(0.0)):       DecoderFunc(DecodeFloat32),
	reflect.TypeOf(float64(0.0)):       DecoderFunc(DecodeFloat64),
	reflect.TypeOf(complex64(0 + 1i)):  DecoderFunc(DecodeComplex64),
	reflect.TypeOf(complex128(0 + 1i)): DecoderFunc(DecodeComplex128),
	reflect.TypeOf(string("0")):        DecoderFunc(DecodeString),
	reflect.TypeOf(time.Now()):         DecoderFunc(DecodeTime),
}

func DecoderOf(t reflect.Type) Decoder {
	return builtinDecoders[t]
}

func DecodeBool(data []byte) (interface{}, error) {
	return strconv.ParseBool(string(data))
}

func DecodeInt(data []byte) (interface{}, error) {
	v, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return nil, err
	}
	return int(v), nil
}

func DecodeInt8(data []byte) (interface{}, error) {
	v, err := strconv.ParseInt(string(data), 10, 8)
	if err != nil {
		return nil, err
	}
	return int8(v), nil
}
func DecodeInt16(data []byte) (interface{}, error) {
	v, err := strconv.ParseInt(string(data), 10, 16)
	if err != nil {
		return nil, err
	}
	return int16(v), nil
}
func DecodeInt32(data []byte) (interface{}, error) {
	v, err := strconv.ParseInt(string(data), 10, 32)
	if err != nil {
		return nil, err
	}
	return int32(v), nil
}
func DecodeInt64(data []byte) (interface{}, error) {
	v, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return nil, err
	}
	return int64(v), nil
}

func DecodeUint(data []byte) (interface{}, error) {
	v, err := strconv.ParseUint(string(data), 10, 64)
	if err != nil {
		return nil, err
	}
	return uint(v), nil
}

func DecodeUint8(data []byte) (interface{}, error) {
	v, err := strconv.ParseUint(string(data), 10, 8)
	if err != nil {
		return nil, err
	}
	return uint8(v), nil
}
func DecodeUint16(data []byte) (interface{}, error) {
	v, err := strconv.ParseUint(string(data), 10, 16)
	if err != nil {
		return nil, err
	}
	return uint16(v), nil
}
func DecodeUint32(data []byte) (interface{}, error) {
	v, err := strconv.ParseUint(string(data), 10, 32)
	if err != nil {
		return nil, err
	}
	return uint32(v), nil
}
func DecodeUint64(data []byte) (interface{}, error) {
	v, err := strconv.ParseUint(string(data), 10, 64)
	if err != nil {
		return nil, err
	}
	return uint64(v), nil
}

func DecodeFloat32(data []byte) (interface{}, error) {
	v, err := strconv.ParseFloat(string(data), 32)
	if err != nil {
		return nil, err
	}
	return float32(v), nil
}

func DecodeFloat64(data []byte) (interface{}, error) {
	v, err := strconv.ParseFloat(string(data), 64)
	if err != nil {
		return nil, err
	}
	return float64(v), nil
}

func DecodeComplex64(data []byte) (interface{}, error) {
	v, err := strconv.ParseComplex(string(data), 64)
	if err != nil {
		return nil, err
	}
	return complex64(v), nil
}

func DecodeComplex128(data []byte) (interface{}, error) {
	v, err := strconv.ParseComplex(string(data), 128)
	if err != nil {
		return nil, err
	}
	return complex128(v), nil
}

func DecodeString(data []byte) (interface{}, error) {
	return string(data), nil
}

// DecodeTime parses data bytes as time.Time in UTC timezone.
// Supported formats of the data bytes are:
// 1. RFC3339Nano string, e.g. "2006-01-02T15:04:05-07:00"
// 2. Unix timestamp, e.g. "1136239445"
func DecodeTime(data []byte) (interface{}, error) {
	value := string(data)

	// Try parsing value as RFC3339 format.
	if t, err := time.Parse(time.RFC3339Nano, value); err == nil {
		return t.UTC(), nil
	}

	// Try parsing value as int64 (timestamp).
	// TODO(ggicci): can support float timestamp, e.g. 1618974933.284368
	if timestamp, err := strconv.ParseInt(value, 10, 64); err == nil {
		return time.Unix(timestamp, 0).UTC(), nil
	}

	return time.Time{}, fmt.Errorf("invalid time value")
}
