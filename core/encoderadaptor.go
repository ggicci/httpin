package core

import (
	"reflect"

	"github.com/ggicci/httpin/internal"
)

var byteType = internal.TypeOf[byte]()

type FormValueEncoder interface {
	EncodeX(value reflect.Value) ([]string, error)
}

type EncoderAdaptor struct {
	BaseType    reflect.Type
	BaseEncoder Encoder
}

func AdaptEncoder(baseType reflect.Type, encoder Encoder) *EncoderAdaptor {
	return &EncoderAdaptor{
		BaseType:    baseType,
		BaseEncoder: encoder,
	}
}

func (adaptor *EncoderAdaptor) EncoderByKind(kind TypeKind) FormValueEncoder {
	switch kind {
	case TypeKindT:
		return adaptor.T()
	case TypeKindTSlice:
		return adaptor.TSlice()
	case TypeKindPatchT:
		return adaptor.PatchT()
	case TypeKindPatchTSlice:
		return adaptor.PatchTSlice()
	default:
		return nil
	}
}

func (adaptor *EncoderAdaptor) T() FormValueEncoder {
	return (*singleFormValueEncoder)(adaptor)
}

func (adaptor *EncoderAdaptor) TSlice() FormValueEncoder {
	return (*multiFormValueEncoder)(adaptor)
}

func (adaptor *EncoderAdaptor) PatchT() FormValueEncoder {
	return (*patchFieldFormValueEncoder)(adaptor)
}

func (adaptor *EncoderAdaptor) PatchTSlice() FormValueEncoder {
	return (*patchFieldMultiFormValueEncoder)(adaptor)
}

type singleFormValueEncoder EncoderAdaptor

// EncodeX encodes value of T to []string, a single string element in the result. Where T is the BaseType.
func (e *singleFormValueEncoder) EncodeX(value reflect.Value) ([]string, error) {
	if s, err := e.BaseEncoder.Encode(value); err != nil {
		return nil, err
	} else {
		return []string{s}, nil
	}
}

type multiFormValueEncoder EncoderAdaptor

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

type patchFieldFormValueEncoder EncoderAdaptor

// EncodeX encodes value of patch.Field[T] to []string, where T is the BaseType.
func (e *patchFieldFormValueEncoder) EncodeX(value reflect.Value) ([]string, error) {
	if !value.FieldByName("Valid").Bool() {
		return nil, nil
	}
	innerValue := value.FieldByName("Value")
	return (*singleFormValueEncoder)(e).EncodeX(innerValue)
}

type patchFieldMultiFormValueEncoder EncoderAdaptor

// EncodeX encodes value of patch.Field[[]T] to []string, where T is the BaseType.
func (e *patchFieldMultiFormValueEncoder) EncodeX(value reflect.Value) ([]string, error) {
	if !value.FieldByName("Valid").Bool() {
		return nil, nil
	}
	innerValue := value.FieldByName("Value")
	return (*multiFormValueEncoder)(e).EncodeX(innerValue)
}
