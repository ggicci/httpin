package core

import "github.com/ggicci/httpin/internal"

type FormEncoder struct {
	Setter func(key string, value []string) // form value setter
}

func (e *FormEncoder) Execute(rtm *DirectiveRuntime) error {
	if rtm.IsFieldSet() {
		return nil // skip when already encoded by former directives
	}

	key := rtm.Directive.Argv[0]
	valueType := rtm.Value.Type()
	// When baseType is a file type, we treat it as a file upload.
	if isFileType(valueType) {
		if internal.IsNil(rtm.Value) {
			return nil // skip when nil, which means no file uploaded
		}

		encoder, err := NewFileSlicable(rtm.Value)
		if err != nil {
			return err
		}
		files, err := encoder.ToFileSlice()
		if err != nil {
			return err
		}
		if len(files) == 0 {
			return nil // skip when no file uploaded
		}
		return fileUploadBuilder(rtm, files)
	}

	var adapt AnyStringableAdaptor
	encoderInfo := rtm.GetCustomEncoder()
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

func fileUploadBuilder(rtm *DirectiveRuntime, files []FileMarshaler) error {
	rb := rtm.GetRequestBuilder()
	key := rtm.Directive.Argv[0]
	rb.SetAttachment(key, files)
	rtm.MarkFieldSet(true)
	return nil
}
