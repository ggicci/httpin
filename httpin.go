package httpin

import (
	"fmt"
	"net/http"
	"reflect"
	"sync"
)

type ContextKey int

const (
	Input ContextKey = iota // the primary key to get the input object in the context injected by httpin

	FieldSet
)

var builtEngines sync.Map

type Engine struct {
	// core
	inputType reflect.Type
	tree      *FieldResolver

	// options
	errorStatusCode int
}

func copyEngine(engine *Engine) *Engine {
	return &Engine{inputType: engine.inputType, tree: engine.tree}
}

func New(inputStruct interface{}, opts ...option) (*Engine, error) {
	typ := reflect.TypeOf(inputStruct) // retrieve type information
	if typ == nil {
		return nil, fmt.Errorf("httpin: nil input type")
	}

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.Struct {
		return nil, UnsupportedTypeError{Type: typ}
	}

	var engine *Engine

	builtEngine, built := builtEngines.Load(typ)
	if !built {
		// Build the engine core if not built yet.
		engine = &Engine{
			inputType:       typ,
			errorStatusCode: 422,
		}
		if err := engine.build(); err != nil {
			return nil, fmt.Errorf("httpin: %w", err)
		}
		builtEngines.Store(typ, engine)
	} else {
		// Load the engine core and get a copy.
		engine = copyEngine(builtEngine.(*Engine))
	}

	// Apply default options and user custom options to the engine.
	var allOptions []option
	defaultOptions := []option{
		WithErrorStatusCode(422),
	}
	allOptions = append(allOptions, defaultOptions...)
	allOptions = append(allOptions, opts...)

	for _, opt := range allOptions {
		opt(engine)
	}

	return engine, nil
}

func (e *Engine) Decode(req *http.Request) (interface{}, error) {
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
func (e *Engine) build() error {
	tree, err := buildResolverTree(e.inputType)
	if err != nil {
		return err
	}
	e.tree = tree
	return nil
}
