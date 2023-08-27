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

	// NOTE(ggicci): Array?
	valueType := rtm.Value.Type().Elem()
	if valueType.Kind() == reflect.Slice {
		return e.extractMulti(rtm, valueType, key)
	}

	rtmHelper := &directiveRuntimeHelper{rtm}
	switch decoder := rtmHelper.decoderOf(valueType).(type) {
	case ValueTypeDecoder:
		if inputValue, interfaceValue, err := decodeValueAt(decoder, e.Form.Value[key], 0); err != nil {
			return &fieldError{key, inputValue, err}
		} else {
			rtm.Value.Elem().Set(reflect.ValueOf(interfaceValue))
		}
	case FileTypeDecoder:
		if inputFile, interfaceValue, err := decodeFileAt(decoder, e.Form.File[key], 0); err != nil {
			return &fieldError{key, inputFile, err}
		} else {
			rtm.Value.Elem().Set(reflect.ValueOf(interfaceValue))
		}
	default:
		return UnsupportedTypeError{valueType}
	}

	rtmHelper.DeliverContextValue(FieldSet, true)
	return nil
}

func (e *extractor) extractMulti(rtm *DirectiveRuntime, sliceType reflect.Type, key string) error {
	var (
		rtmHelper = &directiveRuntimeHelper{rtm}
		theSlice  reflect.Value
		values    = e.Form.Value[key]
		files     = e.Form.File[key]
		elemType  = sliceType.Elem()
	)

	switch decoder := rtmHelper.decoderOf(elemType).(type) {
	case ValueTypeDecoder:
		theSlice = reflect.MakeSlice(sliceType, len(values), len(values))
		for i := 0; i < len(values); i++ {
			if _, interfaceValue, err := decodeValueAt(decoder, values, i); err != nil {
				return &fieldError{key, values, fmt.Errorf("at index %d: %w", i, err)}
			} else {
				theSlice.Index(i).Set(reflect.ValueOf(interfaceValue))
			}
		}
	case FileTypeDecoder:
		theSlice = reflect.MakeSlice(sliceType, len(files), len(files))
		for i := 0; i < len(files); i++ {
			if _, interfaceValue, err := decodeFileAt(decoder, files, i); err != nil {
				return &fieldError{key, files, fmt.Errorf("at index %d: %w", i, err)}
			} else {
				theSlice.Index(i).Set(reflect.ValueOf(interfaceValue))
			}
		}
	default:
		return UnsupportedTypeError{sliceType}
	}

	rtm.Value.Elem().Set(theSlice)
	rtmHelper.DeliverContextValue(FieldSet, true)
	return nil
}

func decodeValueAt(decoder ValueTypeDecoder, values []string, index int) (string, interface{}, error) {
	var inputValue = ""
	if index < len(values) {
		inputValue = values[index]
	}
	res, err := decoder.Decode(inputValue)
	return inputValue, res, err
}

func decodeFileAt(decoder FileTypeDecoder, files []*multipart.FileHeader, index int) (*multipart.FileHeader, interface{}, error) {
	var inputFile *multipart.FileHeader
	if index < len(files) {
		inputFile = files[index]
	}
	res, err := decoder.Decode(inputFile)
	return inputFile, res, err
}
