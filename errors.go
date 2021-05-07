package httpin

import (
	"errors"
	"fmt"
	"reflect"
)

var (
	ErrMissingField         = errors.New("field required but missing")
	ErrUnsupporetedType     = errors.New("unsupported type")
	ErrUnregisteredExecutor = errors.New("unregistered executor")
)

type UnsupportedTypeError struct {
	Type reflect.Type
}

func (e UnsupportedTypeError) Error() string {
	return fmt.Sprintf("unsupported type %q", e.Type.String())
}

func (e UnsupportedTypeError) Unwrap() error {
	return ErrUnsupporetedType
}

type InvalidField struct {
	// Field is the name of the field.
	Field string `json:"field"`

	// Source is the directive which causes the error.
	// e.g. form, header, required, etc.
	Source string `json:"source"`

	// Value is the input data.
	Value interface{} `json:"value"`

	// InternalError is the underlying error thrown by the directive executor.
	InternalError error `json:"error"`
}

func (f *InvalidField) Error() string {
	return fmt.Sprintf("invalid field %q: %v", f.Field, f.InternalError)
}

func (f *InvalidField) Unwrap() error {
	return f.InternalError
}
