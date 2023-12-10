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
		return NewStringSlicablePatchFieldWrapper(rv, adapt)
	}

	if isSliceType(rv.Type()) && !isByteSliceType(rv.Type()) {
		return NewStringableSliceWrapper(rv, adapt)
	} else {
		return NewStringSlicableSingleStringableWrapper(rv, adapt)
	}
}

// StringSlicablePatchFieldWrapper wraps a patch.Field[T] to implement
// StringSlicable. The wrapped reflect.Value must be a patch.Field[T].
//
// It works like a proxy. It delegates the ToStringSlice and FromStringSlice
// calls to the internal StringSlicable.
type StringSlicablePatchFieldWrapper struct {
	Value                   reflect.Value // of patch.Field[T]
	internalStringSliceable StringSlicable
}

// NewStringSlicablePatchFieldWrapper creates a StringSlicablePatchFieldWrapper from rv.
// Returns error when patch.Field.Value is not a StringSlicable.
func NewStringSlicablePatchFieldWrapper(rv reflect.Value, adapt AnyStringableAdaptor) (StringSlicablePatchFieldWrapper, error) {
	stringSlicable, err := NewStringSlicable(rv.FieldByName("Value"), adapt)
	if err != nil {
		return StringSlicablePatchFieldWrapper{}, fmt.Errorf("cannot create StringSlicable for PatchField: %w", err)
	} else {
		return StringSlicablePatchFieldWrapper{
			Value:                   rv,
			internalStringSliceable: stringSlicable,
		}, nil
	}
}

func (w StringSlicablePatchFieldWrapper) ToStringSlice() ([]string, error) {
	if w.Value.FieldByName("Valid").Bool() {
		return w.internalStringSliceable.ToStringSlice()
	} else {
		return []string{}, nil
	}
}

func (w StringSlicablePatchFieldWrapper) FromStringSlice(values []string) error {
	if err := w.internalStringSliceable.FromStringSlice(values); err != nil {
		return err
	} else {
		w.Value.FieldByName("Valid").SetBool(true)
		return nil
	}
}

type StringableSlice []Stringable

func (sa StringableSlice) ToStringSlice() ([]string, error) {
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

func (sa StringableSlice) FromStringSlice(values []string) error {
	for i, s := range values {
		if err := sa[i].FromString(s); err != nil {
			return fmt.Errorf("cannot decode %q at index %d: %w", values[i], i, err)
		}
	}
	return nil
}

// StringableSliceWrapper wraps a reflect.Value to implement StringSlicable. The
// wrapped reflect.Value must be a slice of Stringable.
type StringableSliceWrapper struct {
	Value reflect.Value
	Adapt AnyStringableAdaptor
}

// NewStringableSliceWrapper creates a StringableSliceWrapper from rv.
// Returns error when rv is not a slice of Stringable or cannot get address of rv.
func NewStringableSliceWrapper(rv reflect.Value, adapt AnyStringableAdaptor) (StringableSliceWrapper, error) {
	if !rv.CanAddr() {
		return StringableSliceWrapper{}, fmt.Errorf("cannot get address of value %q", rv)
	}
	return StringableSliceWrapper{Value: rv, Adapt: adapt}, nil
}

func (w StringableSliceWrapper) ToStringSlice() ([]string, error) {
	var stringables = make(StringableSlice, w.Value.Len())
	for i := 0; i < w.Value.Len(); i++ {
		if stringable, err := NewStringable(w.Value.Index(i), w.Adapt); err != nil {
			return nil, fmt.Errorf("cannot create Stringable from %q at index %d: %w", w.Value.Index(i), i, err)
		} else {
			stringables[i] = stringable
		}
	}
	return stringables.ToStringSlice()
}

func (w StringableSliceWrapper) FromStringSlice(ss []string) error {
	var stringables = make(StringableSlice, len(ss))
	w.Value.Set(reflect.MakeSlice(w.Value.Type(), len(ss), len(ss)))
	for i := range ss {
		if stringable, err := NewStringable(w.Value.Index(i), w.Adapt); err != nil {
			return fmt.Errorf("cannot create Stringable from %q at index %d: %w", w.Value.Index(i), i, err)
		} else {
			stringables[i] = stringable
		}
	}
	return stringables.FromStringSlice(ss)
}

// StringSlicableSingleStringableWrapper wraps a reflect.Value to implement
// StringSlicable. The wrapped reflect.Value must be a Stringable.
type StringSlicableSingleStringableWrapper struct{ Stringable }

func NewStringSlicableSingleStringableWrapper(rv reflect.Value, adapt AnyStringableAdaptor) (StringSlicable, error) {
	if stringable, err := NewStringable(rv, adapt); err != nil {
		return nil, err
	} else {
		return StringSlicableSingleStringableWrapper{stringable}, nil
	}
}

func (w StringSlicableSingleStringableWrapper) ToStringSlice() ([]string, error) {
	if value, err := w.ToString(); err != nil {
		return nil, err
	} else {
		return []string{value}, nil
	}
}

func (w StringSlicableSingleStringableWrapper) FromStringSlice(values []string) error {
	if len(values) > 0 {
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
