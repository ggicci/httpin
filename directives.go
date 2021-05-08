package httpin

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"strings"
)

type DirectiveExecutor interface {
	Execute(*DirectiveContext) error
}

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

// RegisterDirectiveExecutor registers a named executor globally, which
// implemented the DirectiveExecutor interface. Will panic if the name were
// taken or nil executor.
func RegisterDirectiveExecutor(name string, exe DirectiveExecutor) {
	if _, ok := executors[name]; ok {
		panic(fmt.Sprintf("duplicate executor: %q", name))
	}
	ReplaceDirectiveExecutor(name, exe)
}

// ReplaceDirectiveExecutor works like RegisterDirectiveExecutor without panic
// on duplicate names.
func ReplaceDirectiveExecutor(name string, exe DirectiveExecutor) {
	if exe == nil {
		panic(fmt.Sprintf("nil executor: %q", name))
	}
	executors[name] = exe
	debug("directive executor replaced: %q\n", name)
}

type DirectiveExecutorFunc func(*DirectiveContext) error

func (f DirectiveExecutorFunc) Execute(ctx *DirectiveContext) error {
	return f(ctx)
}

type DirectiveContext struct {
	directive
	ValueType reflect.Type
	Value     reflect.Value
	Request   *http.Request
	Context   context.Context
}

func (c *DirectiveContext) DeliverContextValue(key, val interface{}) {
	c.Context = context.WithValue(c.Context, key, val)
}

type directive struct {
	Executor string   // name of the executor
	Argv     []string // argv
}

// buildDirective builds a Directive instance by parsing a directive string.
// Example directives are:
//   - form=page,page_index -> { Executor: "form", Args: ["page", "page_index"] }
//   - header=x-api-token   -> { Executor: "header", Args: ["x-api-token"] }
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

	// TODO(ggicci): hook custom validators, e.g. dir.Validate()
	return dir, nil
}

func (d *directive) Execute(ctx *DirectiveContext) error {
	return d.getExecutor().Execute(ctx)
}

func (d *directive) getExecutor() DirectiveExecutor {
	return executors[d.Executor]
}
