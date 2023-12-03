package core

import (
	"errors"
	"mime/multipart"
	"reflect"
	"time"

	"github.com/ggicci/httpin/internal"
)

var theBuiltinDecoders = map[reflect.Type]Decoder[any]{
	internal.TypeOf[bool]():       ToAnyDecoder[bool](DecoderFunc[bool](internal.DecodeBool)),
	internal.TypeOf[int]():        ToAnyDecoder[int](DecoderFunc[int](internal.DecodeInt)),
	internal.TypeOf[int8]():       ToAnyDecoder[int8](DecoderFunc[int8](internal.DecodeInt8)),
	internal.TypeOf[int16]():      ToAnyDecoder[int16](DecoderFunc[int16](internal.DecodeInt16)),
	internal.TypeOf[int32]():      ToAnyDecoder[int32](DecoderFunc[int32](internal.DecodeInt32)),
	internal.TypeOf[int64]():      ToAnyDecoder[int64](DecoderFunc[int64](internal.DecodeInt64)),
	internal.TypeOf[uint]():       ToAnyDecoder[uint](DecoderFunc[uint](internal.DecodeUint)),
	internal.TypeOf[uint8]():      ToAnyDecoder[uint8](DecoderFunc[uint8](internal.DecodeUint8)),
	internal.TypeOf[uint16]():     ToAnyDecoder[uint16](DecoderFunc[uint16](internal.DecodeUint16)),
	internal.TypeOf[uint32]():     ToAnyDecoder[uint32](DecoderFunc[uint32](internal.DecodeUint32)),
	internal.TypeOf[uint64]():     ToAnyDecoder[uint64](DecoderFunc[uint64](internal.DecodeUint64)),
	internal.TypeOf[float32]():    ToAnyDecoder[float32](DecoderFunc[float32](internal.DecodeFloat32)),
	internal.TypeOf[float64]():    ToAnyDecoder[float64](DecoderFunc[float64](internal.DecodeFloat64)),
	internal.TypeOf[complex64]():  ToAnyDecoder[complex64](DecoderFunc[complex64](internal.DecodeComplex64)),
	internal.TypeOf[complex128](): ToAnyDecoder[complex128](DecoderFunc[complex128](internal.DecodeComplex128)),
	internal.TypeOf[string]():     ToAnyDecoder[string](DecoderFunc[string](internal.DecodeString)),
	internal.TypeOf[time.Time]():  ToAnyDecoder[time.Time](DecoderFunc[time.Time](internal.DecodeTime)),
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

type DataSource interface {
	string | *multipart.FileHeader
}

type decoderInterface[DT DataSource, RT any] interface {
	Decode(DT) (RT, error)
}

func ToAnyDecoder[T any](decoder Decoder[T]) Decoder[any] {
	if decoder == nil {
		return nil
	}
	return DecoderFunc[any](func(s string) (any, error) {
		return decoder.Decode(s)
	})
}

// SmartDecoder is a decoder that switches the return value of the inner Decoder[DT] to
// WantType. For example, if the inner decoder returns a *T, and WantType is T, then the
// SmartDecoder will return T instead of *T, vice versa.
type SmartDecoder struct {
	Decoder  Decoder[any]
	WantType reflect.Type
}

func NewSmartDecoder(typ reflect.Type, decoder Decoder[any]) Decoder[any] {
	return &SmartDecoder{Decoder: decoder, WantType: typ}
}

func (sd *SmartDecoder) Decode(value string) (any, error) {
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
	return nil, typeMismatchedError(sd.WantType, gotType)
}

func validateDecoder(decoder any) error {
	if decoder == nil || internal.IsNil(reflect.ValueOf(decoder)) {
		return errors.New("nil decoder")
	}
	return nil
}
