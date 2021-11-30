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

	// Set this context value to true to indicate that the field has been set.
	// When multiple executors were applied to a field, if the field value were set by
	// an executor, the latter executors may skip running by consulting this context value.
	FieldSet

	StopRecursion
)

var builtEngines sync.Map

// Engine holds the information on how to decode a request to an instance of a
// concrete struct type.
type Engine struct {
	// core
	inputType reflect.Type
	tree      *fieldResolver

	// options
	errorHandler ErrorHandler
}

// New builds an HTTP request decoder for the specified struct type with custom options.
func New(inputStruct interface{}, opts ...Option) (*Engine, error) {
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

	var core *Engine

	builtEngine, built := builtEngines.Load(typ)
	if !built {
		// Build the engine core if not built yet.
		core = &Engine{inputType: typ}
		if err := core.build(); err != nil {
			return nil, fmt.Errorf("httpin: %w", err)
		}
		builtEngines.Store(typ, core)
	} else {
		// Load the engine core and get a copy.
		core = copyEngineCore(builtEngine.(*Engine))
	}

	// Apply default options and user custom options to the engine.
	var allOptions []Option
	// defaultOptions := []Option{}
	// allOptions = append(allOptions, defaultOptions...)
	allOptions = append(allOptions, opts...)

	for _, opt := range allOptions {
		if err := opt(core); err != nil {
			return nil, fmt.Errorf("httpin: invalid option: %w", err)
		}
	}

	return core, nil
}

// Decode decodes an HTTP request to a struct instance.
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

func copyEngineCore(eng *Engine) *Engine {
	return &Engine{inputType: eng.inputType, tree: eng.tree}
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

func (e *Engine) getErrorHandler() ErrorHandler {
	if e.errorHandler != nil {
		return e.errorHandler
	}

	return globalCustomErrorHandler
}
