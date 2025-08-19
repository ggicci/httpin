package codec

import (
	"reflect"

	"github.com/ggicci/strconvx"
)

type (
	StringCodec        = strconvx.StringCodec
	StringCodecAdaptor = strconvx.AnyAdaptor
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
