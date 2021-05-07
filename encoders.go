package httpin

import "reflect"

// TODO(ggicci):
// InputStruct -> http.Request

type Encoder interface {
	Encode(reflect.Value) ([]string, error)
}
