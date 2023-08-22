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
	RegisterDirectiveExecutor(
		"body",
		DirectiveExecutorFunc(bodyDecoder),
		// DirectiveNormalizerFunc(bodyDirectiveNormalizer),
	)
	RegisterDirectiveExecutor("required", DirectiveExecutorFunc(required))
	RegisterDirectiveExecutor("default", DirectiveExecutorFunc(defaultValueSetter))

	// decoder is a special executor which does nothing, but is an indicator of
	// overriding the decoder for a specific field.
	owl.RegisterDirectiveExecutor("decoder", DirectiveExecutorFunc(noop))
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

// noop is a no-operation directive executor.
func noop(_ *DirectiveRuntime) error {
	return nil
}

type directiveRuntimeHelper struct {
	*DirectiveRuntime
}

func (h *directiveRuntimeHelper) decoderOf(t reflect.Type) interface{} {
	decoder := h.DirectiveRuntime.Context.Value(CustomDecoder)
	if decoder != nil {
		return decoder
	}
	return decoderOf(t)
}

func (h *directiveRuntimeHelper) DeliverContextValue(key interface{}, value interface{}) {
	h.DirectiveRuntime.Context = context.WithValue(h.DirectiveRuntime.Context, key, value)
}
