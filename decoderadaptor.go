package httpin

import (
	"fmt"
	"mime/multipart"
	"reflect"
)

type decoderKindType int

const (
	decoderKindScalar     decoderKindType = iota // T
	decoderKindMulti                             // []T
	decoderKindPatch                             // patch.Field[T]
	decoderKindPatchMulti                        // patch.Field[[]T]
)

type decoderAdaptor[DT DataSource] interface {
	BaseType() reflect.Type            // T
	Scalar(reflect.Type) decoder2D[DT] // takes in a desired return type
	Multi(reflect.Type) decoder2D[DT]
	Patch(reflect.Type) decoder2D[DT]
	PatchMulti(reflect.Type) decoder2D[DT]
	DecoderByKind(kind decoderKindType, returnType reflect.Type) decoder2D[DT]
}

// decoderAdaptorImpl is an implementation of decoderAdaptor.
// It can be adapted to 3 types of Decoder:
//
//   - Call .Scalar() to get a Decoder that can be registered for type T.
//   - Call .Multi() for []T.
//   - Call .Patch() for patch.Field[T].
//   - Call .PatchMulti() for patch.Field[[]T].
//
// Itself is also a ScalarTypeDecoder.
type decoderAdaptorImpl[DT DataSource] struct {
	Decoder[DT]
	baseType reflect.Type
}

func (sva *decoderAdaptorImpl[DT]) BaseType() reflect.Type {
	return sva.baseType
}

func (sva *decoderAdaptorImpl[DT]) Scalar(returnType reflect.Type) decoder2D[DT] {
	return &scalarTypeDecoder[DT]{sva, returnType}
}

func (sva *decoderAdaptorImpl[DT]) Multi(returnType reflect.Type) decoder2D[DT] {
	return &multiTypeDecoder[DT]{sva, returnType}
}

func (sva *decoderAdaptorImpl[DT]) Patch(returnType reflect.Type) decoder2D[DT] {
	return &patchFieldTypeDecoder[DT]{sva, returnType}
}

func (sva *decoderAdaptorImpl[DT]) PatchMulti(returnType reflect.Type) decoder2D[DT] {
	return &patchFieldMultiTypeDecoder[DT]{sva, returnType}
}

func (sva *decoderAdaptorImpl[DT]) DecoderByKind(kind decoderKindType, returnType reflect.Type) decoder2D[DT] {
	switch kind {
	case decoderKindScalar:
		return sva.Scalar(returnType)
	case decoderKindMulti:
		return sva.Multi(returnType)
	case decoderKindPatch:
		return sva.Patch(returnType)
	case decoderKindPatchMulti:
		return sva.PatchMulti(returnType)
	}
	return nil
}

// adaptDecoder adapts a decoder (of Decoder[DT]) to a decoderAdaptor.
// It returns nil if the decoder is not supported.
func adaptDecoder(returnType reflect.Type, decoder interface{}) interface{} {
	switch decoder := decoder.(type) {
	case ValueTypeDecoder:
		return &decoderAdaptorImpl[string]{decoder, returnType}
	case FileTypeDecoder:
		return &decoderAdaptorImpl[*multipart.FileHeader]{decoder, returnType}
	default:
		return nil
	}
}

type scalarTypeDecoder[DT DataSource] struct {
	*decoderAdaptorImpl[DT]
	ReturnType reflect.Type
}

// DecodeX of scalarTypeDecoder[DT] decodes a single value.
// It only decodes the first value in the given slice. Returns T.
func (s *scalarTypeDecoder[DT]) DecodeX(values []DT) (interface{} /* T */, error) {
	return s.Decoder.Decode(values[0])
}

type multiTypeDecoder[DT DataSource] struct {
	*decoderAdaptorImpl[DT]
	ReturnType reflect.Type
}

// DecodeX of multiTypeDecoder[DT] decodes multiple values. Returns []T.
func (m *multiTypeDecoder[DT]) DecodeX(values []DT) (interface{} /* []T */, error) {
	res := reflect.MakeSlice(m.ReturnType, len(values), len(values))
	for i, value := range values {
		if gotValue, err := m.Decoder.Decode(value); err != nil {
			return nil, fmt.Errorf("at index %d: %w", i, err)
		} else {
			res.Index(i).Set(reflect.ValueOf(gotValue))
		}
	}
	return res.Interface(), nil
}

type patchFieldTypeDecoder[DT DataSource] struct {
	*decoderAdaptorImpl[DT]
	ReturnType reflect.Type
}

// DecodeX of patchFieldTypeDecoder[DT] decodes a single value.
// It only decodes the first value in the given slice. Returns patch.Field[T].
func (p *patchFieldTypeDecoder[DT]) DecodeX(values []DT) (interface{} /* patch.Field[T] */, error) {
	res := reflect.New(p.ReturnType)
	if gotValue, err := p.Decoder.Decode(values[0]); err != nil {
		return res.Interface(), err
	} else {
		res.Elem().FieldByName("Value").Set(reflect.ValueOf(gotValue))
		res.Elem().FieldByName("Valid").SetBool(true)
		return res.Elem().Interface(), nil
	}
}

type patchFieldMultiTypeDecoder[DT DataSource] struct {
	*decoderAdaptorImpl[DT]
	ReturnType reflect.Type
}

// DecodeX of patchFieldMultiTypeDecoder[DT] decodes multiple values. Returns patch.Field[[]T].
func (p *patchFieldMultiTypeDecoder[DT]) DecodeX(values []DT) (interface{} /* patch.Field[[]T] */, error) {
	subValue := reflect.MakeSlice(reflect.SliceOf(p.BaseType()), len(values), len(values))
	for i, value := range values {
		if gotValue, err := p.Decoder.Decode(value); err != nil {
			return nil, fmt.Errorf("at index %d: %w", i, err)
		} else {
			subValue.Index(i).Set(reflect.ValueOf(gotValue))
		}
	}
	res := reflect.New(p.ReturnType)
	res.Elem().FieldByName("Value").Set(subValue)
	res.Elem().FieldByName("Valid").SetBool(true)
	return res.Elem().Interface(), nil
}

// typeOf returns the reflect.Type of a given type.
// e.g. typeOf[int]() returns reflect.TypeOf(0)
func typeOf[T any]() reflect.Type {
	var zero [0]T
	return reflect.TypeOf(zero).Elem()
}
