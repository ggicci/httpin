package codec

import (
	"fmt"
	"reflect"
)

func getPointer(rv reflect.Value) (reflect.Value, error) {
	if rv.Kind() == reflect.Pointer {
		return createInstanceIfNil(rv), nil
	} else {
		return addressOf(rv)
	}
}

func createInstanceIfNil(rv reflect.Value) reflect.Value {
	if rv.IsNil() {
		rv.Set(reflect.New(rv.Type().Elem()))
	}
	return rv
}

func addressOf(rv reflect.Value) (reflect.Value, error) {
	if !rv.CanAddr() {
		return rv, fmt.Errorf("cannot get address of value %v", rv)
	}
	rv = rv.Addr()
	return rv, nil
}
