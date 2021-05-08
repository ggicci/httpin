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

type directiveContext int

const (
	fieldSet directiveContext = iota
)

func init() {
	executors = make(map[string]DirectiveExecutor)

	RegisterDirectiveExecutor("form", DirectiveExecutorFunc(FormValueExtractor))
	RegisterDirectiveExecutor("header", DirectiveExecutorFunc(HeaderValueExtractor))
	RegisterDirectiveExecutor("body", DirectiveExecutorFunc(BodyDecoder))
	RegisterDirectiveExecutor("required", DirectiveExecutorFunc(RequireField))
}

// RegisterDirectiveExecutor registers a named executor globally, which
// implemented the DirectiveExecutor interface.
func RegisterDirectiveExecutor(name string, exe DirectiveExecutor) {
	if _, ok := executors[name]; ok {
		panic(fmt.Sprintf("duplicate executor: %q", name))
	}
	if exe == nil {
		panic(fmt.Sprintf("nil executor: %q", name))
	}
	executors[name] = exe
	debug("directive executor registered: %q\n", name)
}

type DirectiveExecutor interface {
	Execute(*DirectiveContext) error
}

type DirectiveExecutorFunc func(*DirectiveContext) error

func (f DirectiveExecutorFunc) Execute(ctx *DirectiveContext) error {
	return f(ctx)
}

type Directive struct {
	Executor string   // name of the executor
	Argv     []string // argv
}

// BuildDirective builds a Directive instance by parsing a directive string.
// Example directives are:
//   - form=page,page_index -> { Executor: "form", Args: ["page", "page_index"] }
//   - header=x-api-token   -> { Executor: "header", Args: ["x-api-token"] }
func BuildDirective(directive string) (*Directive, error) {
	parts := strings.SplitN(directive, "=", 2)
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

	// TODO(ggicci): hook custom validators, e.g. dir.Validate()
	return dir, nil
}

func (d *Directive) Execute(ctx *DirectiveContext) error {
	return d.getExecutor().Execute(ctx)
}

func (d *Directive) getExecutor() DirectiveExecutor {
	return executors[d.Executor]
}

type DirectiveContext struct {
	Directive
	ValueType reflect.Type
	Value     reflect.Value
	Request   *http.Request
	Context   context.Context
}

func (c *DirectiveContext) DeliverContextValue(key, val interface{}) {
	c.Context = context.WithValue(c.Context, key, val)
}
