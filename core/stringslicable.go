package core

import (
	"fmt"
	"reflect"

	"github.com/ggicci/httpin/internal"
)

type StringSlicable interface {
	ToStringSlice() ([]string, error)
	FromStringSlice([]string) error
}

func NewStringSlicable(rv reflect.Value, adapt AnyStringableAdaptor) (StringSlicable, error) {
	if rv.Type().Implements(stringSliceableType) && rv.CanInterface() {
		return rv.Interface().(StringSlicable), nil
	}

	if IsPatchField(rv.Type()) {
		return newPatchFieldStringSlicableWrapper(rv, adapt)
	}

	if isSliceType(rv.Type()) && !isByteSliceType(rv.Type()) {
		return NewStringableArrayBuilder(rv, adapt)
	} else {
		return newStringSlicableFromStringable(rv, adapt)
	}
}

func newStringSlicableFromStringable(rv reflect.Value, adapt AnyStringableAdaptor) (StringSlicable, error) {
	if stringable, err := NewStringable(rv, adapt); err != nil {
		return nil, err
	} else {
		return SingleStringableSlicableWrapper{stringable}, nil
	}
}

type patchFieldStringSlicableWrapper struct {
	Value                   reflect.Value // of patch.Field[T]
	internalStringSliceable StringSlicable
}

func newPatchFieldStringSlicableWrapper(rv reflect.Value, adapt AnyStringableAdaptor) (patchFieldStringSlicableWrapper, error) {
	stringSlicable, err := NewStringSlicable(rv.FieldByName("Value"), adapt)
	if err != nil {
		return patchFieldStringSlicableWrapper{}, fmt.Errorf("cannot create StringSlicable for PatchField: %w", err)
	} else {
		return patchFieldStringSlicableWrapper{
			Value:                   rv,
			internalStringSliceable: stringSlicable,
		}, nil
	}
}

func (w patchFieldStringSlicableWrapper) ToStringSlice() ([]string, error) {
	if w.Value.FieldByName("Valid").Bool() {
		return w.internalStringSliceable.ToStringSlice()
	} else {
		return []string{}, nil
	}
}

func (w patchFieldStringSlicableWrapper) FromStringSlice(values []string) error {
	if err := w.internalStringSliceable.FromStringSlice(values); err != nil {
		return err
	} else {
		w.Value.FieldByName("Valid").SetBool(true)
		return nil
	}
}

type StringableArrayBuilder struct {
	Value reflect.Value
	Adapt AnyStringableAdaptor
}

func NewStringableArrayBuilder(rv reflect.Value, adapt AnyStringableAdaptor) (StringableArrayBuilder, error) {
	if !rv.CanAddr() {
		return StringableArrayBuilder{}, fmt.Errorf("cannot get address of value %q", rv)
	}
	return StringableArrayBuilder{Value: rv, Adapt: adapt}, nil
}

func (b StringableArrayBuilder) ToStringSlice() ([]string, error) {
	var stringables = make(StringableArray, b.Value.Len())
	for i := 0; i < b.Value.Len(); i++ {
		if stringable, err := NewStringable(b.Value.Index(i), b.Adapt); err != nil {
			return nil, fmt.Errorf("cannot create Stringable from %q at index %d: %w", b.Value.Index(i), i, err)
		} else {
			stringables[i] = stringable
		}
	}
	return stringables.ToStringSlice()
}

func (b StringableArrayBuilder) FromStringSlice(ss []string) error {
	var stringables = make(StringableArray, len(ss))
	b.Value.Set(reflect.MakeSlice(b.Value.Type(), len(ss), len(ss)))
	for i := range ss {
		if stringable, err := NewStringable(b.Value.Index(i), b.Adapt); err != nil {
			return fmt.Errorf("cannot create Stringable from %q at index %d: %w", b.Value.Index(i), i, err)
		} else {
			stringables[i] = stringable
		}
	}
	return stringables.FromStringSlice(ss)
}

type StringableArray []Stringable

func (sa StringableArray) ToStringSlice() ([]string, error) {
	values := make([]string, len(sa))
	for i, s := range sa {
		if value, err := s.ToString(); err != nil {
			return nil, fmt.Errorf("cannot encode %q at index %d: %w", sa[i], i, err)
		} else {
			values[i] = value
		}
	}
	return values, nil
}

func (sa StringableArray) FromStringSlice(values []string) error {
	for i, s := range values {
		if err := sa[i].FromString(s); err != nil {
			return fmt.Errorf("cannot decode %q at index %d: %w", values[i], i, err)
		}
	}
	return nil
}

type SingleStringableSlicableWrapper struct{ Stringable }

func (w SingleStringableSlicableWrapper) ToStringSlice() ([]string, error) {
	if value, err := w.ToString(); err != nil {
		return nil, err
	} else {
		return []string{value}, nil
	}
}

func (w SingleStringableSlicableWrapper) FromStringSlice(values []string) error {
	if len(values) >= 1 {
		return w.FromString(values[0])
	}
	return nil
}

var (
	stringSliceableType = internal.TypeOf[StringSlicable]()
	byteType            = internal.TypeOf[byte]()
)

func isSliceType(t reflect.Type) bool {
	return t.Kind() == reflect.Slice || t.Kind() == reflect.Array
}

func isByteSliceType(t reflect.Type) bool {
	if isSliceType(t) && t.Elem() == byteType {
		return true
	}
	return false
}
