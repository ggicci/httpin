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
	if rtm.Value.Type().Kind() == reflect.Slice {
		return e.extractMulti(rtm, key)
	}

	rtmHelper := &directiveRuntimeHelper{rtm}

	switch decoder := rtmHelper.decoderOf(rtm.Value.Type()).(type) {
	case ValueTypeDecoder:
		if gotValue, interfaceValue, err := decodeValueAt(decoder, e.Form.Value[key], 0); err != nil {
			return fieldError{key, gotValue, err}
		} else {
			rtm.Value.Elem().Set(reflect.ValueOf(interfaceValue))
		}
	case FileTypeDecoder:
		if gotFile, interfaceValue, err := decodeFileAt(decoder, e.Form.File[key], 0); err != nil {
			return fieldError{key, gotFile, err}
		} else {
			rtm.Value.Elem().Set(reflect.ValueOf(interfaceValue))
		}
	default:
		return UnsupportedTypeError{rtm.Value.Type()}
	}

	rtmHelper.DeliverContextValue(FieldSet, true)
	return nil
}

func (e *extractor) extractMulti(rtm *DirectiveRuntime, key string) error {
	var (
		theSlice  reflect.Value
		elemType  = rtm.Value.Type().Elem()
		values    = e.Form.Value[key]
		files     = e.Form.File[key]
		rtmHelper = &directiveRuntimeHelper{rtm}
	)

	switch decoder := rtmHelper.decoderOf(elemType).(type) {
	case ValueTypeDecoder:
		theSlice = reflect.MakeSlice(rtm.Value.Type(), len(values), len(values))
		for i := 0; i < len(values); i++ {
			if _, interfaceValue, err := decodeValueAt(decoder, values, i); err != nil {
				return fieldError{key, values, fmt.Errorf("at index %d: %w", i, err)}
			} else {
				theSlice.Index(i).Set(reflect.ValueOf(interfaceValue))
			}
		}
	case FileTypeDecoder:
		theSlice = reflect.MakeSlice(rtm.Value.Type(), len(files), len(files))
		for i := 0; i < len(files); i++ {
			if _, interfaceValue, err := decodeFileAt(decoder, files, i); err != nil {
				return fieldError{key, files, fmt.Errorf("at index %d: %w", i, err)}
			} else {
				theSlice.Index(i).Set(reflect.ValueOf(interfaceValue))
			}
		}
	default:
		return UnsupportedTypeError{rtm.Value.Type()}
	}

	rtm.Value.Elem().Set(theSlice)
	rtmHelper.DeliverContextValue(FieldSet, true)
	return nil
}

func decodeValueAt(decoder ValueTypeDecoder, values []string, index int) (string, interface{}, error) {
	var gotValue = ""
	if index < len(values) {
		gotValue = values[index]
	}
	res, err := decoder.Decode(gotValue)
	return gotValue, res, err
}

func decodeFileAt(decoder FileTypeDecoder, files []*multipart.FileHeader, index int) (*multipart.FileHeader, interface{}, error) {
	var gotFile *multipart.FileHeader
	if index < len(files) {
		gotFile = files[index]
	}
	res, err := decoder.Decode(gotFile)
	return gotFile, res, err
}
