package httpin

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"strings"
)

var (
	executors map[string]DirectiveExecutor
)

func init() {
	executors = make(map[string]DirectiveExecutor)

	RegisterDirectiveExecutor("form", DirectiveExecutorFunc(formValueExtractor))
	RegisterDirectiveExecutor("header", DirectiveExecutorFunc(headerValueExtractor))
	RegisterDirectiveExecutor("body", DirectiveExecutorFunc(bodyDecoder))
	RegisterDirectiveExecutor("required", DirectiveExecutorFunc(required))
}

// DirectiveExecutor is the interface implemented by a "directive executor".
type DirectiveExecutor interface {
	Execute(*DirectiveContext) error
}

// RegisterDirectiveExecutor registers a named executor globally, which
// implemented the DirectiveExecutor interface. Will panic if the name were
// taken or nil executor.
func RegisterDirectiveExecutor(name string, exe DirectiveExecutor) {
	if _, ok := executors[name]; ok {
		panic(fmt.Errorf("%w: %q", ErrDuplicateExecutor, name))
	}
	ReplaceDirectiveExecutor(name, exe)
}

// ReplaceDirectiveExecutor works like RegisterDirectiveExecutor without panic
// on duplicate names.
func ReplaceDirectiveExecutor(name string, exe DirectiveExecutor) {
	if exe == nil {
		panic(fmt.Errorf("%w: %q", ErrNilExecutor, name))
	}
	executors[name] = exe
}

// DirectiveExecutorFunc is an adpator to allow to use of ordinary functions as
// httpin.DirectiveExecutor.
type DirectiveExecutorFunc func(*DirectiveContext) error

// Execute calls f(ctx).
func (f DirectiveExecutorFunc) Execute(ctx *DirectiveContext) error {
	return f(ctx)
}

// DirectiveContext holds essential information about the field being resolved
// and the active HTTP request. Working as the context in a directive executor.
type DirectiveContext struct {
	directive
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

// directive defines the profile to locate an httpin.DirectiveExecutor instance
// and drive it with essential arguments.
type directive struct {
	Executor string   // name of the executor
	Argv     []string // argv
}

// buildDirective builds a `directive` by parsing a directive string extracted
// from the struct tag.
//
// Example directives are:
//    "form=page,page_index" -> { Executor: "form", Args: ["page", "page_index"] }
//    "header=x-api-token"   -> { Executor: "header", Args: ["x-api-token"] }
func buildDirective(directiveStr string) (*directive, error) {
	parts := strings.SplitN(directiveStr, "=", 2)
	executor := parts[0]
	var argv []string
	if len(parts) == 2 {
		// Split the remained string by delimiter `,` as argv.
		argv = strings.Split(parts[1], ",")
	}

	// Ensure that the corresponding executor had been registered.
	dir := &directive{Executor: executor, Argv: argv}
	if dir.getExecutor() == nil {
		return nil, fmt.Errorf("%w: %q", ErrUnregisteredExecutor, dir.Executor)
	}

	return dir, nil
}

// Execute locates the executor and runs it with the specified context.
func (d *directive) Execute(ctx *DirectiveContext) error {
	return d.getExecutor().Execute(ctx)
}

// getExecutor locates the executor by its name. It must exist.
func (d *directive) getExecutor() DirectiveExecutor {
	return executors[d.Executor]
}
