package core

import (
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/ggicci/httpin/internal"
	"github.com/ggicci/strconvx"
)

type Stringable = strconvx.StringConverter

func NewStringable(rv reflect.Value, adapt AnyStringConverterAdaptor) (stringable Stringable, err error) {
	if IsPatchField(rv.Type()) {
		stringable, err = NewStringablePatchFieldWrapper(rv, adapt)
	} else {
		stringable, err = newStringable(rv, adapt)
	}
	if err != nil {
		return nil, err
	}
	return stringable, nil
}

// Create a Stringable from a reflect.Value. If rv is a pointer type, it will
// try to create a Stringable from rv. Otherwise, it will try to create a
// Stringable from rv.Addr(). Only basic built-in types are supported. As a
// special case, time.Time is also supported.
func newStringable(rv reflect.Value, adapt AnyStringConverterAdaptor) (Stringable, error) {
	rv, err := getPointer(rv)
	if err != nil {
		return nil, err
	}

	// Now rv is a pointer type.
	if adapt != nil {
		return adapt(rv.Interface())
	}

	// For the base type time.Time, it is a special case here.
	// We won't use TextMarshaler/TextUnmarshaler for time.Time.
	// if rv.Type().Elem() == timeType {
	// 	return strconvxNS.New(rv)
	// }

	// Fallback to use built-in stringable types.
	return strconvxNS.New(rv)
}

type StringablePatchFieldWrapper struct {
	Value              reflect.Value // of patch.Field[T]
	internalStringable Stringable
}

func NewStringablePatchFieldWrapper(rv reflect.Value, adapt AnyStringConverterAdaptor) (*StringablePatchFieldWrapper, error) {
	stringable, err := NewStringable(rv.FieldByName("Value"), adapt)
	if err != nil {
		return &StringablePatchFieldWrapper{}, fmt.Errorf("cannot create Stringable for PatchField: %w", err)
	} else {
		return &StringablePatchFieldWrapper{
			Value:              rv,
			internalStringable: stringable,
		}, nil
	}
}

func (w *StringablePatchFieldWrapper) ToString() (string, error) {
	if w.Value.FieldByName("Valid").Bool() {
		return w.internalStringable.ToString()
	} else {
		return "", errors.New("invalid value") // when Valid is false
	}
}

// FromString sets the value of the wrapped patch.Field[T] from the given
// string. It returns an error if the given string is not valid. And leaves the
// original value of both Value and Valid unchanged. On the other hand, if no
// error occurs, it sets Valid to true.
func (w *StringablePatchFieldWrapper) FromString(s string) error {
	if err := w.internalStringable.FromString(s); err != nil {
		return err
	} else {
		w.Value.FieldByName("Valid").SetBool(true)
		return nil
	}
}

var timeType = internal.TypeOf[time.Time]()

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
