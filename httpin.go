// Package httpin helps decoding an HTTP request to a custom struct by binding
// data with querystring (query params), HTTP headers, form data, JSON/XML
// payloads, URL path params, and file uploads (multipart/form-data).
package httpin

import (
	"context"
	"errors"
	"io"
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

// Option is a collection of options for creating a Core instance.
var Option coreOptions = coreOptions{
	WithErrorHandler:            core.WithErrorHandler,
	WithMaxMemory:               core.WithMaxMemory,
	WithNestedDirectivesEnabled: core.WithNestedDirectivesEnabled,
}

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
func New(inputStruct any, opts ...core.Option) (*core.Core, error) {
	return core.New(inputStruct, opts...)
}

// File is the builtin type of httpin to manupulate file uploads. On the server
// side, it is used to represent a file in a multipart/form-data request. On the
// client side, it is used to represent a file to be uploaded.
type File = core.File

// UploadFile is a helper function to create a File instance from a file path.
// It is useful when you want to upload a file from the local file system.
func UploadFile(path string) *File {
	return core.UploadFile(path)
}

// UploadStream is a helper function to create a File instance from a io.Reader. It
// is useful when you want to upload a file from a stream.
func UploadStream(r io.ReadCloser) *File {
	return core.UploadStream(r)
}

// DecodeTo decodes an HTTP request to the given input struct. The input must be
// a pointer (no matter how many levels) to a struct instance. For example:
//
//	input := &InputStruct{}
//	if err := DecodeTo(req, input); err != nil { ... }
//
// input is now populated with data from the request.
func DecodeTo(req *http.Request, input any, opts ...core.Option) error {
	co, err := New(internal.DereferencedType(input), opts...)
	if err != nil {
		return err
	}
	return co.DecodeTo(req, input)
}

// Decode decodes an HTTP request to a struct instance. The return value is a
// pointer to the input struct. For example:
//
//	if user, err := Decode[User](req); err != nil { ... }
//	// now user is a *User instance, which has been populated with data from the request.
func Decode[T any](req *http.Request, opts ...core.Option) (*T, error) {
	rt := internal.TypeOf[T]()
	if rt.Kind() != reflect.Struct {
		return nil, errors.New("generic type T must be a struct type")
	}
	co, err := New(rt, opts...)
	if err != nil {
		return nil, err
	}
	if v, err := co.Decode(req); err != nil {
		return nil, err
	} else {
		return v.(*T), nil
	}
}

// NewRequest wraps NewRequestWithContext using context.Background(), see NewRequestWithContext.
func NewRequest(method, url string, input any, opts ...core.Option) (*http.Request, error) {
	return NewRequestWithContext(context.Background(), method, url, input)
}

// NewRequestWithContext turns the given input struct into an HTTP request. The
// input struct with the "in" tags defines how to bind the data from the struct
// to the HTTP request. Use it as the replacement of http.NewRequest().
//
//	addUserPayload := &AddUserRequest{...}
//	addUserRequest, err := NewRequestWithContext(context.Background(), "GET", "http://example.com", addUserPayload)
//	http.DefaultClient.Do(addUserRequest)
func NewRequestWithContext(ctx context.Context, method, url string, input any, opts ...core.Option) (*http.Request, error) {
	co, err := New(input, opts...)
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

type coreOptions struct {
	// WithErrorHandler overrides the default error handler.
	WithErrorHandler func(core.ErrorHandler) core.Option

	// WithMaxMemory overrides the default maximum memory size (32MB) when reading
	// the request body. See https://pkg.go.dev/net/http#Request.ParseMultipartForm
	// for more details.
	WithMaxMemory func(int64) core.Option

	// WithNestedDirectivesEnabled enables/disables nested directives.
	WithNestedDirectivesEnabled func(bool) core.Option
}
