package httpin

import (
	"encoding/base64"
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
	innerValue := value.FieldByName("Value")
	return (*scalarFormValueEncoder)(e).EncodeX(innerValue)
}

type patchMultiFormValueEncoder encoderAdaptor

// EncodeX encodes value of patch.Field[[]T] to []string, where T is the BaseType.
func (e *patchMultiFormValueEncoder) EncodeX(value reflect.Value) ([]string, error) {
	innerValue := value.FieldByName("Value")
	return (*multiFormValueEncoder)(e).EncodeX(innerValue)
}

func toFileEncoders(value reflect.Value, typeKind typeKind) []FileEncoder {
	switch typeKind {
	case typeKindScalar:
		return toFileEncodersOne(value)
	case typeKindPatch:
		return toFileEncodersOne(value.FieldByName("Value"))
	case typeKindMulti:
		return toFileEncodersMulti(value)
	case typeKindPatchMulti:
		return toFileEncodersMulti(value.FieldByName("Value"))
	}
	return nil
}

func toFileEncodersOne(one reflect.Value) []FileEncoder {
	return []FileEncoder{one.Interface().(FileEncoder)}
}

func toFileEncodersMulti(multi reflect.Value) []FileEncoder {
	files := make([]FileEncoder, multi.Len())
	for i := 0; i < multi.Len(); i++ {
		files[i] = multi.Index(i).Interface().(FileEncoder)
	}
	return files
}
