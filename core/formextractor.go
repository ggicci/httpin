package core

import (
	"mime/multipart"

	"github.com/ggicci/httpin/codec"
)

type FormExtractor struct {
	Runtime *DirectiveRuntime
	multipart.Form
	KeyNormalizer func(string) string
}

func (e *FormExtractor) Extract(keys ...string) error {
	if len(keys) == 0 {
		keys = e.Runtime.Directive.Argv
	}
	for _, key := range keys {
		if e.KeyNormalizer != nil {
			key = e.KeyNormalizer(key)
		}
		if err := e.extract(key); err != nil {
			return err
		}
	}
	return nil
}

func (e *FormExtractor) extract(key string) error {
	if e.Runtime.IsFieldSet() {
		return nil // skip when already extracted by former directives
	}

	values := e.Form.Value[key]
	files := e.Form.File[key]

	// Quick fail on empty input.
	if len(values) == 0 && len(files) == 0 {
		return nil
	}

	ns := e.Runtime.GetNamespace()
	var sourceValue any
	var err error
	valueType := e.Runtime.Value.Type().Elem()
	if ns.IsFileType(valueType) {
		// When fileDecoder is not nil, it means that the field is a file upload.
		// We should decode files instead of values.
		if len(files) == 0 {
			return nil // skip when no file uploaded
		}
		sourceValue = files

		var decoder FileSliceCodec
		decoder, err = codec.NewFileSliceCodec(e.Runtime.Value.Elem())
		if err == nil {
			err = decoder.FromFileSlice(codec.ToFileHeaderList(files))
		}
	} else {
		if len(values) == 0 {
			return nil // skip when no value given
		}
		sourceValue = values

		var adapt StringCodecAdaptor
		decoderInfo := e.Runtime.GetCustomCoder() // custom decoder, specified by "decoder" directive
		if decoderInfo != nil {
			adapt = decoderInfo.Adaptor
		}
		var decoder StringSliceCodec
		decoder, err = ns.NewStringSliceCodec(e.Runtime.Value.Elem(), adapt)
		if err == nil {
			err = decoder.FromStringSlice(values)
		}
	}

	if err != nil {
		return &fieldError{key, sourceValue, err}
	}
	e.Runtime.MarkFieldSet(true)
	return nil
}
