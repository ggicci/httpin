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

var reservedExecutorNames = []string{"decoder"}

func init() {
	// Built-in Directives
	RegisterDirectiveExecutor("form", DirectiveExecutorFunc(formValueExtractor))
	RegisterDirectiveExecutor("query", DirectiveExecutorFunc(queryValueExtractor))
	RegisterDirectiveExecutor("header", DirectiveExecutorFunc(headerValueExtractor))
	RegisterDirectiveExecutor("body", DirectiveExecutorFunc(bodyDecoder))
	RegisterDirectiveExecutor("required", DirectiveExecutorFunc(required))
	RegisterDirectiveExecutor("default", DirectiveExecutorFunc(defaultValueSetter))

	// decoder is a special executor which does nothing, but is an indicator of
	// overriding the decoder for a specific field.
	owl.RegisterDirectiveExecutor("decoder", DirectiveExecutorFunc(nil))
}

// RegisterDirectiveExecutor registers a named executor globally, which
// implemented the DirectiveExecutor interface. Will panic if the name were
// taken or nil executor.
func RegisterDirectiveExecutor(name string, exe DirectiveExecutor) {
	panicOnReservedExecutorName(name)
	owl.RegisterDirectiveExecutor(name, exe)
}

// ReplaceDirectiveExecutor works like RegisterDirectiveExecutor without panic
// on duplicate names.
func ReplaceDirectiveExecutor(name string, exe DirectiveExecutor) {
	panicOnReservedExecutorName(name)
	owl.ReplaceDirectiveExecutor(name, exe)
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
	return decoderOf(elemType)
}

func (rw *directiveRuntimeHelper) DeliverContextValue(key, value interface{}) {
	rw.DirectiveRuntime.Context = context.WithValue(
		rw.DirectiveRuntime.Context, key, value,
	)
}
