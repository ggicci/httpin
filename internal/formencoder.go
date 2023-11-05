package internal

import (
	"encoding"
	"fmt"
	"reflect"
)

type FormValueMarshaler interface {
	HttpinFormValue() (string, error)
}

var (
	fallbackEncoder = interfaceEncoder{}
)

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
	if DefaultRegistry.IsFileType(baseType) {
		fileEncoders, err := toFileEncoders(rtm.Value, TypeKind)
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
		encoder = DefaultRegistry.GetEncoder(baseType)
	}

	// As the last resort, use the fallback encoder.
	if encoder == nil {
		encoder = fallbackEncoder
	}

	values, err := AdaptEncoder(baseType, encoder).EncoderByKind(TypeKind).EncodeX(rtm.Value)
	if err != nil {
		return err
	}
	e.Setter(key, values)
	rtm.MarkFieldSet(true)
	return nil
}

// interfaceEncoder utilizes the following interfaces to encode a value in order:
//   - httpin.FormValueMarshaler
//   - encoding.TextMarshaler
//   - fmt.Stringer
type interfaceEncoder struct{}

func (ie interfaceEncoder) Encode(value reflect.Value) (string, error) {
	ivalue := value.Interface()

	if marshaler, ok := ivalue.(FormValueMarshaler); ok {
		return marshaler.HttpinFormValue()
	}

	if marshaler, ok := ivalue.(encoding.TextMarshaler); ok {
		bs, err := marshaler.MarshalText()
		if err != nil {
			return "", err
		}
		return string(bs), nil
	}

	if marshaler, ok := ivalue.(fmt.Stringer); ok {
		return marshaler.String(), nil
	}

	return "", UnsupportedTypeError(value.Type())
}

func fileUploadBuilder(rtm *DirectiveRuntime, files []FileEncoder) error {
	rb := rtm.GetRequestBuilder()
	key := rtm.Directive.Argv[0]
	rb.SetAttachment(key, files)
	rtm.MarkFieldSet(true)
	return nil
}
