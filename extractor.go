package httpin

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"reflect"
)

type extractor struct {
	multipart.Form

	KeyNormalizer func(string) string
}

func newExtractor(r *http.Request) *extractor {
	var form multipart.Form

	if r.MultipartForm != nil {
		form = *r.MultipartForm
	} else {
		if r.Form != nil {
			form.Value = r.Form
		}
	}

	return &extractor{
		Form:          form,
		KeyNormalizer: nil,
	}
}

func (e *extractor) Execute(ctx *DirectiveRuntime) error {
	for _, key := range ctx.Directive.Argv {
		if e.KeyNormalizer != nil {
			key = e.KeyNormalizer(key)
		}
		if err := e.extract(ctx, key); err != nil {
			return err
		}
	}
	return nil
}

func (e *extractor) extract(rtm *DirectiveRuntime, key string) error {
	if rtm.Context.Value(FieldSet) == true {
		return nil
	}

	values := e.Form.Value[key]
	files := e.Form.File[key]

	// Quick fail on empty input.
	if len(values) == 0 && len(files) == 0 {
		return nil
	}

	valueType := rtm.Value.Type().Elem()
	elemType, decoderKind := scalarElemTypeOf(valueType)
	rtmHelper := directiveRuntimeHelper{rtm}
	adaptor := rtmHelper.decoderOf(elemType)
	var decodedValue interface{}
	var err error

	switch ada := adaptor.(type) {
	case decoderAdaptor[string]:
		decodedValue, err = ada.DecoderByKind(decoderKind).Decode(values)
	case decoderAdaptor[*multipart.FileHeader]:
		decodedValue, err = ada.DecoderByKind(decoderKind).Decode(files)
	default:
		err = UnsupportedTypeError{elemType}
	}

	if err == nil {
		err = setDirectiveRuntimeValue(rtm, decodedValue)
	}
	if err != nil {
		return &fieldError{key, values, err}
	}
	rtmHelper.DeliverContextValue(FieldSet, true)
	return nil
}

func setDirectiveRuntimeValue(rtm *DirectiveRuntime, value interface{}) error {
	newValue := reflect.ValueOf(value)
	targetType := rtm.Value.Type().Elem()
	if newValue.Type().AssignableTo(targetType) {
		rtm.Value.Elem().Set(newValue)
		return nil
	}

	return fmt.Errorf("%w: decoded value is of type %v that not assignable to type %v",
		ErrValueTypeMismatch, newValue.Type(), targetType)
}

// scalarElemTypeOf returns the scalar element type of a given type.
//   - T -> T, decoderKindScalar
//   - []T -> T, decoderKindMulti
//   - patch.Field[T] -> T, decoderKindPatch
//   - patch.Field[[]T] -> T, decoderKindPatchMulti
//
// The given type is gonna use the decoder of the scalar element type to decode
// the input values.
func scalarElemTypeOf(valueType reflect.Type) (reflect.Type, decoderKindType) {
	if valueType.Kind() == reflect.Slice {
		return valueType.Elem(), decoderKindMulti
	}
	if isPatchField(valueType) {
		subElemType, isMulti := patchFieldElemType(valueType)
		if isMulti {
			return subElemType, decoderKindPatchMulti
		} else {
			return subElemType, decoderKindPatch
		}
	}
	return valueType, decoderKindScalar
}
