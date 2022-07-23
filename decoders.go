package httpin

import (
	"fmt"
	"reflect"

	"github.com/ggicci/httpin/internal"
)

// ValueTypeDecoder is the interface implemented by types that can decode a
// string to themselves.
type ValueTypeDecoder = internal.ValueTypeDecoder

// FileTypeDecoder is the interface implemented by types that can decode a
// *multipart.FileHeader to themselves.
type FileTypeDecoder = internal.FileTypeDecoder

// ValueTypeDecoderFunc is an adaptor to allow the use of ordinary functions as
// httpin `ValueTypeDecoder`s.
type ValueTypeDecoderFunc = internal.ValueTypeDecoderFunc

// FileTypeDecoderFunc is an adaptor to allow the use of ordinary functions as
// httpin `FileTypeDecoder`s.
type FileTypeDecoderFunc = internal.FileTypeDecoderFunc

var (
	decoders      = make(map[reflect.Type]interface{}) // custom decoders
	namedDecoders = make(map[string]interface{})       // custom decoders (registered by name)
)

func isTypeDecoder(decoder interface{}) bool {
	_, isValueTypeDecoder := decoder.(ValueTypeDecoder)
	_, isFileTypeDecoder := decoder.(FileTypeDecoder)
	return isValueTypeDecoder || isFileTypeDecoder
}

// RegisterTypeDecoder registers a specific type decoder. Panics on conflicts.
func RegisterTypeDecoder(typ reflect.Type, decoder interface{}) {
	if _, ok := decoders[typ]; ok {
		panic(fmt.Errorf("httpin: %w: %q", ErrDuplicateTypeDecoder, typ))
	}

	ReplaceTypeDecoder(typ, decoder)
}

// ReplaceTypeDecoder replaces a specific type decoder.
func ReplaceTypeDecoder(typ reflect.Type, decoder interface{}) {
	ensureValidDecoder(typ, decoder)
	decoders[typ] = decoder
}

// RegisterNamedDecoder registers a decoder by name. Panics on conflicts.
func RegisterNamedDecoder(name string, decoder interface{}) {
	if _, ok := namedDecoders[name]; ok {
		panic(fmt.Errorf("httpin: %w: %q", ErrDuplicateNamedDecoder, name))
	}

	ReplaceNamedDecoder(name, decoder)
}

// ReplaceNamedDecoder replaces a decoder by name.
func ReplaceNamedDecoder(name string, decoder interface{}) {
	ensureValidDecoder(nil, decoder)
	namedDecoders[name] = decoder
}

func ensureValidDecoder(typ reflect.Type, decoder interface{}) {
	if decoder == nil {
		panic(fmt.Errorf("httpin: %w: %q", ErrNilTypeDecoder, typ))
	}

	if !isTypeDecoder(decoder) {
		panic(fmt.Errorf("httpin: %w: %q", ErrInvalidTypeDecoder, typ))
	}
}

// decoderOf retrieves a decoder by type, from the global registerred decoders.
func decoderOf(t reflect.Type) interface{} {
	dec := decoders[t]
	if dec != nil {
		return dec
	}
	return internal.DecoderOf(t)
}

// decoderByName retrieves a decoder by name, from the global registerred named decoders.
func decoderByName(name string) interface{} {
	return namedDecoders[name]
}
