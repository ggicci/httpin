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

var builtResolvers sync.Map // map[reflect.Type]*owl.Resolver

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
		owl.WithNamespace(c.ns.encoders),
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
			removeDecoderDirective,             // backward compatibility, use "coder" instead
			removeCoderDirective,               // "coder" takes precedence over "decoder"
			removeCodecDirective,               // use "codec" instead
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

func removeDecoderDirective(r *owl.Resolver) error {
	return reserveCodecDirective(r, "decoder")
}

func removeCoderDirective(r *owl.Resolver) error {
	return reserveCodecDirective(r, "coder")
}

func removeCodecDirective(r *owl.Resolver) error {
	return reserveCodecDirective(r, "codec")
}

// reserveCodecDirective removes the directive from the resolver. name is "coder" or "decoder".
// The "decoder"/"coder"are two special directives which do nothing, but an indicator of
// overriding the decoder and encoder for a specific field.
func reserveCodecDirective(r *owl.Resolver, name string) error {
	d := r.RemoveDirective(name)
	if d == nil {
		return nil
	}
	if len(d.Argv) == 0 {
		return fmt.Errorf("directive %s: missing coder name", name)
	}

	if isFileType(r.Type) {
		return fmt.Errorf("directive %s: cannot be used on a file type field", name)
	}

	namedAdaptor := namedStringCodecAdaptors[d.Argv[0]]
	if namedAdaptor == nil {
		return fmt.Errorf("directive %s: %w: %q", name, ErrUnregisteredCoder, d.Argv[0])
	}

	r.Context = context.WithValue(r.Context, CtxCustomCoder, namedAdaptor)
	return nil
}

// ensureDirectiveExecutorsRegistered ensures all directives that defined in the
// resolver are registered in the executor registry.
func ensureDirectiveExecutorsRegistered(r *owl.Resolver) error {
	for _, d := range r.Directives {
		if decoderNamespace.LookupExecutor(d.Name) == nil {
			return fmt.Errorf("%w: %q", ErrUnregisteredDirective, d.Name)
		}
		// NOTE: don't need to check encoders namespace because a directive
		// will always be registered in both namespaces. See RegisterDirective().
	}
	return nil
}

type directiveOrderForEncoding []*owl.Directive

func (d directiveOrderForEncoding) Len() int {
	return len(d)
}

func (d directiveOrderForEncoding) Less(i, j int) bool {
	switch d[i].Name {
	case "default":
		return true // always the first one to run
	case "nonzero":
		return true // always the second one to run
	}
	return false
}

func (d directiveOrderForEncoding) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}
