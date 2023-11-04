package httpin

import (
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
	case typeT:
		return adaptor.T()
	case typeTSlice:
		return adaptor.TSlice()
	case typePatchT:
		return adaptor.PatchT()
	case typePatchTSlice:
		return adaptor.PatchTSlice()
	default:
		return nil
	}
}

func (adaptor *encoderAdaptor) T() formValueEncoder {
	return (*singleFormValueEncoder)(adaptor)
}

func (adaptor *encoderAdaptor) TSlice() formValueEncoder {
	return (*multiFormValueEncoder)(adaptor)
}

func (adaptor *encoderAdaptor) PatchT() formValueEncoder {
	return (*patchFieldFormValueEncoder)(adaptor)
}

func (adaptor *encoderAdaptor) PatchTSlice() formValueEncoder {
	return (*patchFieldMultiFormValueEncoder)(adaptor)
}

type singleFormValueEncoder encoderAdaptor

// EncodeX encodes value of T to []string, a single string element in the result. Where T is the BaseType.
func (e *singleFormValueEncoder) EncodeX(value reflect.Value) ([]string, error) {
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
		bs, _ := encodeByteSlice(value.Bytes())
		return []string{bs}, nil
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

type patchFieldFormValueEncoder encoderAdaptor

// EncodeX encodes value of patch.Field[T] to []string, where T is the BaseType.
func (e *patchFieldFormValueEncoder) EncodeX(value reflect.Value) ([]string, error) {
	if !value.FieldByName("Valid").Bool() {
		return nil, nil
	}
	innerValue := value.FieldByName("Value")
	return (*singleFormValueEncoder)(e).EncodeX(innerValue)
}

type patchFieldMultiFormValueEncoder encoderAdaptor

// EncodeX encodes value of patch.Field[[]T] to []string, where T is the BaseType.
func (e *patchFieldMultiFormValueEncoder) EncodeX(value reflect.Value) ([]string, error) {
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
	case typeT:
		return toFileEncodersOne(value)
	case typePatchT:
		if !value.FieldByName("Valid").Bool() {
			return nil, nil // skip no file upload: patch.Field.Valid is false
		}
		return toFileEncodersOne(value.FieldByName("Value"))
	case typeTSlice:
		return toFileEncodersMulti(value)
	case typePatchTSlice:
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
