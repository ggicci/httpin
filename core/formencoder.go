package core

type FormEncoder struct {
	Setter func(key string, value []string) // form value setter
}

func (e *FormEncoder) Execute(rtm *DirectiveRuntime) error {
	if rtm.IsFieldSet() {
		return nil // skip when already encoded by former directives
	}

	key := rtm.Directive.Argv[0]
	valueType := rtm.Value.Type()
	baseType, TypeKind := BaseTypeOf(valueType)

	// When baseType is a file type, we treat it as a file upload.
	if defaultRegistry.IsFileType(baseType) {
		fileEncoders, err := toFileEncoders(rtm.Value, TypeKind)
		if err != nil {
			return err
		}
		if len(fileEncoders) == 0 {
			return nil // skip when no file upload
		}
		return fileUploadBuilder(rtm, fileEncoders)
	}

	var adapt AnyStringableAdaptor
	encoderInfo := rtm.getCustomEncoderV2()
	if encoderInfo != nil {
		adapt = encoderInfo.Adapt
	}
	var encoder StringSlicable
	encoder, err := NewStringSlicable(rtm.Value, adapt)
	if err != nil {
		return err
	}

	if values, err := encoder.ToStringSlice(); err != nil {
		return err
	} else {
		e.Setter(key, values)
		rtm.MarkFieldSet(true)
		return nil
	}
}

func fileUploadBuilder(rtm *DirectiveRuntime, files []FileEncoder) error {
	rb := rtm.GetRequestBuilder()
	key := rtm.Directive.Argv[0]
	rb.SetAttachment(key, files)
	rtm.MarkFieldSet(true)
	return nil
}
