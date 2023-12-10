package core

import (
	"errors"
	"fmt"

	"github.com/ggicci/owl"
)

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

func NewInvalidFieldError(err *owl.ResolveError) *InvalidFieldError {
	r := err.Resolver
	de := err.AsDirectiveExecutionError()

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
		Source:       de.Name, // e.g. form, header, required, etc.
		Key:          inputKey,
		Value:        inputValue,
		ErrorMessage: err.Error(),
	}
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
