package httpin

import (
	"context"
	"fmt"
	"reflect"

	"github.com/ggicci/owl"
)

type (
	Directive             = owl.Directive
	DirectiveRuntime      = owl.DirectiveRuntime
	DirectiveExecutor     = owl.DirectiveExecutor
	DirectiveExecutorFunc = owl.DirectiveExecutorFunc
)

var (
	reservedExecutorNames = []string{"decoder"}

	// decoderNamespace is the namespace for registering directive executors that are
	// used to decode the http request to input struct.
	decoderNamespace = owl.NewNamespace()

	// encoderNamespace is the namespace for registering directive executors that are
	// used to encode the input struct to http request.
	encoderNamespace = owl.NewNamespace()
)

func init() {
	// Built-in Directives
	RegisterDirectiveExecutor("form", DirectiveExecutorFunc(formValueExtractor), noopDirective)
	RegisterDirectiveExecutor("query", DirectiveExecutorFunc(queryValueExtractor), noopDirective)
	RegisterDirectiveExecutor("header", DirectiveExecutorFunc(headerValueExtractor), noopDirective)
	RegisterDirectiveExecutor("body", DirectiveExecutorFunc(bodyDecoder), noopDirective)
	RegisterDirectiveExecutor("required", DirectiveExecutorFunc(required), noopDirective)
	RegisterDirectiveExecutor("default", DirectiveExecutorFunc(defaultValueSetter), noopDirective)

	// decoder is a special executor which does nothing, but is an indicator of
	// overriding the decoder for a specific field.
	decoderNamespace.RegisterDirectiveExecutor("decoder", DirectiveExecutorFunc(nil))
}

// RegisterDirectiveExecutor registers a directive executor globally, which has the
// DirectiveExecutor interface implemented. The executor will be registered to the
// decoder namespace, only used by the decoding, i.e. decode http request to struct.
// Will panic if the name were taken or nil executor. Pass parameter replace (true) to
// ignore the name conflict.
func RegisterDirectiveExecutor(name string, dec, enc DirectiveExecutor, replace ...bool) {
	registerDirectiveExecutorToNamespace(decoderNamespace, name, dec, replace...)
	registerDirectiveExecutorToNamespace(encoderNamespace, name, enc, replace...)
}

func registerDirectiveExecutorToNamespace(ns *owl.Namespace, name string, exe DirectiveExecutor, replace ...bool) {
	panicOnReservedExecutorName(name)
	ns.RegisterDirectiveExecutor(name, exe, replace...)
}

func panicOnReservedExecutorName(name string) {
	for _, reservedName := range reservedExecutorNames {
		if name == reservedName {
			panic(fmt.Errorf("httpin: %w: %q", ErrReservedExecutorName, name))
		}
	}
}

type directiveRuntimeHelper struct {
	*DirectiveRuntime
}

func (rw *directiveRuntimeHelper) decoderOf(elemType reflect.Type) interface{} {
	decoder := rw.DirectiveRuntime.Resolver.Context.Value(CustomDecoder)
	if decoder != nil {
		return decoder
	}
	return decoderByType(elemType)
}

func (rw *directiveRuntimeHelper) DeliverContextValue(key, value interface{}) {
	rw.DirectiveRuntime.Context = context.WithValue(
		rw.DirectiveRuntime.Context, key, value,
	)
}

var noopDirective = DirectiveExecutorFunc(nil)
