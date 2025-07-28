package core

import (
	"context"
	"fmt"
	"reflect"
	"sync"

	"github.com/ggicci/httpin/codec"
	"github.com/ggicci/owl"
)

var defaultNS = NewNamespace()

// Namespace represents an isolated environment that holds all registered directives, codecs, and codec adaptors.
// A Core instance created from this namespace operates solely based on the information contained within it.
type Namespace struct {
	*codec.Namespace            // for registering codecs
	*DirectiveRegistry          // for registering directive executors
	builtResolvers     sync.Map // map[reflect.Type]*owl.Resolver

	fileTypes                map[reflect.Type]struct{}
	namedStringCodecAdaptors map[string]*NamedStringCodecAdaptor
}

func NewNamespace() *Namespace {
	ns := &Namespace{
		Namespace:         codec.NewNamespace(),
		DirectiveRegistry: NewDirectiveRegistry(),

		fileTypes:                make(map[reflect.Type]struct{}),
		namedStringCodecAdaptors: make(map[string]*NamedStringCodecAdaptor),
	}
	return ns
}

func (ns *Namespace) New(inputStruct any, opts ...Option) (*Core, error) {
	resolver, err := ns.buildResolver(inputStruct)
	if err != nil {
		return nil, err
	}

	core := &Core{
		ns:       ns,
		resolver: resolver,
	}

	// Apply default options and user custom options to the
	var allOptions []Option
	defaultOptions := []Option{
		WithMaxMemory(defaultMaxMemory),
		WithNestedDirectivesEnabled(globalNestedDirectivesEnabled),
	}
	allOptions = append(allOptions, defaultOptions...)
	allOptions = append(allOptions, opts...)

	for _, opt := range allOptions {
		if err := opt(core); err != nil {
			return nil, fmt.Errorf("invalid option: %w", err)
		}
	}

	return core, nil
}

// RegisterCodec overrides the builtin codec that registered by default for the given
// type typ. For example, a codec for bool type is registered by default which can
// parse strings like "true", "false", "1", "0", and so on to a bool value.
// If you would like to override the default behaviour for a particular type.
// You should register a custom codec here. For instance:
//
//	func init() {
//		ns.RegisterCodec(strconvx.ToAnyAdaptor(func(b *bool) (core.StringCodec, error) {
//			return (*YesNo)(b), nil
//		}))
//	}
//
//	type YesNo bool
//
//	func (yn YesNo) String() string {
//		if yn {
//			return "yes"
//		}
//		return "no"
//	}
//
//	func (yn *YesNo) FromString(s string) error {
//		switch s {
//		case "yes":
//			*yn = true
//		case "no":
//			*yn = false
//		default:
//			return fmt.Errorf("invalid YesNo value: %q", s)
//		}
//		return nil
//	}
func (ns *Namespace) RegisterCodec(typ reflect.Type, adaptor StringCodecAdaptor) {
	ns.Adapt(typ, adaptor)
}

// TODO(ggicci): add comment
func (ns *Namespace) RegisterNamedCodec(name string, typ reflect.Type, adaptor StringCodecAdaptor) {
	ns.namedStringCodecAdaptors[name] = &NamedStringCodecAdaptor{
		Name:     name,
		BaseType: typ,
		Adaptor:  adaptor,
	}
}

// TODO(ggicci): add comment
func (ns *Namespace) RegisterFileCodec(typ reflect.Type) {
	ns.fileTypes[typ] = struct{}{}
}

// IsFileType checks if the given typ is registered as a "File".
func (ns *Namespace) IsFileType(typ reflect.Type) bool {
	baseType, _ := codec.BaseTypeOf(typ)
	_, ok := ns.fileTypes[baseType]
	return ok
}

// buildResolver builds a resolver for the inputStruct. It will run normalizations
// on the resolver and cache it.
func (ns *Namespace) buildResolver(inputStruct any) (*owl.Resolver, error) {
	resolver, err := owl.New(inputStruct)
	if err != nil {
		return nil, err
	}

	// Returns the cached resolver if it's already built.
	if cached, ok := ns.builtResolvers.Load(resolver.Type); ok {
		return cached.(*owl.Resolver), nil
	}

	// Normalize the resolver before caching it.
	if err := ns.normalizeResolver(resolver); err != nil {
		return nil, err
	}

	// Cache the resolver.
	ns.builtResolvers.Store(resolver.Type, resolver)
	return resolver, nil
}

// normalizeResolver normalizes the resolvers by running a series of
// normalizations on every field resolver.
func (ns *Namespace) normalizeResolver(r *owl.Resolver) error {
	normalize := func(r *owl.Resolver) error {
		for _, fn := range []func(*owl.Resolver) error{
			ns.removeDecoderDirective,             // backward compatibility, use "coder" instead
			ns.removeCoderDirective,               // "coder" takes precedence over "decoder"
			ns.removeCodecDirective,               // use "codec" instead
			ns.ensureDirectiveExecutorsRegistered, // always the last one
		} {
			if err := fn(r); err != nil {
				return err
			}
		}
		return nil
	}

	return r.Iterate(normalize)
}

// ensureDirectiveExecutorsRegistered ensures all directives that defined in the
// resolver are registered in the executor registry.
func (ns *Namespace) ensureDirectiveExecutorsRegistered(r *owl.Resolver) error {
	for _, d := range r.Directives {
		if ns.decoders.LookupExecutor(d.Name) == nil {
			return fmt.Errorf("%w: %q", ErrUnregisteredDirective, d.Name)
		}
		// NOTE: don't need to check encoders namespace because a directive
		// will always be registered in both namespaces. See RegisterDirective().
	}
	return nil
}

func (ns *Namespace) removeDecoderDirective(r *owl.Resolver) error {
	return ns.reserveCodecDirective(r, "decoder")
}

func (ns *Namespace) removeCoderDirective(r *owl.Resolver) error {
	return ns.reserveCodecDirective(r, "coder")
}

func (ns *Namespace) removeCodecDirective(r *owl.Resolver) error {
	return ns.reserveCodecDirective(r, "codec")
}

// reserveCodecDirective removes the directive from the resolver. The given name
// can be one of "codec", "coder", "decoder". All these are special directives
// which do nothing, but as an indicator of overriding the decoder and encoder
// for a specific field.
func (ns *Namespace) reserveCodecDirective(r *owl.Resolver, name string) error {
	d := r.RemoveDirective(name)
	if d == nil {
		return nil
	}
	if len(d.Argv) == 0 {
		return fmt.Errorf("directive %q: %w", name, ErrMissingCodecName)
	}
	if len(d.Argv) > 1 {
		return fmt.Errorf("directive %q: %w, expected one, got %d",
			name, ErrTooManyNamedCodecs, len(d.Argv))
	}
	if ns.IsFileType(r.Type) {
		return fmt.Errorf("directive %q: %w, cannot be used on a field of file type ", name, ErrIncompatibleDirective)
	}

	namedAdaptor := ns.namedStringCodecAdaptors[d.Argv[0]]
	if namedAdaptor == nil {
		return fmt.Errorf("directive %q: %w: %q", name, ErrUnregisteredCodec, d.Argv[0])
	}

	r.Context = context.WithValue(r.Context, CtxCustomCodec, namedAdaptor)
	return nil
}

// unregisterNamedCodec removes the codec registered as the given name.
// It's only used by the unit tests of this package.
func (ns *Namespace) unregisterNamedCodec(name string) {
	delete(ns.namedStringCodecAdaptors, name)
}

// unregisterFileCodec deletes the codec for the given typ.
// It's only used by the unit tests of this package.
func (ns *Namespace) unregisterFileCodec(typ reflect.Type) {
	delete(ns.fileTypes, typ)
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
