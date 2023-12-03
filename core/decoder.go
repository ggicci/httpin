package core

import (
	"errors"
	"mime/multipart"
	"reflect"

	"github.com/ggicci/httpin/internal"
)

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
