package httpin

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/ggicci/owl"
)

var (
	ErrMissingField             = errors.New("missing required field")
	ErrUnsupporetedType         = errors.New("unsupported type")
	ErrUnregisteredExecutor     = errors.New("unregistered executor")
	ErrDuplicateTypeDecoder     = errors.New("duplicate type decoder")
	ErrDuplicateNamedDecoder    = errors.New("duplicate named decoder")
	ErrNilTypeDecoder           = errors.New("nil type decoder")
	ErrInvalidTypeDecoder       = errors.New("invalid type decoder")
	ErrDuplicateExecutor        = errors.New("duplicate executor")
	ErrReservedExecutorName     = errors.New("reserved executor name")
	ErrNilExecutor              = errors.New("nil executor")
	ErrUnknownBodyType          = errors.New("unknown body type")
	ErrDuplicateAnnotationField = errors.New("duplicate annotation field")
	ErrNilErrorHandler          = errors.New("nil error handler")
	ErrMaxMemoryTooSmall        = errors.New("max memory too small")
	ErrNilFile                  = errors.New("nil file")
	ErrDuplicateBodyDecoder     = errors.New("duplicate body decoder")
	ErrMissingDecoderName       = errors.New("missing decoder name")
	ErrDecoderNotFound          = errors.New("decoder not found")
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
	// err is the underlying error thrown by the directive executor.
	err error

	// Field is the name of the field.
	Field string `json:"field"`

	// Source is the directive which causes the error.
	// e.g. form, header, required, etc.
	Source string `json:"source"`

	// Key is the key to get the input data from the source.
	Key string `json:"key"`

	// Value is the input data.
	Value interface{} `json:"value"`

	// ErrorMessage is the string representation of `internalError`.
	ErrorMessage string `json:"error"`
}

func (e *InvalidFieldError) Error() string {
	return fmt.Sprintf("invalid field %q: %v", e.Field, e.err)
}

func (e *InvalidFieldError) Unwrap() error {
	return e.err
}

func NewInvalidFieldError(err *owl.ResolveError) *InvalidFieldError {
	r := err.Resolver
	de := err.AsDirectiveExecutionError()

	var fe *fieldError
	var inputKey string
	var inputValue interface{}
	errors.As(err, &fe)
	if fe != nil {
		inputValue = fe.Value
		inputKey = fe.Key
	}

	return &InvalidFieldError{
		err:          err,
		Field:        r.Field.Name,
		Source:       de.Name, // e.g. form, header, required, etc.
		Key:          inputKey,
		Value:        inputValue,
		ErrorMessage: err.Error(),
	}
}

type fieldError struct {
	Key           string
	Value         interface{}
	internalError error
}

func (e fieldError) Error() string {
	return e.internalError.Error()
}

func (e fieldError) Unwrap() error {
	return e.internalError
}
