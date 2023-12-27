// Package httpin helps decoding an HTTP request to a custom struct by binding
// data with querystring (query params), HTTP headers, form data, JSON/XML
// payloads, URL path params, and file uploads (multipart/form-data).
package httpin

import (
	"context"
	"fmt"
	"net/http"
	"reflect"

	"github.com/ggicci/httpin/core"
	"github.com/ggicci/httpin/internal"
)

type contextKey int

const (
	// Input is the key to get the input object from Request.Context() injected by httpin. e.g.
	//
	//     input := r.Context().Value(httpin.Input).(*InputStruct)
	Input contextKey = iota
)

// New calls core.New to create a new Core instance. Which is responsible for both:
//
//   - decoding an HTTP request to an instance of the inputStruct;
//   - and encoding an instance of the inputStruct to an HTTP request.
//
// Note that the Core instance is bound to the given specific type, it will not
// work for other types. If you want to decode/encode other types, you need to
// create another Core instance. Or directly use the following functions, which are
// just shortcuts of Core's methods, so you don't need to create a Core instance:
//   - httpin.Decode(): decode an HTTP request to an instance of the inputStruct.
//   - httpin.NewRequest() to encode an instance of the inputStruct to an HTTP request.
//
// For best practice, we would recommend using httpin.NewInput() to create an
// HTTP middleware for a specific input type. The middleware can be bound to an
// API, chained with other middlewares, and also reused in other APIs. You even
// don't need to call the Deocde() method explicitly, the middleware will do it
// for you and put the decoded instance to the request's context.
var New = core.New

// WithMaxMemory overrides the default maximum memory size (32MB) when reading
// the request body. See https://pkg.go.dev/net/http#Request.ParseMultipartForm
// for more details.
var WithMaxMemory = core.WithMaxMemory

// WithErrorHandler overrides the default error handler.
// If you want to override the default error handler globally, you can use core.RegisterErrorHandler.
var WithErrorHandler = core.WithErrorHandler

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
	co, err := New(originalType.Elem())
	if err != nil {
		return err
	}
	if value, err := co.Decode(req); err != nil {
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

// NewRequest wraps NewRequestWithContext using context.Background.
func NewRequest(method, url string, input any) (*http.Request, error) {
	return NewRequestWithContext(context.Background(), method, url, input)
}

// NewRequestWithContext returns a new http.Request given a method, url and an
// input struct instance. The fields of the input struct will be encoded to the
// request by resolving the "in" tags and executing the directives.
func NewRequestWithContext(ctx context.Context, method, url string, input any) (*http.Request, error) {
	co, err := New(input)
	if err != nil {
		return nil, err
	}
	return co.NewRequestWithContext(ctx, method, url, input)
}

// NewInput creates an HTTP middleware handler. Which is a function that takes
// in an http.Handler and returns another http.Handler.
//
// The middleware created by NewInput is to add the decoding function to an
// existing http.Handler. This functionality will decode the HTTP request and
// put the decoded struct instance to the request's context. So that the next
// hop can get the decoded struct instance from the request's context.
//
// We recommend using https://github.com/justinas/alice to chain your
// middlewares. If you're using some popular web frameworks, they may have
// already provided a middleware chaining mechanism.
//
// For example:
//
//	type ListUsersRequest struct {
//		Page    int `in:"query=page,page_index,index"`
//		PerPage int `in:"query=per_page,page_size"`
//	}
//
//	func ListUsersHandler(rw http.ResponseWriter, r *http.Request) {
//		input := r.Context().Value(httpin.Input).(*ListUsersRequest)
//		// ...
//	}
//
//	func init() {
//		http.Handle("/users", alice.New(httpin.NewInput(&ListUsersRequest{})).ThenFunc(ListUsersHandler))
//	}
func NewInput(inputStruct any, opts ...core.Option) func(http.Handler) http.Handler {
	co, err := New(inputStruct, opts...)
	internal.PanicOnError(err)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			// Here we read the request and decode it to fill our structure.
			// Once failed, the request should end here.
			input, err := co.Decode(r)
			if err != nil {
				co.GetErrorHandler()(rw, r, err)
				return
			}

			// We put the `input` to the request's context, and it will pass to the next hop.
			ctx := context.WithValue(r.Context(), Input, input)
			next.ServeHTTP(rw, r.WithContext(ctx))
		})
	}
}
