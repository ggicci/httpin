package httpin

import (
	"mime/multipart"

	"github.com/ggicci/httpin/patch"
)

func wrapDecoderForPatchField[T any](decoder interface{}) interface{} {
	switch d := decoder.(type) {
	case ValueTypeDecoder:
		return ValueTypeDecoderFunc(func(value string) (interface{}, error) {
			if gotValue, err := d.Decode(value); err != nil {
				return patch.Field[T]{}, err
			} else {
				return patch.Field[T]{Value: gotValue.(T), Valid: true}, nil
			}
		})
	case FileTypeDecoder:
		return FileTypeDecoderFunc(func(file *multipart.FileHeader) (interface{}, error) {
			if gotValue, err := d.Decode(file); err != nil {
				return patch.Field[T]{}, err
			} else {
				return patch.Field[T]{Value: gotValue.(T), Valid: true}, nil
			}
		})
	default:
		panic("httpin: invalid decoder")
	}
}
