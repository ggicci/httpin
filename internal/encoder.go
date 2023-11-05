package internal

import (
	"errors"
	"reflect"
)

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
	if encoder == nil || IsNil(reflect.ValueOf(encoder)) {
		return errors.New("nil encoder")
	}
	return nil
}
