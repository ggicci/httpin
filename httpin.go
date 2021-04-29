package httpin

import (
	"fmt"
	"net/http"
	"reflect"
)

type ContextKey int

const (
	Input ContextKey = iota // the primary key to get the input object in the context injected by httpin
)

type Core struct {
	inputType reflect.Type
	tree      *FieldResolver
}

func New(inputStruct interface{}, opts ...CoreOption) (*Core, error) {
	typ := reflect.TypeOf(inputStruct) // retrieve type information
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.Struct {
		return nil, UnsupportedTypeError{Type: typ}
	}

	core := &Core{
		inputType: typ,
	}

	if err := core.build(); err != nil {
		return nil, fmt.Errorf("httpin: build: %w", err)
	}

	return core, nil
}

func (e *Core) ReadRequest(r *http.Request) (interface{}, error) {
	return nil, nil
}

// build builds extractors for the exported fields of the input struct.
func (e *Core) build() error {
	tree, err := buildResolverTree(e.inputType)
	if err != nil {
		return err
	}
	e.tree = tree
	return nil
}
