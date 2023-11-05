package internal

import "reflect"

func IsNil(value reflect.Value) bool {
	switch value.Kind() {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.Interface, reflect.Slice:
		return value.IsNil()
	default:
		return false
	}
}

func Pointerize[T any](v T) *T {
	return &v
}
