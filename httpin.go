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
	// Input is the key to get the input object from Request.Context() injected by httpin. e.g.
	//
	//     input := r.Context().Value(httpin.Input).(*InputStruct)
	Input ContextKey = iota

	// RequestValue is the key to get the HTTP request value (of *http.Request)
	// from DirectiveRuntime.Context. The HTTP request value is injected by
	// httpin to the context of DirectiveRuntime before executing the directive.
	// See Core.Decode() for more details.
	RequestValue

	// CustomDecoder is the key to get the custom decoder for a field from
	// Resolver.Context. Which is specified by the "decoder" directive.
	// During resolver building phase, the "decoder" directive will be removed
	// from the resolver, and the targeted decoder by name will be put into
	// Resolver.Context with this key. e.g.
	//
	//    type GreetInput struct {
	//        Message string `httpin:"decoder=custom"`
	//    }
	// For the above example, the decoder named "custom" will be put into the
	// resolver of Message field with this key.
	CustomDecoder

	// FieldSet is used by executors to tell whether a field has been set. When
	// multiple executors were applied to a field, if the field value were set
	// by a former executor, the latter executors MAY skip running by consulting
	// this context value.
	FieldSet
)

var (
	globalCustomErrorHandler ErrorHandler = defaultErrorHandler
)

// ErrorHandler is the type of custom error handler. The error handler is used
// by the http.Handler that created by NewInput() to handle errors during
// decoding the HTTP request.
type ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error)

// Decode decodes an HTTP request to the given input struct. The input must be a
// pointer to a struct instance. For example:
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

func Encode(method, url string, input interface{}) (*http.Request, error) {
	core, err := New(input)
	if err != nil {
		return nil, err
	}
	return core.Encode(method, url, input)
}

// NewInput creates a "Middleware". A middleware is a function that takes a
// http.Handler and returns another http.Handler.
//
// The middleware created by NewInput is to add the decoding function to an
// existing http.Handler. This functionality will decode the HTTP request and
// put the decoded struct instance to the request's context. So that the next
// hop can get the decoded struct instance from the request's context.
//
// We recommend using https://github.com/justinas/alice to chain your
// middlewares. If you're using some popular web frameworks, they may have
// already provided a middleware chaining mechanism.
func NewInput(inputStruct interface{}, opts ...Option) func(http.Handler) http.Handler {
	core, err := New(inputStruct, opts...)
	if err != nil {
		panic(fmt.Errorf("httpin: %w", err))
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

// ReplaceDefaultErrorHandler replaces the default error handler with the given
// custom error handler. The default error handler will be used in the http.Handler
// that decoreated by the middleware created by NewInput().
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
