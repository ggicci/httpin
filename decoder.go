package httpin

import (
	"errors"
	"fmt"
	"mime/multipart"
	"reflect"
	"time"
)

var (
	builtinDecoders = newPriorityPair()                  // builtin decoders, always registered
	customDecoders  = newPriorityPair()                  // custom decoders (by type)
	namedDecoders   = make(map[string]*namedDecoderInfo) // custom decoders (by name)
)

func init() {
	registerBuiltinDecoder[bool](decodeBool)
	registerBuiltinDecoder[int](decodeInt)
	registerBuiltinDecoder[int8](decodeInt8)
	registerBuiltinDecoder[int16](decodeInt16)
	registerBuiltinDecoder[int32](decodeInt32)
	registerBuiltinDecoder[int64](decodeInt64)
	registerBuiltinDecoder[uint](decodeUint)
	registerBuiltinDecoder[uint8](decodeUint8)
	registerBuiltinDecoder[uint16](decodeUint16)
	registerBuiltinDecoder[uint32](decodeUint32)
	registerBuiltinDecoder[uint64](decodeUint64)
	registerBuiltinDecoder[float32](decodeFloat32)
	registerBuiltinDecoder[float64](decodeFloat64)
	registerBuiltinDecoder[complex64](decodeComplex64)
	registerBuiltinDecoder[complex128](decodeComplex128)
	registerBuiltinDecoder[string](decodeString)
	registerBuiltinDecoder[time.Time](decodeTime)
}

type (
	Decoder[T any]         decoderInterface[string, T]
	DecoderFunc[T any]     func(string) (T, error)
	FileDecoder[T any]     decoderInterface[*multipart.FileHeader, T]
	FileDecoderFunc[T any] func(*multipart.FileHeader) (T, error)
)

func (fn DecoderFunc[T]) Decode(value string) (T, error) {
	return fn(value)
}

func (fn FileDecoderFunc[T]) Decode(fh *multipart.FileHeader) (T, error) {
	return fn(fh)
}

// RegisterDecoder registers a FormValueDecoder. The decoder takes in a string
// value and decodes it to value of type T. Used to decode values from
// querystring, form, header. Panics on conflict types and nil decoders.
//
// NOTE: the decoder returns the decoded value as any. For best practice, the
// underlying type of the decoded value should be T, even returning *T also
// works. If the real returned value were not T or *T, the decoder will return
// an error (ErrValueTypeMismatch) while decoding.
func RegisterDecoder[T any](decoder Decoder[T], replace ...bool) {
	force := len(replace) > 0 && replace[0]
	registerDecoderTo[T](customDecoders, decoder, force)
}

// RegisterNamedDecoder registers a decoder by name. Panics on conflict names
// and invalid decoders. It decodes a string value to a value of type T. Use the
// "decoder" directive to override the decoder of a struct field:
//
//	func init() {
//	    httpin.RegisterNamedDecoder[time.Time]("x_time", DecoderFunc[string](decodeTimeInXFormat))
//	}
//
//	type Input struct {
//	    // httpin will use the decoder registered above, instead of the builtin decoder for time.Time.
//	    Time time.Time `in:"query:time;decoder=x_time"`
//	}
//
// Visit https://ggicci.github.io/httpin/directives/decoder for more details.
func RegisterNamedDecoder[T any](name string, decoder Decoder[T], replace ...bool) {
	force := len(replace) > 0 && replace[0]
	if _, ok := namedDecoders[name]; ok && !force {
		panicOnError(fmt.Errorf("duplicate name: %q", name))
	}
	panicOnError(validateDecoder(decoder))
	typ := typeOf[T]()

	namedDecoders[name] = &namedDecoderInfo{
		Name:     name,
		Original: decoder,
		Adapted:  adaptDecoder(typ, newSmartDecoder(typ, toAnyDecoder(decoder))).(valueDecoderAdaptor),
	}
}

type dataSource interface {
	string | *multipart.FileHeader
}

type decoderInterface[DT dataSource, RT any] interface {
	Decode(DT) (RT, error)
}

func registerBuiltinDecoder[T any](fn DecoderFunc[T]) {
	registerDecoderTo[T](builtinDecoders, fn, false)
}

func registerDecoderTo[T any](p priorityPair, decoder Decoder[T], force bool) {
	registerDecoderToX[T](p, toAnyDecoder(decoder), force)
}

func registerDecoderToX[T any](p priorityPair, decoder Decoder[any], force bool) {
	panicOnError(validateDecoder(decoder))

	typ := typeOf[T]()
	primaryDecoder := adaptDecoder(typ, newSmartDecoder(typ, decoder))
	panicOnError(p.SetPair(typ, primaryDecoder, nil, force))

	if typ.Kind() == reflect.Pointer {
		// When we have a pointer type (*T), we also register the decoder for its base
		// type (T). The decoder for the base type (T) will be registered as the
		// secondary decoder.
		baseType := typ.Elem()
		secondaryDecoder := adaptDecoder(baseType, newSmartDecoder(baseType, decoder))
		panicOnError(p.SetPair(baseType, nil, secondaryDecoder, force))
	} else {
		// When we have a non-pointer type (T), we also register the decoder for its
		// pointer type (*T). The decoder for the pointer type (*T) will be registered
		// as the secondary decoder.
		pointerType := reflect.PtrTo(typ)
		secondaryDecoder := adaptDecoder(pointerType, newSmartDecoder(pointerType, decoder))
		panicOnError(p.SetPair(pointerType, nil, secondaryDecoder, force))
	}
}

func toAnyDecoder[T any](decoder Decoder[T]) Decoder[any] {
	if decoder == nil {
		return nil
	}
	return DecoderFunc[any](func(s string) (any, error) {
		return decoder.Decode(s)
	})
}

// smartDecoder is a decoder that switches the return value of the inner Decoder[DT] to
// WantType. For example, if the inner decoder returns a *T, and WantType is T, then the
// smartDecoder will return T instead of *T, vice versa.
type smartDecoder struct {
	Decoder  Decoder[any]
	WantType reflect.Type
}

func newSmartDecoder(typ reflect.Type, decoder Decoder[any]) Decoder[any] {
	return &smartDecoder{Decoder: decoder, WantType: typ}
}

func (sd *smartDecoder) Decode(value string) (any, error) {
	gotValue, err := sd.Decoder.Decode(value)
	if err != nil {
		return nil, err
	}
	if gotValue == nil {
		return nil, nil // nil value, return directly
	}

	gotType := reflect.TypeOf(gotValue) // returns nil if gotValue is nil

	// Returns directly on the same type.
	if gotType == sd.WantType {
		return gotValue, nil
	}

	// Want T, got *T, return T.
	if gotType.Kind() == reflect.Ptr && gotType.Elem() == sd.WantType {
		rv := reflect.ValueOf(gotValue)
		if rv.IsNil() {
			return nil, nil
		}
		return rv.Elem().Interface(), nil
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

func validateDecoder(decoder any) error {
	if decoder == nil || isNil(reflect.ValueOf(decoder)) {
		return errors.New("nil decoder")
	}
	return nil
}

// decoderByName retrieves a decoder by name, from the global registerred named decoders.
func decoderByName(name string) *namedDecoderInfo {
	return namedDecoders[name]
}

// decoderByType retrieves a decoder by type, from the global registerred decoders.
func decoderByType(t reflect.Type) valueDecoderAdaptor {
	d := customDecoders.GetOne(t)
	if d == nil {
		d = builtinDecoders.GetOne(t)
	}
	if d != nil {
		return d.(valueDecoderAdaptor)
	} else {
		return nil
	}
}

type namedDecoderInfo struct {
	Name     string
	Original any
	Adapted  valueDecoderAdaptor
}
