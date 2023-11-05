package httpin

import (
	"context"
	"errors"
	"fmt"
	"mime"
	"net/http"
	"strings"
	"sync"

	"github.com/ggicci/httpin/internal"
	"github.com/ggicci/owl"
)

const (
	minimumMaxMemory = int64(1 << 10)  // 1KB
	defaultMaxMemory = int64(32 << 20) // 32 MB
)

type RequestBuilder = internal.RequestBuilder

var builtResolvers sync.Map // map[reflect.Type]*owl.Resolver

// core is the core of httpin. It holds the resolver of a specific struct type.
// Who is responsible for decoding an HTTP request to an instance of such struct
// type.
type core struct {
	resolver *owl.Resolver

	errorHandler errorHandler
	maxMemory    int64 // in bytes
}

// New creates a new Core instance, which holds the resolver of the inputStruct.
//   - Use Core.Decode() to decode an HTTP request to an instance of the inputStruct.
//   - Use Core.NewRequest() to encode an instance of the inputStruct to an HTTP request.
//   - Use NewInput() to create an HTTP middleware.
func New(inputStruct any, opts ...coreOption) (*core, error) {
	resolver, err := buildResolver(inputStruct)
	if err != nil {
		return nil, err
	}

	core := &core{
		resolver: resolver,
	}

	// Apply default options and user custom options to the core.
	var allOptions []coreOption
	defaultOptions := []coreOption{
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

// Decode decodes an HTTP request to a struct instance.
// The return value is a pointer to the input struct.
// For example:
//
//	New(&Input{}).Decode(req) -> *Input
//	New(Input{}).Decode(req) -> *Input
func (c *core) Decode(req *http.Request) (any, error) {
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

	rv, err := c.resolver.Resolve(
		owl.WithNamespace(decoderNamespace),
		owl.WithValue(internal.CtxRequest, req),
	)
	if err != nil {
		return nil, internal.NewInvalidFieldError(err.(*owl.ResolveError))
	}
	return rv.Interface(), nil
}

// NewRequest wraps NewRequestWithContext using context.Background.
func (c *core) NewRequest(method string, url string, input any) (*http.Request, error) {
	return c.NewRequestWithContext(context.Background(), method, url, input)
}

// NewRequestWithContext returns a new http.Request given a method, url and an
// input struct instance. Note that the Core instance is bound to a specific
// type of struct. Which means when the given input is not the type of the
// struct that the Core instance holds, error of type mismatch will be returned.
// In order to avoid this error, you can always use httpin.NewRequest() function
// instead. Which will create a Core instance for you when needed. There's no
// performance penalty for doing so. Because there's a cache layer for all the
// Core instances.
func (c *core) NewRequestWithContext(ctx context.Context, method string, url string, input any) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, err
	}

	rb := &RequestBuilder{}
	if err = c.resolver.Scan(
		input,
		owl.WithNamespace(encoderNamespace),
		owl.WithValue(internal.CtxRequestBuilder, rb),
	); err != nil {
		return nil, err
	}

	// Populate the request with the encoded values.
	if err := rb.Populate(req); err != nil {
		return nil, err
	}

	return req, nil
}

// buildResolver builds a resolver for the inputStruct. It will run normalizations
// on the resolver and cache it.
func buildResolver(inputStruct any) (*owl.Resolver, error) {
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
			reserveEncoderDirective,
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
		return errors.New("missing decoder name")
	}
	decoder := internal.DefaultRegistry.GetNamedDecoder(d.Argv[0])
	if decoder == nil {
		return fmt.Errorf("unregistered decoder: %q", d.Argv[0])
	}
	if internal.DefaultRegistry.IsFileType(r.Type) {
		return errors.New("cannot use decoder directive on a file type field")
	}
	r.Context = context.WithValue(r.Context, internal.CtxCustomDecoder, decoder)
	return nil
}

func reserveEncoderDirective(r *owl.Resolver) error {
	d := r.RemoveDirective("encoder")
	if d == nil {
		return nil
	}
	if len(d.Argv) == 0 {
		return errors.New("missing encoder name")
	}
	encoder := internal.DefaultRegistry.GetNamedEncoder(d.Argv[0])
	if encoder == nil {
		return fmt.Errorf("unregistered encoder: %q", d.Argv[0])
	}
	if internal.DefaultRegistry.IsFileType(r.Type) {
		return errors.New("cannot use encoder directive on a file type field")
	}
	r.Context = context.WithValue(r.Context, internal.CtxCustomEncoder, encoder)
	return nil
}

// normalizeBodyDirective normalizes the body directive of the resolver.
// If no body format specified, the default type is "json".
func normalizeBodyDirective(r *owl.Resolver) error {
	dir := r.GetDirective("body")
	if dir == nil || dir.Name != "body" {
		return nil
	}
	if len(dir.Argv) == 0 {
		dir.Argv = []string{"json"} // use json as default when no body format specified
	}
	dir.Argv[0] = strings.ToLower(dir.Argv[0])

	var bodyFormat = dir.Argv[0]
	if internal.DefaultRegistry.GetBodyDecoder(bodyFormat) == nil {
		return fmt.Errorf("unknown body format: %q", bodyFormat)
	}
	return nil
}

// ensureDirectiveExecutorsRegistered ensures all directives that defined in the
// resolver are registered in the executor registry.
func ensureDirectiveExecutorsRegistered(r *owl.Resolver) error {
	for _, d := range r.Directives {
		if decoderNamespace.LookupExecutor(d.Name) == nil {
			return fmt.Errorf("unregistered directive: %q (decoder)", d.Name)
		}
		if encoderNamespace.LookupExecutor(d.Name) == nil {
			return fmt.Errorf("unregistered directive: %q (encoder)", d.Name)
		}
	}
	return nil
}

// getErrorHandler returns the error handler of the core if set, or the global
// custom error handler.
func (c *core) getErrorHandler() errorHandler {
	if c.errorHandler != nil {
		return c.errorHandler
	}

	return globalCustomErrorHandler
}
