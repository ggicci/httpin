package httpin

import (
	"encoding/base64"
	"errors"
	"fmt"
	"reflect"
)

var byteType = typeOf[byte]()

type formValueEncoder interface {
	EncodeX(value reflect.Value) ([]string, error)
}

type encoderAdaptor struct {
	BaseType    reflect.Type
	BaseEncoder Encoder
}

func adaptEncoder(baseType reflect.Type, encoder Encoder) *encoderAdaptor {
	return &encoderAdaptor{
		BaseType:    baseType,
		BaseEncoder: encoder,
	}
}

func (adaptor *encoderAdaptor) EncoderByKind(kind typeKind) formValueEncoder {
	switch kind {
	case typeKindScalar:
		return adaptor.Scalar()
	case typeKindMulti:
		return adaptor.Multi()
	case typeKindPatch:
		return adaptor.Patch()
	case typeKindPatchMulti:
		return adaptor.PatchMulti()
	default:
		return nil
	}
}

func (adaptor *encoderAdaptor) Scalar() formValueEncoder {
	return (*scalarFormValueEncoder)(adaptor)
}

func (adaptor *encoderAdaptor) Multi() formValueEncoder {
	return (*multiFormValueEncoder)(adaptor)
}

func (adaptor *encoderAdaptor) Patch() formValueEncoder {
	return (*patchFormValueEncoder)(adaptor)
}

func (adaptor *encoderAdaptor) PatchMulti() formValueEncoder {
	return (*patchMultiFormValueEncoder)(adaptor)
}

type scalarFormValueEncoder encoderAdaptor

// EncodeX encodes value of T to []string, a single string element in the result. Where T is the BaseType.
func (e *scalarFormValueEncoder) EncodeX(value reflect.Value) ([]string, error) {
	if s, err := e.BaseEncoder.Encode(value); err != nil {
		return nil, err
	} else {
		return []string{s}, nil
	}
}

type multiFormValueEncoder encoderAdaptor

// EncodeX encodes value of []T to []string, where T is the BaseType.
func (e *multiFormValueEncoder) EncodeX(value reflect.Value) ([]string, error) {
	// Special case: []byte => base64.
	if e.BaseType == byteType {
		return []string{base64.URLEncoding.EncodeToString(value.Bytes())}, nil
	}

	var res []string
	for i := 0; i < value.Len(); i++ {
		if s, err := e.BaseEncoder.Encode(value.Index(i)); err != nil {
			return nil, err
		} else {
			res = append(res, s)
		}
	}
	return res, nil
}

type patchFormValueEncoder encoderAdaptor

// EncodeX encodes value of patch.Field[T] to []string, where T is the BaseType.
func (e *patchFormValueEncoder) EncodeX(value reflect.Value) ([]string, error) {
	if !value.FieldByName("Valid").Bool() {
		return nil, nil
	}
	innerValue := value.FieldByName("Value")
	return (*scalarFormValueEncoder)(e).EncodeX(innerValue)
}

type patchMultiFormValueEncoder encoderAdaptor

// EncodeX encodes value of patch.Field[[]T] to []string, where T is the BaseType.
func (e *patchMultiFormValueEncoder) EncodeX(value reflect.Value) ([]string, error) {
	if !value.FieldByName("Valid").Bool() {
		return nil, nil
	}
	innerValue := value.FieldByName("Value")
	return (*multiFormValueEncoder)(e).EncodeX(innerValue)
}

func toFileEncoders(value reflect.Value, typeKind typeKind) ([]FileEncoder, error) {
	if isNil(value) {
		return nil, nil // skip no file upload: value is nil
	}

	switch typeKind {
	case typeKindScalar:
		return toFileEncodersOne(value)
	case typeKindPatch:
		if !value.FieldByName("Valid").Bool() {
			return nil, nil // skip no file upload: patch.Field.Valid is false
		}
		return toFileEncodersOne(value.FieldByName("Value"))
	case typeKindMulti:
		return toFileEncodersMulti(value)
	case typeKindPatchMulti:
		if !value.FieldByName("Valid").Bool() {
			return nil, nil // skip no file upload: patch.Field.Valid is false
		}
		return toFileEncodersMulti(value.FieldByName("Value"))
	}
	return nil, nil
}

func toFileEncodersOne(one reflect.Value) ([]FileEncoder, error) {
	if err := validateFileEncoderValue(one); err != nil {
		return nil, err
	}
	return []FileEncoder{one.Interface().(FileEncoder)}, nil
}

func toFileEncodersMulti(multi reflect.Value) ([]FileEncoder, error) {
	files := make([]FileEncoder, multi.Len())
	for i := 0; i < multi.Len(); i++ {
		if err := validateFileEncoderValue(multi.Index(i)); err != nil {
			return nil, fmt.Errorf("at index %d: %v", i, err)
		} else {
			files[i] = multi.Index(i).Interface().(FileEncoder)
		}
	}
	return files, nil
}

func validateFileEncoderValue(value reflect.Value) error {
	if isNil(value) {
		return errors.New("file encoder cannot be nil")
	}
	return nil
}
