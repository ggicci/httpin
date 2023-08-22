// Package httpin helps decoding an HTTP request to a custom struct by binding
// data with querystring (query params), HTTP headers, form data, JSON/XML
// payloads, URL path params, and file uploads (multipart/form-data).
package httpin

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
)

type ContextKey int

const (
	minimumMaxMemory = 1 << 10  // 1KB
	defaultMaxMemory = 32 << 20 // 32 MB

	// Input is the key to get the input object from Request.Context() injected by httpin. e.g.
	//
	//     input := r.Context().Value(httpin.Input).(*InputStruct)
	Input ContextKey = iota

	RequestValue

	CustomDecoder

	// FieldSet is used by executors to tell whether a field has been set. When
	// multiple executors were applied to a field, if the field value were set
	// by a former executor, the latter executors MAY skip running by consulting
	// this context value.
	FieldSet

	StopRecursion
)

var (
	globalCustomErrorHandler ErrorHandler = defaultErrorHandler
)

type ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error)

// Decode decodes an HTTP request to a struct instance.
// e.g.
//
//	input := &InputStruct{}
//	if err := Decode(req, &input); err != nil { ... }
//
// input is now populated with data from the request.
func Decode(req *http.Request, input interface{}) error {
	originalType := reflect.TypeOf(input)
	if originalType.Kind() != reflect.Ptr {
		return fmt.Errorf("httpin: input must be a pointer")
	}
	core, err := New(originalType.Elem())
	if err != nil {
		return err
	}
	if value, err := core.Decode(req); err != nil {
		return err
	} else {
		if originalType.Elem().Kind() == reflect.Ptr {
			reflect.ValueOf(input).Elem().Set(reflect.ValueOf(value))
		} else {
			reflect.ValueOf(input).Elem().Set(reflect.ValueOf(value).Elem())
		}
		return nil
	}
}

// NewInput creates a "Middleware Constructor" for making a chain, which acts as
// a list of http.Handler constructors. We recommend using
// https://github.com/justinas/alice to chain your HTTP middleware functions and
// the app handler.
func NewInput(inputStruct interface{}, opts ...Option) func(http.Handler) http.Handler {
	core, err := New(inputStruct, opts...)
	if err != nil {
		panic(err)
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			// Here we read the request and decode it to fill our structure.
			// Once failed, the request should end here.
			input, err := core.Decode(r)
			if err != nil {
				core.getErrorHandler()(rw, r, err)
				return
			}

			// We put the `input` to the request's context, and it will pass to the next hop.
			ctx := context.WithValue(r.Context(), Input, input)
			next.ServeHTTP(rw, r.WithContext(ctx))
		})
	}
}

func ReplaceDefaultErrorHandler(custom ErrorHandler) {
	if custom == nil {
		panic(fmt.Errorf("httpin: %w", ErrNilErrorHandler))
	}
	globalCustomErrorHandler = custom
}

func defaultErrorHandler(rw http.ResponseWriter, r *http.Request, err error) {
	var invalidFieldError *InvalidFieldError
	if errors.As(err, &invalidFieldError) {
		rw.Header().Add("Content-Type", "application/json")
		rw.WriteHeader(http.StatusUnprocessableEntity) // status: 422
		json.NewEncoder(rw).Encode(invalidFieldError)
		return
	}

	http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError) // status: 500
}
