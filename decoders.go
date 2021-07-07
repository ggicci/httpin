package httpin

import (
	"fmt"
	"reflect"

	"github.com/ggicci/httpin/internal"
)

// Decoder is the interface implemented by types that can decode bytes to
// themselves.
type Decoder = internal.Decoder

// DecoderFunc is an adaptor to allow the use of ordinary functions as httpin
// decoders.
type DecoderFunc = internal.DecoderFunc

var decoders = map[reflect.Type]Decoder{} // custom decoders

// RegisterDecoder registers a decoder to decode the specific type. Panics on conflicts.
func RegisterDecoder(typ reflect.Type, decoder Decoder) {
	if _, ok := decoders[typ]; ok {
		panic(fmt.Sprintf("duplicate decoder for type %q", typ))
	}
	ReplaceDecoder(typ, decoder)
}

// ReplaceDecoder replaces a decoder to decode the specific type.
func ReplaceDecoder(typ reflect.Type, decoder Decoder) {
	if decoder == nil {
		panic("nil decoder")
	}
	decoders[typ] = decoder
}

// decoderOf retrieves a decoder by type.
func decoderOf(t reflect.Type) Decoder {
	dec := decoders[t]
	if dec != nil {
		return dec
	}
	return internal.DecoderOf(t)
}
