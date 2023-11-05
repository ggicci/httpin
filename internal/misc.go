package internal

import (
	"fmt"
	"reflect"
)

func IsNil(value reflect.Value) bool {
	switch value.Kind() {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.Interface, reflect.Slice:
		return value.IsNil()
	default:
		return false
	}
}

func PanicOnError(err error) {
	if err != nil {
		panic(fmt.Errorf("httpin: %w", err))
	}
}

// TypeOf returns the reflect.Type of a given type.
// e.g. TypeOf[int]() returns reflect.TypeOf(0)
func TypeOf[T any]() reflect.Type {
	var zero [0]T
	return reflect.TypeOf(zero).Elem()
}

func Pointerize[T any](v T) *T {
	return &v
}
