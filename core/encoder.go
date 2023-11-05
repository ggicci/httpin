package core

import (
	"encoding/base64"
	"errors"
	"reflect"
	"time"

	"github.com/ggicci/httpin/internal"
)

var theBuiltinEncoders = map[reflect.Type]any{
	internal.TypeOf[bool]():       EncoderFunc[bool](internal.EncodeBool),
	internal.TypeOf[int]():        EncoderFunc[int](internal.EncodeInt),
	internal.TypeOf[int8]():       EncoderFunc[int8](internal.EncodeInt8),
	internal.TypeOf[int16]():      EncoderFunc[int16](internal.EncodeInt16),
	internal.TypeOf[int32]():      EncoderFunc[int32](internal.EncodeInt32),
	internal.TypeOf[int64]():      EncoderFunc[int64](internal.EncodeInt64),
	internal.TypeOf[uint]():       EncoderFunc[uint](internal.EncodeUint),
	internal.TypeOf[uint8]():      EncoderFunc[uint8](internal.EncodeUint8),
	internal.TypeOf[uint16]():     EncoderFunc[uint16](internal.EncodeUint16),
	internal.TypeOf[uint32]():     EncoderFunc[uint32](internal.EncodeUint32),
	internal.TypeOf[uint64]():     EncoderFunc[uint64](internal.EncodeUint64),
	internal.TypeOf[float32]():    EncoderFunc[float32](internal.EncodeFloat32),
	internal.TypeOf[float64]():    EncoderFunc[float64](internal.EncodeFloat64),
	internal.TypeOf[complex64]():  EncoderFunc[complex64](internal.EncodeComplex64),
	internal.TypeOf[complex128](): EncoderFunc[complex128](internal.EncodeComplex128),
	internal.TypeOf[string]():     EncoderFunc[string](internal.EncodeString),
	internal.TypeOf[time.Time]():  EncoderFunc[time.Time](internal.EncodeTime),
	internal.TypeOf[[]byte]():     EncoderFunc[[]byte](encodeByteSlice), // []byte is a special case
}

// Encoder is a type that can encode a value of type T to a string. It is
// used by the "form", "query", and "header" directives to encode a value.
type Encoder interface {
	Encode(value reflect.Value) (string, error)
}

// EncoderFunc is a function that encodes a value of type T to a string.
// It implements the Encoder interface.
type EncoderFunc[T any] func(value T) (string, error)

func (fn EncoderFunc[T]) Encode(value reflect.Value) (string, error) {
	return fn(value.Interface().(T))
}

// ToPointerEncoder makes an encoder for a type (T) be able to used as an
// encoder for a T's pointer type (*T).
type ToPointerEncoder struct {
	Encoder
}

func (pe ToPointerEncoder) Encode(value reflect.Value) (string, error) {
	return pe.Encoder.Encode(value.Elem())
}

func validateEncoder(encoder any) error {
	if encoder == nil || internal.IsNil(reflect.ValueOf(encoder)) {
		return errors.New("nil encoder")
	}
	return nil
}

func encodeByteSlice(bytes []byte) (string, error) {
	// NOTE: we're using base64.StdEncoding here, not base64.URLEncoding.
	return base64.StdEncoding.EncodeToString(bytes), nil
}
