package core

import (
	"errors"
	"fmt"
	"strings"

	"github.com/ggicci/httpin/codec"
	"github.com/ggicci/owl"
)

var (
	ErrUnregisteredDirective = errors.New("unregistered directive")
	ErrUnregisteredCodec     = errors.New("unregistered codec")
	ErrFieldTypeMismatch     = codec.ErrFieldTypeMismatch
	ErrUnsupportedFieldType  = codec.ErrUnsupportedFieldType
	ErrUnsupportedType       = owl.ErrUnsupportedType

	ErrMissingCodecName      = errors.New("missing codec name")
	ErrTooManyNamedCodecs    = errors.New("too many named codecs")
	ErrIncompatibleDirective = errors.New("incompatible directive")
)

type InvalidFieldError struct {
	// err is the underlying error thrown by the directive executor.
	err error

	// Field is the name of the field.
	Field string `json:"field"`

	// Source is the directive which causes the error.
	// e.g. form, header, required, etc.
	Directive string `json:"directive"`

	// Key is the key to get the input data from the source.
	Key string `json:"key"`

	// Value is the input data.
	Value any `json:"value"`

	// ErrorMessage is the string representation of `internalError`.
	ErrorMessage string `json:"error"`
}

func (e *InvalidFieldError) Error() string {
	return fmt.Sprintf("invalid field %q: %v", e.Field, e.err)
}

func (e *InvalidFieldError) Unwrap() error {
	return e.err
}

func NewInvalidFieldError(err error) *InvalidFieldError {
	var (
		r  *owl.Resolver
		de *owl.DirectiveExecutionError
	)

	switch err := err.(type) {
	case *InvalidFieldError:
		return err
	case *owl.ResolveError:
		r = err.Resolver
		de = err.AsDirectiveExecutionError()
	case *owl.ScanError:
		r = err.Resolver
		de = err.AsDirectiveExecutionError()
	default:
		return &InvalidFieldError{
			err:          err,
			ErrorMessage: err.Error(),
		}
	}

	var fe *fieldError
	var inputKey string
	var inputValue any
	errors.As(err, &fe)
	if fe != nil {
		inputValue = fe.Value
		inputKey = fe.Key
	}

	return &InvalidFieldError{
		err:          err,
		Field:        r.Field.Name,
		Directive:    de.Name, // e.g. form, header, required, etc.
		Key:          inputKey,
		Value:        inputValue,
		ErrorMessage: err.Error(),
	}
}

type MultiInvalidFieldError []*InvalidFieldError

func (me MultiInvalidFieldError) Error() string {
	if len(me) == 1 {
		return me[0].Error()
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d invalid fields: ", len(me)))
	for i, e := range me {
		if i > 0 {
			sb.WriteString("; ")
		}
		sb.WriteString(e.Error())
	}
	return sb.String()
}

func (me MultiInvalidFieldError) Unwrap() []error {
	var errs []error
	for _, e := range me {
		errs = append(errs, e)
	}
	return errs
}

type fieldError struct {
	Key           string
	Value         any
	internalError error
}

func (e fieldError) Error() string {
	return e.internalError.Error()
}

func (e fieldError) Unwrap() error {
	return e.internalError
}
