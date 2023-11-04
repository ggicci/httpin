package httpin

import (
	"fmt"
	"mime/multipart"
	"reflect"
)

// decoder2d is the interface implemented by types that can decode a slice of
// DataSource to themselves. DecodeX[DT] takes in a slice of DT values and
// decodes them to some type of value. DecodeX[DT] is usually derived from
// Decoder[DT], using Decoder[DT] to decode each element of the slice.
type decoder2d[DT dataSource] interface {
	DecodeX(values []DT) (any, error)
}

type (
	valueDecoderAdaptor = *decoderAdaptor[string]
	fileDecoderAdaptor  = *decoderAdaptor[*multipart.FileHeader]
)

type decoderAdaptor[DT dataSource] struct {
	BaseDecoder decoderInterface[DT, any]
	BaseType    reflect.Type
}

// adaptDecoder adapts a decoder of baseType to a decoderAdaptor.
func adaptDecoder(baseType reflect.Type, decoder any) any {
	switch decoder := decoder.(type) {
	case Decoder[any]:
		return &decoderAdaptor[string]{decoder, baseType}
	case FileDecoder[any]:
		return &decoderAdaptor[*multipart.FileHeader]{decoder, baseType}
	default:
		return nil
	}
}

// T returns an adapted decoder for type T.
func (sva *decoderAdaptor[DT]) T(returnType reflect.Type) decoder2d[DT] {
	return &singleTypeDecoder[DT]{sva, returnType}
}

// TSlice returns an adapted decoder for []T.
func (sva *decoderAdaptor[DT]) TSlice(returnType reflect.Type) decoder2d[DT] {
	return &multiTypeDecoder[DT]{sva, returnType}
}

// PatchT returns an adapted decoder for patch.Field[T].
func (sva *decoderAdaptor[DT]) PatchT(returnType reflect.Type) decoder2d[DT] {
	return &patchFieldTypeDecoder[DT]{sva, returnType}
}

// PatchTSlice returns an adapted decoder for patch.Field[[]T].
func (sva *decoderAdaptor[DT]) PatchTSlice(returnType reflect.Type) decoder2d[DT] {
	return &patchFieldMultiTypeDecoder[DT]{sva, returnType}
}

func (sva *decoderAdaptor[DT]) DecoderByKind(kind typeKind, returnType reflect.Type) decoder2d[DT] {
	switch kind {
	case typeT:
		return sva.T(returnType)
	case typeTSlice:
		return sva.TSlice(returnType)
	case typePatchT:
		return sva.PatchT(returnType)
	case typePatchTSlice:
		return sva.PatchTSlice(returnType)
	default:
		return nil
	}
}

type singleTypeDecoder[DT dataSource] struct {
	*decoderAdaptor[DT]
	ReturnType reflect.Type
}

// DecodeX of singleTypeDecoder decodes a single value.
// It only decodes the first value in the given slice. Returns T.
func (s *singleTypeDecoder[DT]) DecodeX(values []DT) (any /* T */, error) {
	return s.BaseDecoder.Decode(values[0])
}

type multiTypeDecoder[DT dataSource] struct {
	*decoderAdaptor[DT]
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

type patchFieldTypeDecoder[DT dataSource] struct {
	*decoderAdaptor[DT]
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

type patchFieldMultiTypeDecoder[DT dataSource] struct {
	*decoderAdaptor[DT]
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
