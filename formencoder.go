package httpin

type FormValueMarshaler interface {
	HttpinFormValue() (string, error)
}

type formEncoder struct {
	Setter func(key string, value []string) // form value setter
}

func (e *formEncoder) Execute(rtm *DirectiveRuntime) error {
	if rtm.IsFieldSet() {
		return nil // skip when already encoded by former directives
	}

	key := rtm.Directive.Argv[0]
	valueType := rtm.Value.Type()
	baseType, typeKind := baseTypeOf(valueType)

	// When baseType is a file type, we treat it as a file upload.
	if isFileType(baseType) {
		fileEncoders, err := toFileEncoders(rtm.Value, typeKind)
		if err != nil {
			return err
		}
		if len(fileEncoders) == 0 {
			return nil // skip when no file upload
		}
		return fileUploadBuilder(rtm, fileEncoders)
	}

	_, encoder := rtm.GetCustomEncoder() // custom encoder, specified by "encoder" directive
	// If no named encoder specified, check if there is a custom encoder for the
	// type of this field, if so, use it.
	if encoder == nil {
		encoder = encoderByType(baseType)
	}

	// As the last resort, use the fallback encoder.
	if encoder == nil {
		encoder = fallbackEncoder
	}

	values, err := adaptEncoder(baseType, encoder).EncoderByKind(typeKind).EncodeX(rtm.Value)
	if err != nil {
		return err
	}
	e.Setter(key, values)
	rtm.MarkFieldSet(true)
	return nil
}
