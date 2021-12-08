package httpin

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"strings"
)

var (
	executors   map[string]DirectiveExecutor
	normalizers map[string]DirectiveNormalizer
)

func init() {
	executors = make(map[string]DirectiveExecutor)
	normalizers = make(map[string]DirectiveNormalizer)

	// Built-in Directives
	RegisterDirectiveExecutor("form", DirectiveExecutorFunc(formValueExtractor), nil)
	RegisterDirectiveExecutor("query", DirectiveExecutorFunc(queryValueExtractor), nil)
	RegisterDirectiveExecutor("header", DirectiveExecutorFunc(headerValueExtractor), nil)
	RegisterDirectiveExecutor(
		"body",
		DirectiveExecutorFunc(bodyDecoder),
		DirectiveNormalizerFunc(bodyDirectiveNormalizer),
	)
	RegisterDirectiveExecutor("required", DirectiveExecutorFunc(required), nil)
	// RegisterDirectiveExecutor("file", DirectiveExecutorFunc(fileValueExtractor), nil)
}

// DirectiveExecutor is the interface implemented by a "directive executor".
type DirectiveExecutor interface {
	Execute(*DirectiveContext) error
}

type DirectiveNormalizer interface {
	Normalize(*Directive) error
}

// RegisterDirectiveExecutor registers a named executor globally, which
// implemented the DirectiveExecutor interface. Will panic if the name were
// taken or nil executor.
func RegisterDirectiveExecutor(name string, exe DirectiveExecutor, norm DirectiveNormalizer) {
	if _, ok := executors[name]; ok {
		panic(fmt.Errorf("%w: %q", ErrDuplicateExecutor, name))
	}
	ReplaceDirectiveExecutor(name, exe, norm)
}

// ReplaceDirectiveExecutor works like RegisterDirectiveExecutor without panic
// on duplicate names.
func ReplaceDirectiveExecutor(name string, exe DirectiveExecutor, norm DirectiveNormalizer) {
	if exe == nil {
		panic(fmt.Errorf("%w: %q", ErrNilExecutor, name))
	}
	executors[name] = exe
	normalizers[name] = norm
}

// DirectiveExecutorFunc is an adpator to allow to use of ordinary functions as
// httpin.DirectiveExecutor.
type DirectiveExecutorFunc func(*DirectiveContext) error

// Execute calls f(ctx).
func (f DirectiveExecutorFunc) Execute(ctx *DirectiveContext) error {
	return f(ctx)
}

// DirectiveNormalizerFunc is an adaptor to allow to use of ordinary functions as
// httpin.DirectiveNormalizer.
type DirectiveNormalizerFunc func(*Directive) error

// Normalize calls f(dir).
func (f DirectiveNormalizerFunc) Normalize(dir *Directive) error {
	return f(dir)
}

// DirectiveContext holds essential information about the field being resolved
// and the active HTTP request. Working as the context in a directive executor.
type DirectiveContext struct {
	Directive
	ValueType reflect.Type
	Value     reflect.Value
	Request   *http.Request
	Context   context.Context
}

// DeliverContextValue binds a value to the specified key in the context. And it
// will be delivered among the executors in the same field resolver.
func (c *DirectiveContext) DeliverContextValue(key, value interface{}) {
	c.Context = context.WithValue(c.Context, key, value)
}

// Directive defines the profile to locate an httpin.DirectiveExecutor instance
// and drive it with essential arguments.
type Directive struct {
	Executor string   // name of the executor
	Argv     []string // argv
}

// buildDirective builds a `directive` by parsing a directive string extracted
// from the struct tag.
//
// Example directives are:
//    "form=page,page_index" -> { Executor: "form", Args: ["page", "page_index"] }
//    "header=x-api-token"   -> { Executor: "header", Args: ["x-api-token"] }
func buildDirective(directiveStr string) (*Directive, error) {
	parts := strings.SplitN(directiveStr, "=", 2)
	executor := parts[0]
	var argv []string
	if len(parts) == 2 {
		// Split the remained string by delimiter `,` as argv.
		argv = strings.Split(parts[1], ",")
	}

	// Ensure that the corresponding executor had been registered.
	dir := &Directive{Executor: executor, Argv: argv}
	if dir.getExecutor() == nil {
		return nil, fmt.Errorf("%w: %q", ErrUnregisteredExecutor, dir.Executor)
	}

	// Normalize the directive.
	norm := dir.getNormalizer()
	if norm != nil {
		if err := norm.Normalize(dir); err != nil {
			return nil, fmt.Errorf("invalid directive %q: %w", dir.Executor, err)
		}
	}

	return dir, nil
}

// Execute locates the executor and runs it with the specified context.
func (d *Directive) Execute(ctx *DirectiveContext) error {
	return d.getExecutor().Execute(ctx)
}

// getExecutor locates the executor by its name. It must exist.
func (d *Directive) getExecutor() DirectiveExecutor {
	return executors[d.Executor]
}

// getNormalizer locates the directive normalizer by its name.
func (d *Directive) getNormalizer() DirectiveNormalizer {
	return normalizers[d.Executor]
}
