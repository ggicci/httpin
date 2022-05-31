package httpin

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

var (
	globalCustomErrorHandler ErrorHandler = defaultErrorHandler
)

type middleware = func(http.Handler) http.Handler
type ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error)

// NewInput creates a "Middleware Constructor" for making a chain, which acts as
// a list of http.Handler constructors. We recommend using
// https://github.com/justinas/alice to chain your HTTP middleware functions and
// the app handler.
func NewInput(inputStruct interface{}, opts ...Option) middleware {
	engine, err := New(inputStruct, opts...)
	if err != nil {
		panic(err)
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			// Here we read the request and decode it to fill our structure.
			// Once failed, the request should end here.
			input, err := engine.Decode(r)
			if err != nil {
				engine.getErrorHandler()(rw, r, err)
				return
			}

			// We put the `input` to the request's context, and it will pass to the next hop.
			ctx := context.WithValue(r.Context(), Input, input)
			next.ServeHTTP(rw, r.WithContext(ctx))
		})
	}
}

func ReplaceDefaultErrorHandler(custom ErrorHandler) {
	if custom == nil {
		panic(fmt.Errorf("httpin: %w", ErrNilErrorHandler))
	}
	globalCustomErrorHandler = custom
}

func defaultErrorHandler(rw http.ResponseWriter, r *http.Request, err error) {
	var invalidFieldError *InvalidFieldError
	if errors.As(err, &invalidFieldError) {
		rw.Header().Add("Content-Type", "application/json")
		rw.WriteHeader(http.StatusUnprocessableEntity) // status: 422
		json.NewEncoder(rw).Encode(invalidFieldError)
		return
	}

	http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError) // status: 500
}
