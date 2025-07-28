package codec

import (
	"errors"
	"fmt"
	"reflect"
)

// NewStringCodec creates a StringCodec instance from a reflect.Value. It allows
// overriding the underlying StringCodec behaviour by passing through a
// StringCodecAdaptor.
func (ns *Namespace) NewStringCodec(rv reflect.Value, adaptor StringCodecAdaptor) (codec StringCodec, err error) {
	if IsPatchField(rv.Type()) {
		codec, err = ns.NewStringCodec4PatchField(rv, adaptor)
	} else {
		codec, err = ns.newStringCodec(rv, adaptor)
	}
	if err != nil {
		return nil, err
	}
	return codec, nil
}

// Create a StringCodec from a reflect.Value. If rv is a pointer type, it will
// try to create a StringCodec from rv. Otherwise, it will try to create a
// StringCodec from rv.Addr(). Only basic built-in types are supported. As a
// special case, time.Time is also supported. For more details, see
// github.com/ggicci/strconvx package.
func (ns *Namespace) newStringCodec(rv reflect.Value, adaptor StringCodecAdaptor) (StringCodec, error) {
	rv, err := getPointer(rv)
	if err != nil {
		return nil, err
	}

	// Now rv is a pointer type.
	if adaptor != nil {
		return adaptor(rv.Interface())
	}

	// Fallback to use built-in StringCodec types.
	return ns.Namespace.New(rv)
}

// StringCodec4PatchField makes patch.Field[T] implement StringCodec as long as
// T implements StringCodec. It is used to eliminate the effort of implementing
// StringCodec for patch.Field[T] for every type T.
type StringCodec4PatchField struct {
	Value reflect.Value // of patch.Field[T]
	codec StringCodec
}

func (ns *Namespace) NewStringCodec4PatchField(rv reflect.Value, adapt StringCodecAdaptor) (*StringCodec4PatchField, error) {
	StringCodec, err := ns.NewStringCodec(rv.FieldByName("Value"), adapt)
	if err != nil {
		return &StringCodec4PatchField{}, fmt.Errorf("cannot create StringCodec for PatchField: %w", err)
	} else {
		return &StringCodec4PatchField{
			Value: rv,
			codec: StringCodec,
		}, nil
	}
}

func (w *StringCodec4PatchField) ToString() (string, error) {
	if w.Value.FieldByName("Valid").Bool() {
		return w.codec.ToString()
	} else {
		return "", errors.New("invalid value") // when Valid is false
	}
}

// FromString sets the value of the wrapped patch.Field[T] from the given
// string. It returns an error if the given string is not valid. And leaves the
// original value of both Value and Valid unchanged. On the other hand, if no
// error occurs, it sets Valid to true.
func (w *StringCodec4PatchField) FromString(s string) error {
	if err := w.codec.FromString(s); err != nil {
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
		return rv, fmt.Errorf("cannot get address of value %v", rv)
	}
	rv = rv.Addr()
	return rv, nil
}
