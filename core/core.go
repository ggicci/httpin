package core

import (
	"context"
	"errors"
	"fmt"
	"mime"
	"net/http"
	"reflect"
	"sort"
	"sync"

	"github.com/ggicci/owl"
)

// Core is the Core of httpin. It holds the resolver of a specific struct type.
// Who is responsible for decoding an HTTP request to an instance of such struct
// type.
type Core struct {
	ns                     *Namespace
	resolver               *owl.Resolver // for decoding
	scanResolver           *owl.Resolver // for encoding
	errorHandler           ErrorHandler
	maxMemory              int64 // in bytes
	enableNestedDirectives bool
	resolverMu             sync.RWMutex
}

// New creates a new Core instance for the given intpuStruct. It will build a resolver
// for the inputStruct and apply the given options to the Core instance. The Core instance
// is responsible for both:
//
//   - decoding an HTTP request to an instance of the inputStruct;
//   - encoding an instance of the inputStruct to an HTTP request.
func New(inputStruct any, opts ...Option) (*Core, error) {
	return defaultNS.New(inputStruct, opts...)
}

// Decode decodes an HTTP request to an instance of the input struct and returns
// its pointer. For example:
//
//	New(Input{}).Decode(req) -> *Input
func (c *Core) Decode(req *http.Request) (any, error) {
	// Create the input struct instance. Used to be created by owl.Resolve().
	value := reflect.New(c.resolver.Type).Interface()
	if err := c.DecodeTo(req, value); err != nil {
		return nil, err
	} else {
		return value, nil
	}
}

// DecodeTo decodes an HTTP request to the given value. The value must be a pointer
// to the struct instance of the type that the Core instance holds.
func (c *Core) DecodeTo(req *http.Request, value any) (err error) {
	if err = c.parseRequestForm(req); err != nil {
		return fmt.Errorf("failed to parse request form: %w", err)
	}

	err = c.resolver.ResolveTo(
		value,
		owl.WithValue(CtxNamespace, c.ns),
		owl.WithNamespace(c.ns.decoders),
		owl.WithValue(CtxRequest, req),
		owl.WithNestedDirectivesEnabled(c.enableNestedDirectives),
	)
	if err != nil && !errors.Is(err, owl.ErrInvalidResolveTarget) {
		return NewInvalidFieldError(err)
	}
	return err
}

// NewRequest wraps NewRequestWithContext using context.Background(), see
// NewRequestWithContext.
func (c *Core) NewRequest(method string, url string, input any) (*http.Request, error) {
	return c.NewRequestWithContext(context.Background(), method, url, input)
}

// NewRequestWithContext turns the given input struct into an HTTP request. Note
// that the Core instance is bound to a specific type of struct. Which means
// when the given input is not the type of the struct that the Core instance
// holds, error of type mismatch will be returned. In order to avoid this error,
// you can always use httpin.NewRequest() instead. Which will create a Core
// instance for you on demand. There's no performance penalty for doing so.
// Because there's a cache layer for all the Core instances.
func (c *Core) NewRequestWithContext(ctx context.Context, method string, url string, input any) (*http.Request, error) {
	c.prepareScanResolver()
	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, err
	}

	rb := NewRequestBuilder(ctx)

	// NOTE(ggicci): the error returned a joined error by using errors.Join.
	if err = c.scanResolver.Scan(
		input,
		owl.WithValue(CtxNamespace, c.ns),
		owl.WithNamespace(c.ns.encoders), // NOTE: this is owl's namespace
		owl.WithValue(CtxRequestBuilder, rb),
		owl.WithNestedDirectivesEnabled(c.enableNestedDirectives),
	); err != nil {
		// err is a list of *owl.ScanError that joined by errors.Join.
		if errs, ok := err.(interface{ Unwrap() []error }); ok {
			var invalidFieldErrors MultiInvalidFieldError
			for _, err := range errs.Unwrap() {
				invalidFieldErrors = append(invalidFieldErrors, NewInvalidFieldError(err))
			}
			return nil, invalidFieldErrors
		} else {
			return nil, err // should never happen, just in case
		}
	}

	// Populate the request with the encoded values.
	if err := rb.Populate(req); err != nil {
		return nil, fmt.Errorf("failed to populate request: %w", err)
	}

	return req, nil
}

// GetErrorHandler returns the error handler of the core if set, or the global
// custom error handler.
func (c *Core) GetErrorHandler() ErrorHandler {
	if c.errorHandler != nil {
		return c.errorHandler
	}

	return globalCustomErrorHandler
}

func (c *Core) prepareScanResolver() {
	c.resolverMu.RLock()
	if c.scanResolver == nil {
		c.resolverMu.RUnlock()
		c.resolverMu.Lock()
		defer c.resolverMu.Unlock()

		if c.scanResolver == nil {
			c.scanResolver = c.resolver.Copy()

			// Reorder the directives to make sure the "default" and "nonzero" directives work properly.
			c.scanResolver.Iterate(func(r *owl.Resolver) error {
				sort.Sort(directiveOrderForEncoding(r.Directives))
				return nil
			})
		}
	} else {
		c.resolverMu.RUnlock()
	}
}

func (c *Core) parseRequestForm(req *http.Request) (err error) {
	ct, _, _ := mime.ParseMediaType(req.Header.Get("Content-Type"))
	if ct == "multipart/form-data" {
		err = req.ParseMultipartForm(c.maxMemory)
	} else {
		err = req.ParseForm()
	}
	return
}
