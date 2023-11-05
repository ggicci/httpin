package httpin

import (
	"errors"
	"fmt"

	"github.com/ggicci/httpin/directive"
	"github.com/ggicci/httpin/internal"
	"github.com/ggicci/owl"
)

var (
	// decoderNamespace is the namespace for registering directive executors that are
	// used to decode the http request to input struct.
	decoderNamespace = owl.NewNamespace()

	// encoderNamespace is the namespace for registering directive executors that are
	// used to encode the input struct to http request.
	encoderNamespace = owl.NewNamespace()

	reservedExecutorNames = []string{"decoder", "encoder"}
	noopDirective         = &directiveNoop{}
)

type DirectiveRuntime = internal.DirectiveRuntime

type DirectiveExecutor interface {
	Encode(*DirectiveRuntime) error
	Decode(*DirectiveRuntime) error
}

func init() {
	// Register bulit-in directives.
	Customizer().RegisterDirective("form", &directive.DirectvieForm{})
	Customizer().RegisterDirective("query", &directive.DirectiveQuery{})
	Customizer().RegisterDirective("header", &directive.DirectiveHeader{})
	Customizer().RegisterDirective("body", &directive.DirectiveBody{})
	Customizer().RegisterDirective("required", &directive.DirectiveRequired{})
	Customizer().RegisterDirective("default", &directive.DirectiveDefault{})

	// decoder is a special executor which does nothing, but is an indicator of
	// overriding the decoder for a specific field.
	decoderNamespace.RegisterDirectiveExecutor("decoder", asOwlDirectiveExecutor(noopDirective.Decode))
	encoderNamespace.RegisterDirectiveExecutor("encoder", asOwlDirectiveExecutor(noopDirective.Encode))
}

func registerDirectiveExecutorToNamespace(ns *owl.Namespace, name string, exe DirectiveExecutor, force ...bool) {
	panicOnReservedExecutorName(name)
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

// directiveNoop is a DirectiveExecutor that does nothing, "noop" stands for "no operation".
type directiveNoop struct{}

func (*directiveNoop) Encode(*DirectiveRuntime) error { return nil }
func (*directiveNoop) Decode(*DirectiveRuntime) error { return nil }
