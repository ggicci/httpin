package httpin

import (
	"errors"
	"net/http"
	"reflect"
)

const (
	QueryTag  string = "query"
	HeaderTag string = "header"
	BodyTag   string = "body"
)

type Engine struct {
	inputType reflect.Type
}

func NewEngine(inputStruct interface{}, opts ...EngineOption) (*Engine, error) {
	typ := reflect.TypeOf(inputStruct) // retrieve type information
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.Struct {
		return nil, UnsupportedTypeError{Type: typ}
	}

	engine := &Engine{
		inputType: typ,
	}

	// if err := engine.build(); err != nil {
	// 	return nil, fmt.Errorf("httpin: build: %w", err)
	// }

	return engine, nil
}

func (e *Engine) ReadRequest(r *http.Request) (interface{}, error) {
	return nil, nil
}

// newInstance creates a new instance of the input struct.
func (e *Engine) newInstance() reflect.Value {
	return reflect.New(e.inputType)
}

// build builds extractors for the exported fields of the input struct.
func (e *Engine) build() error {
	return errors.New("not implemented")
}
