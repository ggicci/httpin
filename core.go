package httpin

import (
	"context"
	"fmt"
	"mime"
	"net/http"
	"sync"

	"github.com/ggicci/owl"
)

const (
	minimumMaxMemory = int64(1 << 10)  // 1KB
	defaultMaxMemory = int64(32 << 20) // 32 MB
)

var builtResolvers sync.Map // map[reflect.Type]*owl.Resolver

// Core is the core of httpin. It holds the resolver of a specific struct type.
// Who is responsible for decoding an HTTP request to an instance of such struct
// type.
type Core struct {
	resolver *owl.Resolver

	errorHandler ErrorHandler
	maxMemory    int64 // in bytes
}

// New creates a new Core instance, which holds the resolver of the inputStruct.
//   - Use Core.Decode() to decode an HTTP request to an instance of the inputStruct.
//   - Use NewInput() to create an HTTP middleware.
func New(inputStruct interface{}, opts ...Option) (*Core, error) {
	resolver, err := buildResolver(inputStruct)
	if err != nil {
		return nil, err
	}

	core := &Core{
		resolver: resolver,
	}

	// Apply default options and user custom options to the core.
	var allOptions []Option
	defaultOptions := []Option{
		WithMaxMemory(defaultMaxMemory),
	}
	allOptions = append(allOptions, defaultOptions...)
	allOptions = append(allOptions, opts...)

	for _, opt := range allOptions {
		if err := opt(core); err != nil {
			return nil, fmt.Errorf("httpin: invalid option: %w", err)
		}
	}

	return core, nil
}

// buildResolver builds a resolver for the inputStruct. It will run normalizations
// on the resolver and cache it.
func buildResolver(inputStruct interface{}) (*owl.Resolver, error) {
	resolver, err := owl.New(inputStruct)
	if err != nil {
		return nil, err
	}

	// Returns the cached resolver if it's already built.
	if cached, ok := builtResolvers.Load(resolver.Type); ok {
		return cached.(*owl.Resolver), nil
	}

	// Normalize the resolver before caching it.
	if err := normalizeResolver(resolver); err != nil {
		return nil, err
	}

	// Cache the resolver.
	builtResolvers.Store(resolver.Type, resolver)
	return resolver, nil
}

// normalizeResolver normalizes the resolvers by running a series of
// normalizations on every field resolver.
func normalizeResolver(r *owl.Resolver) error {
	normalize := func(r *owl.Resolver) error {
		for _, fn := range []func(*owl.Resolver) error{
			reserveDecoderDirective,
			normalizeBodyDirective,
			ensureDirectiveExecutorsRegistered, // always the last one
		} {
			if err := fn(r); err != nil {
				return err
			}
		}
		return nil
	}

	return r.Iterate(normalize)
}

// reserveDecoderDirective removes the "decoder" directive from the resolver.
// The "decoder" is a special directive which does nothing, but an indicator of
// overriding the decoder for a specific field.
func reserveDecoderDirective(r *owl.Resolver) error {
	d := r.RemoveDirective("decoder")
	if d == nil {
		return nil
	}
	if len(d.Argv) == 0 {
		return ErrMissingDecoderName
	}

	decoder := decoderByName(d.Argv[0])
	if decoder == nil {
		return ErrDecoderNotFound
	}
	r.Context = context.WithValue(r.Context, CustomDecoder, decoder)
	return nil
}

// ensureDirectiveExecutorsRegistered ensures all directives that defined in the
// resolver are registered in the executor registry.
func ensureDirectiveExecutorsRegistered(r *owl.Resolver) error {
	for _, d := range r.Directives {
		if owl.LookupExecutor(d.Name) == nil {
			return fmt.Errorf("%w: %q", ErrUnregisteredExecutor, d.Name)
		}
	}
	return nil
}

// Decode decodes an HTTP request to a struct instance.
// The return value is a pointer to the input struct.
// For example:
//
//	New(&Input{}).Decode(req) -> *Input
//	New(Input{}).Decode(req) -> *Input
func (c *Core) Decode(req *http.Request) (interface{}, error) {
	var err error
	ct, _, _ := mime.ParseMediaType(req.Header.Get("Content-Type"))
	if ct == "multipart/form-data" {
		err = req.ParseMultipartForm(c.maxMemory)
	} else {
		err = req.ParseForm()
	}
	if err != nil {
		return nil, err
	}

	rv, err := c.resolver.Resolve(owl.WithValue(RequestValue, req))
	if err != nil {
		return nil, NewInvalidFieldError(err.(*owl.ResolveError))
	}
	return rv.Interface(), nil
}

// getErrorHandler returns the error handler of the core if set, or the global
// custom error handler.
func (c *Core) getErrorHandler() ErrorHandler {
	if c.errorHandler != nil {
		return c.errorHandler
	}

	return globalCustomErrorHandler
}
