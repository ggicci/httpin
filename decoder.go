package httpin

import (
	"fmt"
	"mime/multipart"
	"reflect"
	"time"
)

// primarySecondaryDecodersMap is a map of type to a pair of decoders. The first
// decoder is the primary decoder, while the second decoder is the secondary.
type primarySecondaryDecodersMap map[reflect.Type][2]interface{}

var (
	builtinDecoders = make(primarySecondaryDecodersMap) // builtin decoders, always registered
	customDecoders  = make(primarySecondaryDecodersMap) // custom decoders (by type)
	namedDecoders   = make(map[string]interface{})      // custom decoders (by name)
)

func init() {
	registerTypeDecoderTo[bool](builtinDecoders, DecoderFunc[string](decodeBool), false)
	registerTypeDecoderTo[int](builtinDecoders, DecoderFunc[string](decodeInt), false)
	registerTypeDecoderTo[int8](builtinDecoders, DecoderFunc[string](decodeInt8), false)
	registerTypeDecoderTo[int16](builtinDecoders, DecoderFunc[string](decodeInt16), false)
	registerTypeDecoderTo[int32](builtinDecoders, DecoderFunc[string](decodeInt32), false)
	registerTypeDecoderTo[int64](builtinDecoders, DecoderFunc[string](decodeInt64), false)
	registerTypeDecoderTo[uint](builtinDecoders, DecoderFunc[string](decodeUint), false)
	registerTypeDecoderTo[uint8](builtinDecoders, DecoderFunc[string](decodeUint8), false)
	registerTypeDecoderTo[uint16](builtinDecoders, DecoderFunc[string](decodeUint16), false)
	registerTypeDecoderTo[uint32](builtinDecoders, DecoderFunc[string](decodeUint32), false)
	registerTypeDecoderTo[uint64](builtinDecoders, DecoderFunc[string](decodeUint64), false)
	registerTypeDecoderTo[float32](builtinDecoders, DecoderFunc[string](decodeFloat32), false)
	registerTypeDecoderTo[float64](builtinDecoders, DecoderFunc[string](decodeFloat64), false)
	registerTypeDecoderTo[complex64](builtinDecoders, DecoderFunc[string](decodeComplex64), false)
	registerTypeDecoderTo[complex128](builtinDecoders, DecoderFunc[string](decodeComplex128), false)
	registerTypeDecoderTo[string](builtinDecoders, DecoderFunc[string](decodeString), false)
	registerTypeDecoderTo[time.Time](builtinDecoders, DecoderFunc[string](decodeTime), false)
}

// DataSource is the type of the input data. It can be string or *multipart.FileHeader.
//   - string: when the input data is from a querystring, form, or header.
//   - *multipart.FileHeader: when the input data is from a file upload.
type DataSource interface{ string | *multipart.FileHeader }

// Decoder is the interface implemented by types that can decode a DataSource to
// themselves.
type Decoder[DT DataSource] interface {
	Decode(value DT) (interface{}, error)
}

// ValueTypeDecoder is the interface implemented by types that can decode a
// string to themselves. Take querystring as an example, the decoder takes in a
// single string value and decodes it to value of type T.
type ValueTypeDecoder = Decoder[string]

// FileTypeDecoder is the interface implemented by types that can decode a
// *multipart.FileHeader to themselves.
type FileTypeDecoder = Decoder[*multipart.FileHeader]

// DecoderFunc is a function that implements Decoder[DT]. It can be used to turn
// a function into a Decoder[DT]. For instance:
//
//	func decodeInt(value string) (interface{}, error) { ... }
//	myIntDecoder := DecoderFunc[string](decodeInt)
type DecoderFunc[DT DataSource] func(value DT) (interface{}, error)

func (fn DecoderFunc[DT]) Decode(value DT) (interface{}, error) {
	return fn(value)
}

// decoder2D is the interface implemented by types that can decode a slice of
// DataSource to themselves. DecodeX[DT] takes in a slice of DT values and
// decodes them to some type of value. DecodeX[DT] is usually derived from
// Decoder[DT], using Decoder[DT] to decode each element of the slice.
type decoder2D[DT DataSource] interface {
	DecodeX(values []DT) (interface{}, error)
}

// RegisterValueTypeDecoder registers a ValueTypeDecoder. The decoder takes in a
// string value and decodes it to value of type T. Panics on conflict types and
// nil decoders.
//
// NOTE: the decoder returns the decoded value as interface{}. For best
// practice, the underlying type of the decoded value should be T, even
// returning *T also works. If the real returned value were not T or *T, the
// decoder will return an error (ErrValueTypeMismatch) while decoding.
func RegisterValueTypeDecoder[T any](decoder Decoder[string], replace ...bool) {
	force := len(replace) > 0 && replace[0]
	registerTypeDecoderTo[T](customDecoders, decoder, force)
}

// RegisterFileTypeDecoder registers a FileTypeDecoder. The decoder takes in a
// *multipart.FileHeader (when uploading files from an HTTP request) and decodes
// it to value of type T. Panics on conflict types and nil decoders.
//
// NOTE: the decoder returns the decoded value as interface{}. For best
// practice, the underlying type of the decoded value should be T, even
// returning *T also works. If the real returned value were not T or *T, the
// decoder will return an error (ErrValueTypeMismatch) while decoding.
func RegisterFileTypeDecoder[T any](decoder Decoder[*multipart.FileHeader], replace ...bool) {
	force := len(replace) > 0 && replace[0]
	registerTypeDecoderTo[T](customDecoders, decoder, force)
}

// RegisterNamedDecoder registers a decoder by name. Panics on conflict names
// and invalid decoders. The decoder can be a ValueTypeDecoder or a
// FileTypeDecoder. It decodes the input value to a value of type T. Use the
// *decoder directive* to override the decoder of a struct field:
//
//	RegisterNamedDecoder[time.Time]("x_time", DecoderFunc[string](decodeTimeInXFormat))
//	type Input struct {
//	    // httpin will use the decoder registered above, instead of the builtin decoder for time.Time.
//	    Time time.Time `in:"query:time;decoder=x_time"`
//	}
//
// Visit https://ggicci.github.io/httpin/directives/decoder for more details.
func RegisterNamedDecoder[T any](name string, decoder interface{}, replace ...bool) {
	force := len(replace) > 0 && replace[0]
	if _, ok := namedDecoders[name]; ok && !force {
		panic(fmt.Errorf("httpin: %w: %q", ErrDuplicateNamedDecoder, name))
	}

	panicOnInvalidDecoder(decoder)
	typ := typeOf[T]()
	namedDecoders[name] = adaptDecoder(typ, newSmartDecoderX(typ, decoder))
}

func registerTypeDecoderTo[T any](m primarySecondaryDecodersMap, decoder interface{}, force bool) {
	typ := typeOf[T]()
	panicOnInvalidDecoder(decoder)

	primaryDecoder := adaptDecoder(typ, newSmartDecoderX(typ, decoder))
	updateDecodersMap(m, typ, primaryDecoder, nil, force)

	if typ.Kind() == reflect.Pointer {
		// When we have a pointer type (*T), we also register the decoder for
		// its base type (T). The decoder for the base type will be registered
		// as the secondary decoder.
		baseType := typ.Elem()
		secondaryDecoder := adaptDecoder(baseType, newSmartDecoderX(baseType, decoder))
		updateDecodersMap(m, baseType, nil, secondaryDecoder, force)
	} else {
		// When we have a non-pointer type (T), we also register the decoder
		// for its pointer type (*T). The decoder for the pointer type will be
		// registered as the secondary decoder.
		pointerType := reflect.PtrTo(typ)
		secondaryDecoder := adaptDecoder(pointerType, newSmartDecoderX(pointerType, decoder))
		updateDecodersMap(m, pointerType, nil, secondaryDecoder, force)
	}
}

// updateDecodersMap updates the decoders map with the given primary and
// secondary decoder. The given nil decoders will be ignored. The secondary
// decoder is always set. While the primary decoder is only set when the primary
// decoder of the given type is not set or force is true.
func updateDecodersMap(m map[reflect.Type][2]interface{}, typ reflect.Type, primary, secondary interface{}, force bool) {
	olds, ok := m[typ]

	if !ok {
		m[typ] = [2]interface{}{primary, secondary}
		return
	}

	oldPrimary, _ := olds[0], olds[1]
	if primary != nil { // set primary
		if oldPrimary != nil && !force { // conflict
			panic(fmt.Errorf("httpin: %w: %q", ErrDuplicateTypeDecoder, typ))
		}
		olds[0] = primary
	}

	if secondary != nil { // always set secondary
		olds[1] = secondary
	}
}

// smartDecoder is a decoder that switches the return value of the inner
// Decoder[DT] to WantType. For example, if the inner decoder returns a *T, and
// WantType is T, then the smartDecoder will return T instead of *T, vice versa.
type smartDecoder[DT DataSource] struct {
	Decoder[DT]
	WantType reflect.Type
}

func newSmartDecoder[DT DataSource](typ reflect.Type, decoder interface{}) Decoder[DT] {
	return &smartDecoder[DT]{decoder.(Decoder[DT]), typ}
}

func newSmartDecoderX(typ reflect.Type, decoder interface{}) interface{} {
	switch decoder := decoder.(type) {
	case ValueTypeDecoder:
		return newSmartDecoder[string](typ, decoder)
	case FileTypeDecoder:
		return newSmartDecoder[*multipart.FileHeader](typ, decoder)
	default:
		return nil
	}
}

func (sd *smartDecoder[DT]) Decode(value DT) (interface{}, error) {
	if gotValue, err := sd.Decoder.Decode(value); err != nil {
		return nil, err
	} else {
		// Returns directly on nil.
		if gotValue == nil {
			return nil, nil
		}

		gotType := reflect.TypeOf(gotValue)

		// Returns directly on the same type.
		if gotType == sd.WantType {
			return gotValue, nil
		}

		// Want T, got *T, return T.
		if gotType.Kind() == reflect.Ptr && gotType.Elem() == sd.WantType {
			return reflect.ValueOf(gotValue).Elem().Interface(), nil
		}

		// Want *T, got T, return &T.
		if sd.WantType.Kind() == reflect.Ptr && sd.WantType.Elem() == gotType {
			res := reflect.New(gotType)
			res.Elem().Set(reflect.ValueOf(gotValue))
			return res.Interface(), nil
		}

		// Can't convert, return error.
		return nil, invalidDecodeReturnType(sd.WantType, gotType)
	}
}

// panicOnInvalidDecoder panics when the decoder is invalid, check by
// validateDecoder.
func panicOnInvalidDecoder(decoder interface{}) {
	if err := validateDecoder(decoder); err != nil {
		panic(fmt.Errorf("httpin: %w", err))
	}
}

// validateDecoder validates the decoder. It returns an error if the decoder is
// invalid, otherwise nil.
//  1. nil decoder --> ErrNilDecoder
//  2. not a ValueTypeDecoder or a FileTypeDecoder --> ErrInvalidDecoder
func validateDecoder(decoder interface{}) error {
	if decoder == nil {
		return ErrNilDecoder
	}
	if !isDecoder(decoder) {
		return ErrInvalidDecoder
	}
	return nil
}

// decoderByName retrieves a decoder by name, from the global registerred named decoders.
func decoderByName(name string) interface{} {
	return namedDecoders[name]
}

// decoderByType retrieves a decoder by type, from the global registerred decoders.
func decoderByType(t reflect.Type) interface{} {
	if d := decoderByTypeFrom(customDecoders, t); d != nil {
		return d
	}
	return decoderByTypeFrom(builtinDecoders, t)
}

// decoderByTypeFrom retrieves a decoder by type, from a specific decoders map.
// It prioritizes the primary decoder over the secondary decoder.
func decoderByTypeFrom(m primarySecondaryDecodersMap, t reflect.Type) interface{} {
	if decoders, ok := m[t]; ok {
		if decoders[0] != nil {
			return decoders[0]
		}
		if decoders[1] != nil {
			return decoders[1]
		}
	}
	return nil
}

// isDecoder checks if the decoder is a ValueTypeDecoder or a FileTypeDecoder.
func isDecoder(decoder interface{}) bool {
	_, isValueTypeDecoder := decoder.(Decoder[string])
	if isValueTypeDecoder {
		return true
	}
	_, isFileTypeDecoder := decoder.(Decoder[*multipart.FileHeader])
	return isFileTypeDecoder
}
