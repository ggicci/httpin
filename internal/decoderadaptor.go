package internal

import (
	"fmt"
	"mime/multipart"
	"reflect"
)

// Decoder2D is the interface implemented by types that can decode a slice of
// DataSource to themselves. DecodeX[DT] takes in a slice of DT values and
// decodes them to some type of value. DecodeX[DT] is usually derived from
// Decoder[DT], using Decoder[DT] to decode each element of the slice.
type Decoder2D[DT DataSource] interface {
	DecodeX(values []DT) (any, error)
}

type (
	ValueDecoderAdaptor = *DecoderAdaptor[string]
	FileDecoderAdaptor  = *DecoderAdaptor[*multipart.FileHeader]
)

type DecoderAdaptor[DT DataSource] struct {
	BaseDecoder decoderInterface[DT, any]
	BaseType    reflect.Type
}

// AdaptDecoder adapts a decoder of baseType to a DecoderAdaptor.
func AdaptDecoder(baseType reflect.Type, decoder any) any {
	switch decoder := decoder.(type) {
	case Decoder[any]:
		return &DecoderAdaptor[string]{decoder, baseType}
	case FileDecoder[any]:
		return &DecoderAdaptor[*multipart.FileHeader]{decoder, baseType}
	default:
		return nil
	}
}

// T returns an adapted decoder for type T.
func (sva *DecoderAdaptor[DT]) T(returnType reflect.Type) Decoder2D[DT] {
	return &singleTypeDecoder[DT]{sva, returnType}
}

// TSlice returns an adapted decoder for []T.
func (sva *DecoderAdaptor[DT]) TSlice(returnType reflect.Type) Decoder2D[DT] {
	return &multiTypeDecoder[DT]{sva, returnType}
}

// PatchT returns an adapted decoder for patch.Field[T].
func (sva *DecoderAdaptor[DT]) PatchT(returnType reflect.Type) Decoder2D[DT] {
	return &patchFieldTypeDecoder[DT]{sva, returnType}
}

// PatchTSlice returns an adapted decoder for patch.Field[[]T].
func (sva *DecoderAdaptor[DT]) PatchTSlice(returnType reflect.Type) Decoder2D[DT] {
	return &patchFieldMultiTypeDecoder[DT]{sva, returnType}
}

func (sva *DecoderAdaptor[DT]) DecoderByKind(kind TypeKind, returnType reflect.Type) Decoder2D[DT] {
	switch kind {
	case TypeKindT:
		return sva.T(returnType)
	case TypeKindTSlice:
		return sva.TSlice(returnType)
	case TypeKindPatchT:
		return sva.PatchT(returnType)
	case TypeKindPatchTSlice:
		return sva.PatchTSlice(returnType)
	default:
		return nil
	}
}

type singleTypeDecoder[DT DataSource] struct {
	*DecoderAdaptor[DT]
	ReturnType reflect.Type
}

// DecodeX of singleTypeDecoder decodes a single value.
// It only decodes the first value in the given slice. Returns T.
func (s *singleTypeDecoder[DT]) DecodeX(values []DT) (any /* T */, error) {
	return s.BaseDecoder.Decode(values[0])
}

type multiTypeDecoder[DT DataSource] struct {
	*DecoderAdaptor[DT]
	ReturnType reflect.Type
}

// DecodeX of multiTypeDecoder decodes multiple values. Returns []T.
func (m *multiTypeDecoder[DT]) DecodeX(values []DT) (any /* []T */, error) {
	res := reflect.MakeSlice(m.ReturnType, len(values), len(values))
	for i, value := range values {
		if gotValue, err := m.BaseDecoder.Decode(value); err != nil {
			return nil, fmt.Errorf("at index %d: %w", i, err)
		} else {
			res.Index(i).Set(reflect.ValueOf(gotValue))
		}
	}
	return res.Interface(), nil
}

type patchFieldTypeDecoder[DT DataSource] struct {
	*DecoderAdaptor[DT]
	ReturnType reflect.Type
}

// DecodeX of patchFieldTypeDecoder decodes a single value.
// It only decodes the first value in the given slice. Returns patch.Field[T].
func (p *patchFieldTypeDecoder[DT]) DecodeX(values []DT) (any /* patch.Field[T] */, error) {
	res := reflect.New(p.ReturnType)
	if gotValue, err := p.BaseDecoder.Decode(values[0]); err != nil {
		return res.Interface(), err
	} else {
		res.Elem().FieldByName("Value").Set(reflect.ValueOf(gotValue))
		res.Elem().FieldByName("Valid").SetBool(true)
		return res.Elem().Interface(), nil
	}
}

type patchFieldMultiTypeDecoder[DT DataSource] struct {
	*DecoderAdaptor[DT]
	ReturnType reflect.Type
}

// DecodeX of patchFieldMultiTypeDecoder decodes multiple values. Returns patch.Field[[]T].
func (pm *patchFieldMultiTypeDecoder[DT]) DecodeX(values []DT) (any /* patch.Field[[]T] */, error) {
	subValue := reflect.MakeSlice(reflect.SliceOf(pm.BaseType), len(values), len(values))
	for i, value := range values {
		if gotValue, err := pm.BaseDecoder.Decode(value); err != nil {
			return nil, fmt.Errorf("at index %d: %w", i, err)
		} else {
			subValue.Index(i).Set(reflect.ValueOf(gotValue))
		}
	}
	res := reflect.New(pm.ReturnType)
	res.Elem().FieldByName("Value").Set(subValue)
	res.Elem().FieldByName("Valid").SetBool(true)
	return res.Elem().Interface(), nil
}
