package httpin

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
)

func New(inputStruct interface{}) Middleware {
	engine, err := NewEngine(inputStruct)
	if err != nil {
		panic(fmt.Errorf("httpin: unable to create engine: %w", err))
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			// Here we read the request and decode it to fill our structure.
			// Once failed, the request should end here.
			input, err := engine.ReadRequest(r)
			if err != nil {
				http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}

			// We put the `input` to the request's context, and it will pass to the next hop.
			ctx := context.WithValue(r.Context(), "httpin", input)
			next.ServeHTTP(rw, r.WithContext(ctx))
		})
	}
}

func NewEngine(inputStruct interface{}) (*Engine, error) {
	typ := reflect.TypeOf(inputStruct) // retrieve type information
	// TODO(ggicci): check typ
	engine := &Engine{
		inputType: typ,
	}

	if err := engine.build(); err != nil {
		return nil, fmt.Errorf("httpin: build: %w", err)
	}

	return engine, nil
}

type Engine struct {
	inputType reflect.Type
}

func (e *Engine) ReadRequest(r *http.Request) (interface{}, error) {
	return nil, nil
}

func (e *Engine) ReadForm(form url.Values) (interface{}, error) {
	return nil, errors.New("not implemented")
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
