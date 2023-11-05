// Package httpin helps decoding an HTTP request to a custom struct by binding
// data with querystring (query params), HTTP headers, form data, JSON/XML
// payloads, URL path params, and file uploads (multipart/form-data).
package httpin

import (
	"context"
	"fmt"
	"net/http"
	"reflect"

	"github.com/ggicci/httpin/internal"
)

type InvalidFieldError = internal.InvalidFieldError

type ContextKey int

const (
	// Input is the key to get the input object from Request.Context() injected by httpin. e.g.
	//
	//     input := r.Context().Value(httpin.Input).(*InputStruct)
	Input ContextKey = iota
)

// Decode decodes an HTTP request to the given input struct. The input must be a
// pointer to a struct instance. For example:
//
//	input := &InputStruct{}
//	if err := Decode(req, &input); err != nil { ... }
//
// input is now populated with data from the request.
func Decode(req *http.Request, input any) error {
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

// Encode is an alias of NewRequest.
func Encode(method, url string, input any) (*http.Request, error) {
	return NewRequest(method, url, input)
}

// NewRequest wraps NewRequestWithContext using context.Background.
func NewRequest(method, url string, input any) (*http.Request, error) {
	return NewRequestWithContext(context.Background(), method, url, input)
}

// NewRequestWithContext returns a new http.Request given a method, url and an
// input struct instance. The fields of the input struct will be encoded to the
// request by resolving the "in" tags and executing the directives.
func NewRequestWithContext(ctx context.Context, method, url string, input any) (*http.Request, error) {
	core, err := New(input)
	if err != nil {
		return nil, err
	}
	return core.NewRequestWithContext(ctx, method, url, input)
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
func NewInput(inputStruct any, opts ...Option) func(http.Handler) http.Handler {
	core, err := New(inputStruct, opts...)
	internal.PanicOnError(err)

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
