package httpin

import (
	"mime/multipart"

	"github.com/ggicci/httpin/patch"
)

// type validSetter interface {
// 	SetValid(bool)
// }

// func setField(rtm *DirectiveRuntime) error {
// 	if rtm.Context.Value(FieldSet) == true {
// 		elem := rtm.Value.Interface()
// 		if setter, ok := elem.(validSetter); ok {
// 			setter.SetValid(true)
// 		}
// 	}
// 	return nil
// }

// func appendSetFieldDirective(r *owl.Resolver) error {
// 	r.Directives = append(r.Directives, owl.NewDirective("_setfield"))
// 	return nil
// }

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
