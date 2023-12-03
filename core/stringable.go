package core

import (
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/ggicci/httpin/internal"
)

type Stringable = internal.Stringable

var timeType = internal.TypeOf[time.Time]()

func NewStringable(rv reflect.Value, adapt AnyStringableAdaptor) (stringable Stringable, err error) {
	if IsPatchField(rv.Type()) {
		stringable, err = newPatchFieldStringableWrapper(rv, adapt)
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
func newStringable(rv reflect.Value, adapt AnyStringableAdaptor) (Stringable, error) {
	rv, err := getPointer(rv)
	if err != nil {
		return nil, err
	}

	// Now rv is a pointer type.
	if adapt != nil {
		return adapt(rv.Interface())
	}

	// Custom type adaptors go first. Which means the coder of a specific type
	// has already been registered/overridden by user.
	if adapt, ok := customStringableAdaptors[rv.Type().Elem()]; ok {
		return adapt(rv.Interface())
	}

	// For the base type time.Time, it is a special case here.
	// We won't use TextMarshaler/TextUnmarshaler for time.Time.
	if rv.Type().Elem() == timeType {
		return internal.NewStringable(rv)
	}

	if hybridCoder := hybridizeCoder(rv); hybridCoder != nil {
		return hybridCoder, nil
	}

	// Fallback to use built-in stringable types.
	return internal.NewStringable(rv)
}

type patchFieldStringableWrapper struct {
	Value              reflect.Value // of patch.Field[T]
	internalStringable Stringable
}

func newPatchFieldStringableWrapper(rv reflect.Value, adapt AnyStringableAdaptor) (patchFieldStringableWrapper, error) {
	stringable, err := newStringable(rv.FieldByName("Value"), adapt)
	if err != nil {
		return patchFieldStringableWrapper{}, fmt.Errorf("cannot create Stringable for PatchField: %w", err)
	} else {
		return patchFieldStringableWrapper{
			Value:              rv,
			internalStringable: stringable,
		}, nil
	}
}

func (w patchFieldStringableWrapper) ToString() (string, error) {
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
func (w patchFieldStringableWrapper) FromString(s string) error {
	if err := w.internalStringable.FromString(s); err != nil {
		return err
	} else {
		w.Value.FieldByName("Valid").SetBool(true)
		return nil
	}
}

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
		return rv, fmt.Errorf("cannot get address of value %q", rv)
	}
	rv = rv.Addr()
	return rv, nil
}
