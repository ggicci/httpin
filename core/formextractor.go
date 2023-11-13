package core

import (
	"mime/multipart"
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

	valueType := e.Runtime.Value.Type().Elem()
	baseType, TypeKind := BaseTypeOf(valueType)
	fileDecoder := defaultRegistry.GetFileDecoder(baseType) // file decoder, for file uploads

	var decodedValue any
	var sourceValue any
	var err error

	if fileDecoder != nil {
		// When fileDecoder is not nil, it means that the field is a file upload.
		// We should decode files instead of values.
		if len(files) == 0 {
			return nil // skip when no file uploaded
		}
		sourceValue = files

		// Adapt the fileDecoder which is for the baseType, to a fileDecoder
		// which is for the valueType.
		decodedValue, err = fileDecoder.DecoderByKind(TypeKind, valueType).DecodeX(files)
	} else {
		if len(values) == 0 {
			return nil // skip when no value given
		}
		sourceValue = values

		var decoder ValueDecoderAdaptor
		decoderInfo := e.Runtime.getCustomDecoder() // custom decoder, specified by "decoder" directive
		// Fallback to use default decoders for registered types.
		if decoderInfo != nil {
			decoder = decoderInfo.Adapted
		} else {
			decoder = defaultRegistry.GetDecoder(baseType)
		}

		if decoder != nil {
			decodedValue, err = decoder.DecoderByKind(TypeKind, valueType).DecodeX(values)
		} else {
			err = UnsupportedTypeError(valueType)
		}
	}

	if err == nil {
		err = e.Runtime.SetValue(decodedValue)
	}
	if err != nil {
		return &fieldError{key, sourceValue, err}
	}

	e.Runtime.MarkFieldSet(true)
	return nil
}
