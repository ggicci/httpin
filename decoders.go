package httpin

import (
	"fmt"
	"reflect"

	"github.com/ggicci/httpin/internal"
)

// TypeDecoder is the interface implemented by types that can decode bytes to
// themselves.
type TypeDecoder = internal.TypeDecoder

// TypeDecoderFunc is an adaptor to allow the use of ordinary functions as httpin
// TypeDecoders.
type TypeDecoderFunc = internal.TypeDecoderFunc

var decoders = map[reflect.Type]TypeDecoder{} // custom decoders

// RegisterTypeDecoder registers a specific type decoder. Panics on conflicts.
func RegisterTypeDecoder(typ reflect.Type, decoder TypeDecoder) {
	if _, ok := decoders[typ]; ok {
		panic(fmt.Errorf("%w: %q", ErrDuplicateTypeDecoder, typ))
	}
	ReplaceTypeDecoder(typ, decoder)
}

// ReplaceTypeDecoder replaces a specific type decoder.
func ReplaceTypeDecoder(typ reflect.Type, decoder TypeDecoder) {
	if decoder == nil {
		panic(fmt.Errorf("%w: %q", typ, ErrNilTypeDecoder))
	}
	decoders[typ] = decoder
}

// decoderOf retrieves a decoder by type.
func decoderOf(t reflect.Type) TypeDecoder {
	dec := decoders[t]
	if dec != nil {
		return dec
	}
	return internal.DecoderOf(t)
}
