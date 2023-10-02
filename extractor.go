package httpin

import (
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

	rtmHelper := directiveRuntimeHelper{rtm}
	valueType := rtm.Value.Type().Elem()
	elemType, decoderKind := scalarElemTypeOf(valueType)
	decoder := rtmHelper.decoderOf(elemType)

	var decodedValue interface{}
	var err error

	switch ada := decoder.(type) {
	case decoderAdaptor[string]:
		decodedValue, err = ada.DecoderByKind(decoderKind, valueType).DecodeX(values)
	case decoderAdaptor[*multipart.FileHeader]:
		decodedValue, err = ada.DecoderByKind(decoderKind, valueType).DecodeX(files)
	default:
		err = unsupportedTypeError(elemType)
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
	if value == nil {
		// NOTE: should we wipe the value here? i.e. set the value to nil if necessary.
		// No case found yet, at lease for now.
		return nil
	}
	newValue := reflect.ValueOf(value)
	targetType := rtm.Value.Type().Elem()
	if newValue.Type().AssignableTo(targetType) {
		rtm.Value.Elem().Set(newValue)
		return nil
	}
	return invalidDecodeReturnType(targetType, reflect.TypeOf(value))
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
