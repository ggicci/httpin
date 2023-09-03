package httpin

import (
	"fmt"
	"mime/multipart"
	"reflect"
	"strconv"
	"time"

	"github.com/ggicci/httpin/patch"
)

var (
	builtinDecoders = make(map[reflect.Type]interface{}) // builtin decoders, always registered
	customDecoders  = make(map[reflect.Type]interface{}) // custom decoders (by type)
	namedDecoders   = make(map[string]interface{})       // custom decoders (by name)
)

func init() {
	registerTypeDecoderTo[bool](builtinDecoders, ValueTypeDecoderFunc(decodeBool))
	registerTypeDecoderTo[int](builtinDecoders, ValueTypeDecoderFunc(decodeInt))
	registerTypeDecoderTo[int8](builtinDecoders, ValueTypeDecoderFunc(decodeInt8))
	registerTypeDecoderTo[int16](builtinDecoders, ValueTypeDecoderFunc(decodeInt16))
	registerTypeDecoderTo[int32](builtinDecoders, ValueTypeDecoderFunc(decodeInt32))
	registerTypeDecoderTo[int64](builtinDecoders, ValueTypeDecoderFunc(decodeInt64))
	registerTypeDecoderTo[uint](builtinDecoders, ValueTypeDecoderFunc(decodeUint))
	registerTypeDecoderTo[uint8](builtinDecoders, ValueTypeDecoderFunc(decodeUint8))
	registerTypeDecoderTo[uint16](builtinDecoders, ValueTypeDecoderFunc(decodeUint16))
	registerTypeDecoderTo[uint32](builtinDecoders, ValueTypeDecoderFunc(decodeUint32))
	registerTypeDecoderTo[uint64](builtinDecoders, ValueTypeDecoderFunc(decodeUint64))
	registerTypeDecoderTo[float32](builtinDecoders, ValueTypeDecoderFunc(decodeFloat32))
	registerTypeDecoderTo[float64](builtinDecoders, ValueTypeDecoderFunc(decodeFloat64))
	registerTypeDecoderTo[complex64](builtinDecoders, ValueTypeDecoderFunc(decodeComplex64))
	registerTypeDecoderTo[complex128](builtinDecoders, ValueTypeDecoderFunc(decodeComplex128))
	registerTypeDecoderTo[string](builtinDecoders, ValueTypeDecoderFunc(decodeString))
	registerTypeDecoderTo[time.Time](builtinDecoders, ValueTypeDecoderFunc(decodeTime))
}

// ValueTypeDecoder is the interface implemented by types that can decode a
// string to themselves.
type ValueTypeDecoder interface {
	Decode(value string) (interface{}, error)
}

// ValueTypeDecoderFunc is an adaptor to allow the use of ordinary functions as
// httpin `ValueTypeDecoder`s.
type ValueTypeDecoderFunc func(string) (interface{}, error)

func (fn ValueTypeDecoderFunc) Decode(value string) (interface{}, error) {
	return fn(value)
}

// FileTypeDecoder is the interface implemented by types that can decode a
// *multipart.FileHeader to themselves.
type FileTypeDecoder interface {
	Decode(file *multipart.FileHeader) (interface{}, error)
}

// FileTypeDecoderFunc is an adaptor to allow the use of ordinary functions as
// httpin `FileTypeDecoder`s.
type FileTypeDecoderFunc func(*multipart.FileHeader) (interface{}, error)

func (fn FileTypeDecoderFunc) Decode(file *multipart.FileHeader) (interface{}, error) {
	return fn(file)
}

func isTypeDecoder(decoder interface{}) bool {
	_, isValueTypeDecoder := decoder.(ValueTypeDecoder)
	_, isFileTypeDecoder := decoder.(FileTypeDecoder)
	return isValueTypeDecoder || isFileTypeDecoder
}

// RegisterDecoder registers a specific type decoder.
// The decoder can be type of `ValueTypeDecoder` or `FileTypeDecoder`.
// Panics on conflicts or invalid decoder.
func RegisterTypeDecoder[T any](decoder interface{}) {
	registerTypeDecoderTo[T](customDecoders, decoder)
}

func registerTypeDecoderTo[T any](m map[reflect.Type]interface{}, decoder interface{}) {
	replaceTypeDecoderTo[T](m, decoder, false)
	// Always register patch.Field[T] by force, as long as T is registered.
	replaceTypeDecoderTo[patch.Field[T]](m, wrapDecoderForPatchField[T](decoder), true)
}

func ReplaceTypeDecoder[T any](decoder interface{}) {
	replaceTypeDecoderTo[T](customDecoders, decoder, true)
	replaceTypeDecoderTo[patch.Field[T]](customDecoders, wrapDecoderForPatchField[T](decoder), true)
}

func replaceTypeDecoderTo[T any](m map[reflect.Type]interface{}, decoder interface{}, force bool) {
	var zero [0]T
	typ := reflect.TypeOf(zero).Elem()
	ensureValidDecoder(typ, decoder)

	if _, conflict := m[typ]; conflict && !force {
		panic(fmt.Errorf("httpin: %w: %q", ErrDuplicateTypeDecoder, typ))
	}
	m[typ] = decoder
}

// RegisterNamedDecoder registers a decoder by name. Panics on conflicts.
func RegisterNamedDecoder(name string, decoder interface{}) {
	if _, ok := namedDecoders[name]; ok {
		panic(fmt.Errorf("httpin: %w: %q", ErrDuplicateNamedDecoder, name))
	}

	ReplaceNamedDecoder(name, decoder)
}

// ReplaceNamedDecoder replaces a decoder by name.
func ReplaceNamedDecoder(name string, decoder interface{}) {
	ensureValidDecoder(nil, decoder)
	namedDecoders[name] = decoder
}

func ensureValidDecoder(typ reflect.Type, decoder interface{}) {
	if decoder == nil {
		panic(fmt.Errorf("httpin: %w: %q", ErrNilTypeDecoder, typ))
	}

	if !isTypeDecoder(decoder) {
		panic(fmt.Errorf("httpin: %w: %q", ErrInvalidTypeDecoder, typ))
	}
}

// decoderOf retrieves a decoder by type, from the global registerred decoders.
func decoderOf(t reflect.Type) interface{} {
	if decoder := customDecoders[t]; decoder != nil {
		return decoder
	} else {
		return builtinDecoders[t]
	}
}

// decoderByName retrieves a decoder by name, from the global registerred named decoders.
func decoderByName(name string) interface{} {
	return namedDecoders[name]
}

// All the builtin decoders:

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
