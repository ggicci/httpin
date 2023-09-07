package httpin

import (
	"fmt"
	"mime/multipart"
	"reflect"
	"time"
)

var (
	builtinDecoders = make(map[reflect.Type]interface{}) // builtin decoders, always registered
	customDecoders  = make(map[reflect.Type]interface{}) // custom decoders (by type)
	namedDecoders   = make(map[string]interface{})       // custom decoders (by name)
)

func init() {
	registerTypeDecoderTo[bool](builtinDecoders, AdaptDecoderFunc[bool, string](decodeBool), false)
	registerTypeDecoderTo[int](builtinDecoders, AdaptDecoderFunc[int, string](decodeInt), false)
	registerTypeDecoderTo[int8](builtinDecoders, AdaptDecoderFunc[int8, string](decodeInt8), false)
	registerTypeDecoderTo[int16](builtinDecoders, AdaptDecoderFunc[int16, string](decodeInt16), false)
	registerTypeDecoderTo[int32](builtinDecoders, AdaptDecoderFunc[int32, string](decodeInt32), false)
	registerTypeDecoderTo[int64](builtinDecoders, AdaptDecoderFunc[int64, string](decodeInt64), false)
	registerTypeDecoderTo[uint](builtinDecoders, AdaptDecoderFunc[uint, string](decodeUint), false)
	registerTypeDecoderTo[uint8](builtinDecoders, AdaptDecoderFunc[uint8, string](decodeUint8), false)
	registerTypeDecoderTo[uint16](builtinDecoders, AdaptDecoderFunc[uint16, string](decodeUint16), false)
	registerTypeDecoderTo[uint32](builtinDecoders, AdaptDecoderFunc[uint32, string](decodeUint32), false)
	registerTypeDecoderTo[uint64](builtinDecoders, AdaptDecoderFunc[uint64, string](decodeUint64), false)
	registerTypeDecoderTo[float32](builtinDecoders, AdaptDecoderFunc[float32, string](decodeFloat32), false)
	registerTypeDecoderTo[float64](builtinDecoders, AdaptDecoderFunc[float64, string](decodeFloat64), false)
	registerTypeDecoderTo[complex64](builtinDecoders, AdaptDecoderFunc[complex64, string](decodeComplex64), false)
	registerTypeDecoderTo[complex128](builtinDecoders, AdaptDecoderFunc[complex128, string](decodeComplex128), false)
	registerTypeDecoderTo[string](builtinDecoders, AdaptDecoderFunc[string, string](decodeString), false)
	registerTypeDecoderTo[time.Time](builtinDecoders, AdaptDecoderFunc[time.Time, string](decodeTime), false)
}

type DataSource interface{ string | *multipart.FileHeader }

type Decoder[DT DataSource] interface {
	Decode(value DT) (interface{}, error)
}

type ValueTypeDecoder = Decoder[string]
type FileTypeDecoder = Decoder[*multipart.FileHeader]

// decoder2D is the interface implemented by types that can decode a slice of
// DataSource to themselves. DataSource can be string or *multipart.FileHeader.
type decoder2D[DT DataSource] interface {
	Decode(values []DT) (interface{}, error)
}

// RegisterDecoder registers a specific type decoder. The decoder can be a
// TypeDecoder or a ScalarTypeDecoder.
//
// When the decoder is a ScalarTypeDecoder, it will be adapted to 3 decoders
// and will be registered to T, []T and patch.Field[T] separately.
//
// When the decoder is a TypeDecoder, it will be registered to T only.
//
// Panics on conflicts or invalid decoder.
func registerTypeDecoder[T any, DT DataSource](decoder Decoder[DT]) {
	registerTypeDecoderTo[T](customDecoders, decoder, false)
}

func RegisterValueTypeDecoder[T any](decoder Decoder[string]) {
	registerTypeDecoder[T, string](decoder)
}

func RegisterFileTypeDecoder[T any](decoder Decoder[*multipart.FileHeader]) {
	registerTypeDecoder[T, *multipart.FileHeader](decoder)
}

func replaceTypeDecoder[T any, DT DataSource](decoder Decoder[DT]) {
	registerTypeDecoderTo[T](customDecoders, decoder, true)
}

func ReplaceValueTypeDecoder[T any](decoder Decoder[string]) {
	replaceTypeDecoder[T, string](decoder)
}

func ReplaceFileTypeDecoder[T any](decoder Decoder[*multipart.FileHeader]) {
	replaceTypeDecoder[T, *multipart.FileHeader](decoder)
}

func registerTypeDecoderTo[T any](m map[reflect.Type]interface{}, decoder interface{}, force bool) {
	var zero [0]T
	typ := reflect.TypeOf(zero).Elem()
	panicOnInvalidDecoder(decoder)

	if _, conflict := m[typ]; conflict && !force {
		panic(fmt.Errorf("httpin: %w: %q", ErrDuplicateTypeDecoder, typ))
	}

	m[typ] = adaptDecoder[T](decoder)
}

// RegisterNamedDecoder registers a decoder by name. Panics on conflicts.
func RegisterNamedDecoder[T any](name string, decoder interface{}) {
	if _, ok := namedDecoders[name]; ok {
		panic(fmt.Errorf("httpin: %w: %q", ErrDuplicateNamedDecoder, name))
	}

	ReplaceNamedDecoder[T](name, decoder)
}

// ReplaceNamedDecoder replaces a decoder by name.
func ReplaceNamedDecoder[T any](name string, decoder interface{}) {
	panicOnInvalidDecoder(decoder)
	namedDecoders[name] = adaptDecoder[T](decoder)
}

func panicOnInvalidDecoder(decoder interface{}) {
	if decoder == nil {
		panic(fmt.Errorf("httpin: %w", ErrNilDecoder))
	}

	if !isDecoder(decoder) {
		panic(fmt.Errorf("httpin: %w", ErrInvalidDecoder))
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

func isDecoder(decoder interface{}) bool {
	_, isValueTypeDecoder := decoder.(Decoder[string])
	_, isFileTypeDecoder := decoder.(Decoder[*multipart.FileHeader])
	return isValueTypeDecoder || isFileTypeDecoder
}
