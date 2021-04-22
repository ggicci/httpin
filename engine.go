package httpin

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
)

func NewEngine(inputStruct interface{}) (*Engine, error) {
	typ := reflect.TypeOf(inputStruct) // retrieve type information
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.Struct {
		return nil, UnsupportedType(typ.Name())
	}

	// TODO(ggicci): check typ
	engine := &Engine{
		inputType: typ,
	}

	// if err := engine.build(); err != nil {
	// 	return nil, fmt.Errorf("httpin: build: %w", err)
	// }

	return engine, nil
}

type Engine struct {
	inputType reflect.Type
}

func (e *Engine) ReadRequest(r *http.Request) (interface{}, error) {
	return nil, nil
}

func (e *Engine) ReadForm(form url.Values) (interface{}, error) {
	rv, err := readForm(e.inputType, form)
	if err != nil {
		return nil, err
	}
	return rv.Interface(), nil
}

func (e *Engine) ReadBody(body io.Reader) (interface{}, error) {
	rv := e.newInstance()

	if err := json.NewDecoder(body).Decode(rv.Interface()); err != nil {
		return nil, fmt.Errorf("httpin: json decode: %w", err)
	}
	return rv.Interface(), nil
}

// newInstance creates a new instance of the input struct.
func (e *Engine) newInstance() reflect.Value {
	return reflect.New(e.inputType)
}

// build builds extractors for the exported fields of the input struct.
func (e *Engine) build() error {
	return errors.New("not implemented")
}
