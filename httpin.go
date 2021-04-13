package httpin

import (
	"context"
	"encoding/json"
	"net/http"
	"reflect"
)

func Input(inputStruct interface{}) Middleware {
	typ := reflect.TypeOf(inputStruct) // retrieve type information
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			rv := reflect.New(typ) // create new instance of the type

			// Here we read the request body and decode it to fill our structure.
			// Once failed, the request should end here.
			if err := json.NewDecoder(r.Body).Decode(rv.Interface()); err != nil {
				http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest) // 400
				return
			}

			// We put the `input` to the request's context, and it will pass to the next hop.
			ctx := context.WithValue(r.Context(), "httpin", rv.Interface())
			next.ServeHTTP(rw, r.WithContext(ctx))
		})
	}
}
