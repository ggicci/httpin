package httpin

import (
	"errors"
	"fmt"
	"reflect"
)

var (
	ErrMissingField             = errors.New("missing required field")
	ErrUnsupporetedType         = errors.New("unsupported type")
	ErrUnregisteredExecutor     = errors.New("unregistered executor")
	ErrDuplicateTypeDecoder     = errors.New("duplicate type decoder")
	ErrNilTypeDecoder           = errors.New("nil type decoder")
	ErrDuplicateExecutor        = errors.New("duplicate executor")
	ErrNilExecutor              = errors.New("nil executor")
	ErrUnknownBodyType          = errors.New("unknown body type")
	ErrDuplicateAnnotationField = errors.New("duplicate annotation field")
	ErrNilErrorHandler          = errors.New("nil error handler")
)

type UnsupportedTypeError struct {
	Type reflect.Type
}

func (e UnsupportedTypeError) Error() string {
	return fmt.Sprintf("unsupported type: %q", e.Type)
}

func (e UnsupportedTypeError) Unwrap() error {
	return ErrUnsupporetedType
}

type InvalidFieldError struct {
	// Field is the name of the field.
	Field string `json:"field"`

	// Source is the directive which causes the error.
	// e.g. form, header, required, etc.
	Source string `json:"source"`

	// Value is the input data.
	Value interface{} `json:"value"`

	// internalError is the underlying error thrown by the directive executor.
	internalError error

	// ErrorMessage is the string representation of `internalError`.
	ErrorMessage string `json:"error"`

	// directives is the list of directives bound to the field.
	Directives []*Directive `json:"-"`
}

func (f *InvalidFieldError) Error() string {
	return fmt.Sprintf("invalid field %q: %v", f.Field, f.internalError)
}

func (f *InvalidFieldError) Unwrap() error {
	return f.internalError
}

type fieldError struct {
	Key           string
	Value         interface{}
	internalError error
}

func (e fieldError) Error() string {
	return e.internalError.Error()
}
