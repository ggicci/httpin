package core

import (
	"encoding"
	"errors"
	"reflect"

	"github.com/ggicci/httpin/internal"
)

type HybridCoder struct {
	internal.StringMarshaler
	internal.StringUnmarshaler
}

func (c *HybridCoder) ToString() (string, error) {
	if c.StringMarshaler != nil {
		return c.StringMarshaler.ToString()
	}
	return "", errors.New("StringMarshaler not implemented")
}

func (c *HybridCoder) FromString(s string) error {
	if c.StringUnmarshaler != nil {
		return c.StringUnmarshaler.FromString(s)
	}
	return errors.New("StringUnmarshaler not implemented")
}

// Hybridize a reflect.Value to a Stringable if possible.
func hybridizeCoder(rv reflect.Value) Stringable {
	if !rv.CanInterface() {
		return nil
	}

	coder := &HybridCoder{}

	// Interface: StringMarshaler.
	if rv.Type().Implements(stringMarshalerType) {
		coder.StringMarshaler = rv.Interface().(internal.StringMarshaler)
	} else if rv.Type().Implements(textMarshalerType) {
		coder.StringMarshaler = &textMarshalerWrapper{rv.Interface().(encoding.TextMarshaler), nil}
	}

	// Interface: StringUnmarshaler.
	if rv.Type().Implements(stringUnmarshalerType) {
		coder.StringUnmarshaler = rv.Interface().(internal.StringUnmarshaler)
	} else if rv.Type().Implements(textUnmarshalerType) {
		coder.StringUnmarshaler = &textMarshalerWrapper{nil, rv.Interface().(encoding.TextUnmarshaler)}
	}

	if coder.StringMarshaler == nil && coder.StringUnmarshaler == nil {
		return nil
	}

	return coder
}

type textMarshalerWrapper struct {
	encoding.TextMarshaler
	encoding.TextUnmarshaler
}

func (w textMarshalerWrapper) ToString() (string, error) {
	b, err := w.TextMarshaler.MarshalText()
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (w textMarshalerWrapper) FromString(s string) error {
	return w.TextUnmarshaler.UnmarshalText([]byte(s))
}

var (
	stringMarshalerType   = internal.TypeOf[internal.StringMarshaler]()
	stringUnmarshalerType = internal.TypeOf[internal.StringUnmarshaler]()
	textMarshalerType     = internal.TypeOf[encoding.TextMarshaler]()
	textUnmarshalerType   = internal.TypeOf[encoding.TextUnmarshaler]()
)
