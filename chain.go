package httpin

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
)

// NewInput creates a "Middleware Constructor" for making a chain, which acts as
// a list of http.Handler constructors. We recommend using
// https://github.com/justinas/alice to chain your HTTP middleware functions and
// the app handler.
func NewInput(inputStruct interface{}) func(http.Handler) http.Handler {
	core, err := New(inputStruct)
	if err != nil {
		panic(err)
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			// Here we read the request and decode it to fill our structure.
			// Once failed, the request should end here.
			input, err := core.Decode(r)
			if err != nil {
				var invalidFieldError *InvalidField
				if errors.As(err, &invalidFieldError) {
					// TODO(ggicci): options to tweak the response
					rw.Header().Add("Content-Type", "application/json")
					rw.WriteHeader(422)
					json.NewEncoder(rw).Encode(invalidFieldError)
					return
				}

				http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			// We put the `input` to the request's context, and it will pass to the next hop.
			ctx := context.WithValue(r.Context(), Input, input)
			next.ServeHTTP(rw, r.WithContext(ctx))
		})
	}
}
