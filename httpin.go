package httpin

import (
	"context"
	"fmt"
	"net/http"
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
