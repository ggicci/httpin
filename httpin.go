package httpin

import (
	"fmt"
	"net/http"
	"reflect"
)

type ContextKey int

const (
	Input ContextKey = iota // the primary key to get the input object in the context injected by httpin

	FieldSet
)

type core struct {
	inputType reflect.Type
	tree      *FieldResolver
}

func New(inputStruct interface{}) (*core, error) {
	typ := reflect.TypeOf(inputStruct) // retrieve type information
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.Struct {
		return nil, UnsupportedTypeError{Type: typ}
	}

	engine := &core{
		inputType: typ,
	}

	if err := engine.build(); err != nil {
		return nil, fmt.Errorf("httpin: %w", err)
	}

	return engine, nil
}

func (e *core) Decode(req *http.Request) (interface{}, error) {
	if err := req.ParseForm(); err != nil {
		return nil, err
	}
	rv, err := e.tree.resolve(req)
	if err != nil {
		return nil, fmt.Errorf("httpin: %w", err)
	}
	return rv.Interface(), nil
}

// build builds extractors for the exported fields of the input struct.
func (e *core) build() error {
	tree, err := buildResolverTree(e.inputType)
	if err != nil {
		return err
	}
	e.tree = tree
	return nil
}
