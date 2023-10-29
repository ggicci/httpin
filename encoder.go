package httpin

import (
	"encoding"
	"errors"
	"fmt"
	"reflect"
	"time"
)

var (
	builtinEncoders = newPriorityPair()                  // builtin encoders, always registered
	customEncoders  = newPriorityPair()                  // custom encoders (by type)
	namedEncoders   = make(map[string]*namedEncoderInfo) // custom encoders (by name)
	fallbackEncoder = interfaceEncoder{}
)

func init() {
	registerBuiltinEncoder[bool](encodeBool)
	registerBuiltinEncoder[int](encodeInt)
	registerBuiltinEncoder[int8](encodeInt8)
	registerBuiltinEncoder[int16](encodeInt16)
	registerBuiltinEncoder[int32](encodeInt32)
	registerBuiltinEncoder[int64](encodeInt64)
	registerBuiltinEncoder[uint](encodeUint)
	registerBuiltinEncoder[uint8](encodeUint8)
	registerBuiltinEncoder[uint16](encodeUint16)
	registerBuiltinEncoder[uint32](encodeUint32)
	registerBuiltinEncoder[uint64](encodeUint64)
	registerBuiltinEncoder[float32](encodeFloat32)
	registerBuiltinEncoder[float64](encodeFloat64)
	registerBuiltinEncoder[complex64](encodeComplex64)
	registerBuiltinEncoder[complex128](encodeComplex128)
	registerBuiltinEncoder[string](encodeString)
	registerBuiltinEncoder[time.Time](encodeTime)
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

// RegisterEncoder registers a Encoder, which is used to encode a value of type T to a string.
func RegisterEncoder[T any](encoder EncoderFunc[T], replace ...bool) {
	force := len(replace) > 0 && replace[0]
	registerEncoderTo[T](customEncoders, encoder, force)
}

// RegisterNamedEncoder registers a Encoder with a name. The name is used to
// reference the encoder in the "encoder" directive. Ex:
//
//	func init() {
//	    RegisterNamedEncoder[time.Time]("mydate", encodeMyDate)
//	}
//
//	type ListUsersRequest struct {
//	    MemberSince time.Time `in:"query,encoder=mydate"` // use "mydate" encoder
//	}
func RegisterNamedEncoder[T any](name string, encoder EncoderFunc[T], replace ...bool) {
	force := len(replace) > 0 && replace[0]
	if _, ok := namedEncoders[name]; ok && !force {
		panicOnError(fmt.Errorf("duplicate name: %q", name))
	}
	panicOnError(validateEncoder(encoder))
	namedEncoders[name] = &namedEncoderInfo{
		Name:     name,
		Original: encoder,
	}
}

func registerBuiltinEncoder[T any](encoder EncoderFunc[T]) {
	registerEncoderTo[T](builtinEncoders, encoder, false)
}

func registerEncoderTo[T any](p priorityPair, encoder Encoder, force bool) {
	panicOnError(validateEncoder(encoder))

	typ := typeOf[T]()
	panicOnError(p.SetPair(typ, encoder, nil, force))

	if typ.Kind() != reflect.Pointer {
		// When we have a non-pointer type (T), we also register the encoder for its
		// pointer type (*T). The encoder for the pointer type (*T) will be registered
		// as the secondary encoder.
		panicOnError(p.SetPair(reflect.PtrTo(typ), nil, scalar2pointerEncoder{encoder}, force))
	}
}

// scalar2pointerEncoder makes an encoder for a scalar type (T) be able to used as an
// encoder for a pointer type (*T).
type scalar2pointerEncoder struct {
	Encoder
}

func (pe scalar2pointerEncoder) Encode(value reflect.Value) (string, error) {
	return pe.Encoder.Encode(value.Elem())
}

// interfaceEncoder utilizes the following interfaces to encode a value in order:
//   - httpin.FormValueMarshaler
//   - encoding.TextMarshaler
//   - fmt.Stringer
type interfaceEncoder struct{}

func (ie interfaceEncoder) Encode(value reflect.Value) (string, error) {
	ivalue := value.Interface()

	if marshaler, ok := ivalue.(FormValueMarshaler); ok {
		return marshaler.HttpinFormValue()
	}

	if marshaler, ok := ivalue.(encoding.TextMarshaler); ok {
		bs, err := marshaler.MarshalText()
		if err != nil {
			return "", err
		}
		return string(bs), nil
	}

	if marshaler, ok := ivalue.(fmt.Stringer); ok {
		return marshaler.String(), nil
	}

	return "", unsupportedTypeError(value.Type())
}

func validateEncoder(encoder any) error {
	if encoder == nil || reflect.ValueOf(encoder).IsNil() {
		return errors.New("nil encoder")
	}
	return nil
}

type namedEncoderInfo struct {
	Name     string
	Original Encoder
}

func encoderByName(name string) *namedEncoderInfo {
	return namedEncoders[name]
}

func encoderByType(t reflect.Type) Encoder {
	if e := customEncoders.GetOne(t); e != nil {
		return e.(Encoder)
	}
	if e := builtinEncoders.GetOne(t); e != nil {
		return e.(Encoder)
	}
	return nil
}
