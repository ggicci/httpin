package core

import (
	"errors"
	"fmt"

	"github.com/ggicci/httpin/internal"
	"github.com/ggicci/owl"
)

type DirectiveExecutor interface {
	// Encode encodes the field of the input struct to the HTTP request.
	Encode(*DirectiveRuntime) error

	// Decode decodes the field of the input struct from the HTTP request.
	Decode(*DirectiveRuntime) error
}

// FIXME(ggicci): remove the following decoderNamespace & encoderNamespace.
var (
	// decoderNamespace is the namespace for registering directive executors that are
	// used to decode the http request to input struct.
	decoderNamespace = owl.NewNamespace()

	// encoderNamespace is the namespace for registering directive executors that are
	// used to encode the input struct to http request.
	encoderNamespace = owl.NewNamespace()
)

// reservedExecutorNames are the names that cannot be used to register user defined directives
var reservedExecutorNames = []string{"codec", "decoder", "coder"}

// DirectiveNamespace holds the namespaces for directive executors. It is used to
// register directive executors that are used to decode and encode the HTTP request
// to/from the input struct.
type DirectiveNamespace struct {
	// decoders is the namespace for registering directive executors that are
	// used to decode the HTTP request to input struct.
	decoders *owl.Namespace

	// encoders is the namespace for registering directive executors that are
	// used to encode the input struct to HTTP request.
	encoders *owl.Namespace
}

func NewDirectiveNamespace() *DirectiveNamespace {
	ns := &DirectiveNamespace{
		decoders: owl.NewNamespace(),
		encoders: owl.NewNamespace(),
	}
	ns.registerBuiltinDirectives()
	return ns
}

func (ns *DirectiveNamespace) registerBuiltinDirectives() {
	ns.RegisterDirective("form", &DirectvieForm{})
	ns.RegisterDirective("query", &DirectiveQuery{})
	ns.RegisterDirective("header", &DirectiveHeader{})
	ns.RegisterDirective("body", &DirectiveBody{})
	ns.RegisterDirective("required", &DirectiveRequired{})
	ns.RegisterDirective("default", &DirectiveDefault{})
	ns.RegisterDirective("nonzero", &DirectiveNonzero{})
	ns.registerDirective("path", defaultPathDirective)
	ns.registerDirective("omitempty", &DirectiveOmitEmpty{})

	// The following are 3 special executors which do nothing. Each of them just
	// serves as an indicator of overriding the codec for a specific field. For
	// historical reasons, we have three of them, but we should only use "codec"
	// in the future. The "decoder" and "coder" are kept for backward
	// compatibility, and will be removed soon.
	ns.registerDirective("decoder", &DirectiveNoop{})
	ns.registerDirective("coder", &DirectiveNoop{})
	ns.registerDirective("codec", &DirectiveNoop{})
}

// RegisterDirective registers a DirectiveExecutor with the given directive name. The
// directive should be able to both extract the value from the HTTP request and build
// the HTTP request from the value. The Decode API is used to decode data from the HTTP
// request to a field of the input struct, and Encode API is used to encode the field of
// the input struct to the HTTP request.
//
// Will panic if the name were taken or given executor is nil. Pass parameter force
// (true) to ignore the name conflict.
func (ns *DirectiveNamespace) RegisterDirective(name string, executor DirectiveExecutor, force ...bool) {
	panicOnReservedExecutorName(name)
	ns.registerDirective(name, executor, force...)
}

func (ns *DirectiveNamespace) registerDirective(name string, executor DirectiveExecutor, force ...bool) {
	registerDirectiveExecutorToNamespace(ns.decoders, name, executor, force...)
	registerDirectiveExecutorToNamespace(ns.encoders, name, executor, force...)
}

func registerDirectiveExecutorToNamespace(ns *owl.Namespace, name string, exe DirectiveExecutor, force ...bool) {
	if exe == nil {
		internal.PanicOnError(errors.New("nil directive executor"))
	}
	if ns == decoderNamespace {
		ns.RegisterDirectiveExecutor(name, asOwlDirectiveExecutor(exe.Decode), force...)
	} else {
		ns.RegisterDirectiveExecutor(name, asOwlDirectiveExecutor(exe.Encode), force...)
	}
}

func asOwlDirectiveExecutor(directiveFunc func(*DirectiveRuntime) error) owl.DirectiveExecutor {
	return owl.DirectiveExecutorFunc(func(dr *owl.DirectiveRuntime) error {
		return directiveFunc((*DirectiveRuntime)(dr))
	})
}

func panicOnReservedExecutorName(name string) {
	for _, reservedName := range reservedExecutorNames {
		if name == reservedName {
			internal.PanicOnError(fmt.Errorf("reserved executor name: %q", name))
		}
	}
}
