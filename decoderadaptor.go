package httpin

import (
	"fmt"
	"mime/multipart"

	"github.com/ggicci/httpin/patch"
)

type decoderKindType int

const (
	decoderKindScalar     decoderKindType = iota // T
	decoderKindMulti                             // []T
	decoderKindPatch                             // patch.Field[T]
	decoderKindPatchMulti                        // patch.Field[[]T]
)

type decoderAdaptor[DT DataSource] interface {
	Scalar() decoder2D[DT]
	Multi() decoder2D[DT]
	Patch() decoder2D[DT]
	PatchMulti() decoder2D[DT]
	DecoderByKind(kind decoderKindType) decoder2D[DT]
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
type decoderAdaptorImpl[T any, DT DataSource] struct {
	Decoder[DT]
}

func (sva *decoderAdaptorImpl[T, DT]) Scalar() decoder2D[DT] {
	return &scalarTypeDecoder[T, DT]{sva.Decoder}
}

func (sva *decoderAdaptorImpl[T, DT]) Multi() decoder2D[DT] {
	return &multiTypeDecoder[T, DT]{sva.Decoder}
}

func (sva *decoderAdaptorImpl[T, DT]) Patch() decoder2D[DT] {
	return &patchFieldTypeDecoder[T, DT]{sva.Decoder}
}

func (sva *decoderAdaptorImpl[T, DT]) PatchMulti() decoder2D[DT] {
	return &patchFieldMultiTypeDecoder[T, DT]{sva.Decoder}
}

func (sva *decoderAdaptorImpl[T, DT]) DecoderByKind(kind decoderKindType) decoder2D[DT] {
	switch kind {
	case decoderKindScalar:
		return sva.Scalar()
	case decoderKindMulti:
		return sva.Multi()
	case decoderKindPatch:
		return sva.Patch()
	case decoderKindPatchMulti:
		return sva.PatchMulti()
	}
	return nil
}

func adaptDecoder[T any, DT DataSource](decoder Decoder[DT]) *decoderAdaptorImpl[T, DT] {
	return &decoderAdaptorImpl[T, DT]{decoder}
}

func adaptDecoderX[T any](decoder interface{}) interface{} {
	switch decoder := decoder.(type) {
	case ValueTypeDecoder:
		return adaptDecoder[T, string](decoder)
	case FileTypeDecoder:
		return adaptDecoder[T, *multipart.FileHeader](decoder)
	default:
		return decoder // noop
	}
}

type DecoderFunc[DT DataSource] func(value DT) (interface{}, error)

func (fn DecoderFunc[DT]) Decode(value DT) (interface{}, error) {
	return fn(value)
}

type scalarTypeDecoder[T any, DT DataSource] struct {
	Decoder[DT]
}

func (s *scalarTypeDecoder[T, DT]) Decode(values []DT) (interface{} /* T */, error) {
	if len(values) == 0 {
		var zero DT
		return s.Decoder.Decode(zero) // "" or nil
	}
	return s.Decoder.Decode(values[0])
}

type multiTypeDecoder[T any, DT DataSource] struct {
	Decoder[DT]
}

func (m *multiTypeDecoder[T, DT]) Decode(values []DT) (interface{} /* []T */, error) {
	res := make([]T, len(values))
	for i, v := range values {
		if gotValue, err := m.Decoder.Decode(v); err != nil {
			return nil, fmt.Errorf("at index %d: %w", i, err)
		} else {
			res[i] = gotValue.(T)
		}
	}
	return res, nil
}

type patchFieldTypeDecoder[T any, DT DataSource] struct {
	Decoder[DT]
}

func (p *patchFieldTypeDecoder[T, DT]) Decode(values []DT) (interface{} /* patch.Field[T] */, error) {
	if len(values) == 0 {
		return patch.Field[T]{}, nil
	}
	if gotValue, err := p.Decoder.Decode(values[0]); err != nil {
		return patch.Field[T]{}, err
	} else {
		return patch.Field[T]{Value: gotValue.(T), Valid: true}, nil
	}
}

type patchFieldMultiTypeDecoder[T any, DT DataSource] struct {
	Decoder[DT]
}

func (p *patchFieldMultiTypeDecoder[T, DT]) Decode(values []DT) (interface{} /* patch.Field[[]T] */, error) {
	res := make([]T, len(values))
	for i, v := range values {
		if gotValue, err := p.Decoder.Decode(v); err != nil {
			return nil, fmt.Errorf("at index %d: %w", i, err)
		} else {
			res[i] = gotValue.(T)
		}
	}
	return patch.Field[[]T]{Value: res, Valid: true}, nil
}
