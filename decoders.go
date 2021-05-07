package httpin

import (
	"fmt"
	"reflect"
	"strconv"
	"time"
)

type Decoder interface {
	Decode([]byte, reflect.Value) error
}

type DecoderFunc func([]byte, reflect.Value) error

func (fn DecoderFunc) Decode(data []byte, rv reflect.Value) error {
	return fn(data, rv)
}

var builtinDecoders = map[reflect.Type]Decoder{
	reflect.TypeOf(true):               DecoderFunc(DecodeBool),
	reflect.TypeOf(int(0)):             DecoderFunc(DecodeInt),
	reflect.TypeOf(int8(0)):            DecoderFunc(DecodeInt),
	reflect.TypeOf(int16(0)):           DecoderFunc(DecodeInt),
	reflect.TypeOf(int32(0)):           DecoderFunc(DecodeInt),
	reflect.TypeOf(int64(0)):           DecoderFunc(DecodeInt),
	reflect.TypeOf(uint(0)):            DecoderFunc(DecodeUint),
	reflect.TypeOf(uint8(0)):           DecoderFunc(DecodeUint),
	reflect.TypeOf(uint16(0)):          DecoderFunc(DecodeUint),
	reflect.TypeOf(uint32(0)):          DecoderFunc(DecodeUint),
	reflect.TypeOf(uint64(0)):          DecoderFunc(DecodeUint),
	reflect.TypeOf(float32(0.0)):       DecoderFunc(DecodeFloat),
	reflect.TypeOf(float64(0.0)):       DecoderFunc(DecodeFloat),
	reflect.TypeOf(complex64(0 + 1i)):  DecoderFunc(DecodeComplex),
	reflect.TypeOf(complex128(0 + 1i)): DecoderFunc(DecodeComplex),
	reflect.TypeOf(string("0")):        DecoderFunc(DecodeString),
	reflect.TypeOf(time.Now()):         DecoderFunc(DecodeTime),
}

var customDecoders = map[reflect.Type]Decoder{}

func decoderOf(t reflect.Type) Decoder {
	dec := customDecoders[t]
	if dec != nil {
		return dec
	}
	return builtinDecoders[t]
}

func DecodeBool(data []byte, rv reflect.Value) error {
	v, err := strconv.ParseBool(string(data))
	if err != nil {
		return err
	}
	rv.SetBool(v)
	return nil
}

func DecodeInt(data []byte, rv reflect.Value) error {
	v, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return err
	}
	rv.SetInt(v)
	return nil
}

func DecodeUint(data []byte, rv reflect.Value) error {
	v, err := strconv.ParseUint(string(data), 10, 64)
	if err != nil {
		return err
	}
	rv.SetUint(v)
	return nil
}

func DecodeFloat(data []byte, rv reflect.Value) error {
	v, err := strconv.ParseFloat(string(data), 64)
	if err != nil {
		return err
	}
	rv.SetFloat(v)
	return nil
}

func DecodeComplex(data []byte, rv reflect.Value) error {
	v, err := strconv.ParseComplex(string(data), 128)
	if err != nil {
		return err
	}
	rv.SetComplex(v)
	return nil
}

func DecodeString(data []byte, rv reflect.Value) error {
	rv.SetString(string(data))
	return nil
}

func parseTime(value string) (time.Time, error) {
	// Try parsing value as RFC3339 format.
	if t, err := time.Parse(time.RFC3339Nano, value); err == nil {
		return t.UTC(), nil
	}

	// Try parsing value as int64 (timestamp).
	// TODO(ggicci): can support float timestamp, e.g. 1618974933.284368
	if timestamp, err := strconv.ParseInt(value, 10, 64); err == nil {
		return time.Unix(timestamp, 0).UTC(), nil
	}

	return time.Time{}, fmt.Errorf("invalid time value, use time.RFC3339Nano format or timestamp")
}

func DecodeTime(data []byte, rv reflect.Value) error {
	timeValue, err := parseTime(string(data))
	if err != nil {
		return err
	}
	rv.Set(reflect.ValueOf(timeValue))
	return nil
}
