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

func RegisterDirectiveExecutor(name string, exe DirectiveExecutor) {
	if _, ok := executors[name]; ok {
		panic(fmt.Sprintf("duplicate executor: %q", name))
	}
	executors[name] = exe
	fmt.Printf("directive executor registered: %q\n", name)
}

func buildDirective(directive string) (*Directive, error) {
	// e.g. form=page, header=x-api-token
	parts := strings.SplitN(directive, "=", 2)
	executor := parts[0]
	var argv []string
	if len(parts) == 2 {
		// Split remained string as argv.
		// e.g. form=page,index, argv = ["page", "index"]
		argv = strings.Split(parts[1], ",")
	}

	// Validate the directive.
	dir := &Directive{Executor: executor, Argv: argv}
	if dir.getExecutor() == nil {
		return nil, fmt.Errorf("invalid directive %q since executor %q not registered",
			directive, dir.Executor)
	}
	return dir, nil
}

func (d *Directive) Execute(ctx *DirectiveContext) error {
	return d.getExecutor().Execute(ctx)
}

func extractFromKVSWithKey(ctx *DirectiveContext, kvs map[string][]string, key string) error {
	if ctx.Context.Value(fieldSet) == true {
		fmt.Printf("    > field already set, skip\n")
		return nil
	}

	// FIXME(ggicci): get converter and check type first, fail fast here
	got, exists := kvs[key]
	if !exists || len(got) == 0 {
		fmt.Printf("    > key %q not found in %s\n", key, ctx.Executor)
		return nil
	}

	if isBasicType(ctx.ValueType) {
		if err := setBasicValue(ctx.Value.Elem(), ctx.ValueType, got[0]); err != nil {
			return err
		}
		ctx.DeliverContextValue(fieldSet, true)
		return nil
	}

	if isTimeType(ctx.ValueType) {
		if err := setTimeValue(ctx.Value.Elem(), ctx.ValueType, got[0]); err != nil {
			return err
		}
		ctx.DeliverContextValue(fieldSet, true)
		return nil
	}

	if isArrayType(ctx.ValueType) {
		if err := setSliceValue(ctx.Value.Elem(), ctx.ValueType, got); err != nil {
			return err
		}
		ctx.DeliverContextValue(fieldSet, true)
		return nil
	}

	return UnsupportedTypeError{ctx.ValueType, ctx.Directive.Executor}
}

func extractFromKVS(ctx *DirectiveContext, kvs map[string][]string, headerKey bool) error {
	for _, key := range ctx.Directive.Argv {
		fmt.Printf("    > execute directive %q with key %q\n", ctx.Directive.Executor, key)
		if headerKey {
			key = http.CanonicalHeaderKey(key)
		}
		if err := extractFromKVSWithKey(ctx, kvs, key); err != nil {
			return err
		}
	}
	return nil
}

func FormValueExtractor(ctx *DirectiveContext) error {
	return extractFromKVS(ctx, ctx.Request.Form, false)
}

func HeaderValueExtractor(ctx *DirectiveContext) error {
	return extractFromKVS(ctx, ctx.Request.Header, true)
}

func BodyDecoder(ctx *DirectiveContext) error {
	// TODO(ggicci): implement this
	return nil
}

func RequireField(ctx *DirectiveContext) error {
	if ctx.Context.Value(fieldSet) == nil {
		return MissingFieldError
	}
	return nil
}
