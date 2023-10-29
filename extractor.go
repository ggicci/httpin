package httpin

import (
	"mime/multipart"
)

type extractor struct {
	Runtime *DirectiveRuntime
	multipart.Form
	KeyNormalizer func(string) string
}

func (e *extractor) Extract(keys ...string) error {
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

func (e *extractor) extract(key string) error {
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
	baseType, typeKind := baseTypeOf(valueType)
	fileDecoder := fileDecoderByType(baseType) // file decoder, for file uploads

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
		decodedValue, err = fileDecoder.DecoderByKind(typeKind, valueType).DecodeX(files)
	} else {
		if len(values) == 0 {
			return nil // skip when no value given
		}
		sourceValue = values

		var decoder *decoderAdaptor[string]
		decoderInfo := e.Runtime.getCustomDecoder() // custom decoder, specified by "decoder" directive
		// Fallback to use default decoders for registered types.
		if decoderInfo != nil {
			decoder = decoderInfo.Adapted
		} else {
			decoder = decoderByType(baseType)
		}

		if decoder != nil {
			decodedValue, err = decoder.DecoderByKind(typeKind, valueType).DecodeX(values)
		} else {
			err = unsupportedTypeError(valueType)
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
