package core

import (
	"fmt"
	"sync"
)

var defaultNS = NewNamespace()

type Namespace struct {
	*DirectiveNamespace          // for registering directive executors
	builtResolvers      sync.Map // map[reflect.Type]*owl.Resolver
}

func NewNamespace() *Namespace {
	ns := &Namespace{
		DirectiveNamespace: NewDirectiveNamespace(),
	}
	return ns
}

func (ns *Namespace) New(inputStruct any, opts ...Option) (*Core, error) {
	resolver, err := buildResolver(inputStruct)
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

// RegisterDirective registers a DirectiveExecutor with the given directive name. The
// directive should be able to both extract the value from the HTTP request and build
// the HTTP request from the value. The Decode API is used to decode data from the HTTP
// request to a field of the input struct, and Encode API is used to encode the field of
// the input struct to the HTTP request.
//
// Will panic if the name were taken or given executor is nil. Pass parameter force
// (true) to ignore the name conflict.
func RegisterDirective(name string, executor DirectiveExecutor, force ...bool) {
	defaultNS.RegisterDirective(name, executor, force...)
}
