package httpin

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"reflect"
)

type Extractor struct {
	multipart.Form

	KeyNormalizer func(string) string
}

func NewExtractor(r *http.Request) *Extractor {
	var form multipart.Form

	if r.MultipartForm != nil {
		form = *r.MultipartForm
	} else {
		if r.Form != nil {
			form.Value = r.Form
		}
	}

	return &Extractor{
		Form:          form,
		KeyNormalizer: nil,
	}
}

func (e *Extractor) Execute(ctx *DirectiveContext) error {
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

func (e *Extractor) extract(ctx *DirectiveContext, key string) error {
	if ctx.Context.Value(FieldSet) == true {
		return nil
	}

	values := e.Form.Value[key]
	files := e.Form.File[key]
	if len(values) == 0 && len(files) == 0 {
		return nil
	}

	// NOTE(ggicci): Array?
	if ctx.ValueType.Kind() == reflect.Slice {
		return e.extractMulti(ctx, key)
	}

	decoder := decoderOf(ctx.ValueType)
	if decoder == nil {
		return UnsupportedTypeError{ctx.ValueType}
	}

	switch dec := decoder.(type) {
	case ValueTypeDecoder:
		if gotValue, interfaceValue, err := decodeValueAt(dec, e.Form.Value[key], 0); err != nil {
			return fieldError{key, gotValue, err}
		} else {
			ctx.Value.Elem().Set(reflect.ValueOf(interfaceValue))
		}
	case FileTypeDecoder:
		if gotFile, interfaceValue, err := decodeFileAt(dec, e.Form.File[key], 0); err != nil {
			return fieldError{key, gotFile, err}
		} else {
			ctx.Value.Elem().Set(reflect.ValueOf(interfaceValue))
		}
	default:
		return UnsupportedTypeError{ctx.ValueType}
	}

	ctx.DeliverContextValue(FieldSet, true)
	return nil
}

func (e *Extractor) extractMulti(ctx *DirectiveContext, key string) error {
	elemType := ctx.ValueType.Elem()
	decoder := decoderOf(elemType)
	if decoder == nil {
		return UnsupportedTypeError{ctx.ValueType}
	}

	var theSlice reflect.Value
	values := e.Form.Value[key]
	files := e.Form.File[key]

	switch dec := decoder.(type) {
	case ValueTypeDecoder:
		theSlice = reflect.MakeSlice(ctx.ValueType, len(values), len(values))
		for i := 0; i < len(values); i++ {
			if _, interfaceValue, err := decodeValueAt(dec, values, i); err != nil {
				return fieldError{key, values, fmt.Errorf("at index %d: %w", i, err)}
			} else {
				theSlice.Index(i).Set(reflect.ValueOf(interfaceValue))
			}
		}
	case FileTypeDecoder:
		theSlice = reflect.MakeSlice(ctx.ValueType, len(files), len(files))
		for i := 0; i < len(files); i++ {
			if _, interfaceValue, err := decodeFileAt(dec, files, i); err != nil {
				return fieldError{key, files, fmt.Errorf("at index %d: %w", i, err)}
			} else {
				theSlice.Index(i).Set(reflect.ValueOf(interfaceValue))
			}
		}
	default:
		return UnsupportedTypeError{ctx.ValueType}
	}

	ctx.Value.Elem().Set(theSlice)
	ctx.DeliverContextValue(FieldSet, true)
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
