package httpin

import (
	"errors"
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

	RegisterDirectiveExecutor("form", DirectiveExecutorFunc(FormValueExtractor))
	RegisterDirectiveExecutor("header", DirectiveExecutorFunc(HeaderValueExtractor))
	RegisterDirectiveExecutor("body", DirectiveExecutorFunc(BodyDecoder))
}

type DirectiveExecutor interface {
	Execute(*DirectiveContext) error
}

type DirectiveExecutorFunc func(*DirectiveContext) error

func (f DirectiveExecutorFunc) Execute(ctx *DirectiveContext) error {
	return f(ctx)
}

type Directive struct {
	Executor string // name of the executor
	Args     string // args
}

type DirectiveContext struct {
	Request *http.Request
	Value   reflect.Value
}

func RegisterDirectiveExecutor(name string, exe DirectiveExecutor) {
	if _, ok := executors[name]; ok {
		panic(fmt.Sprintf("duplicate executor: %q", name))
	}

	executors[name] = exe
}

func buildDirective(directive string) (*Directive, error) {
	// e.g. form=page, header=x-api-token
	// TODO(ggicci): validate executor
	parts := strings.SplitN(directive, "=", 2)
	if len(parts) == 1 {
		return &Directive{Executor: parts[0]}, nil
	}
	return &Directive{Executor: parts[0], Args: parts[1]}, nil
}

func (d *Directive) Execute(ctx *DirectiveContext) error {
	// Lookup an executor and execute it (MUST exist).
	return executors[d.Executor].Execute(ctx)
}

func FormValueExtractor(ctx *DirectiveContext) error {
	return errors.New("not implemented")
}

func HeaderValueExtractor(ctx *DirectiveContext) error {
	return errors.New("not implemented")
}

func BodyDecoder(ctx *DirectiveContext) error {
	return errors.New("not implemented")
}
