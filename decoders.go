package httpin

import (
	"reflect"

	"github.com/ggicci/httpin/internal"
)

// Decoder is the interface implemented by types that can decode bytes to
// themselves.
type Decoder internal.Decoder

// DecoderFunc is an adaptor to allow the use of ordinary functions as httpin
// decoders.
type DecoderFunc internal.DecoderFunc

var decoders = map[reflect.Type]Decoder{} // custom decoders

// decoderOf retrieves a decoder by type.
func decoderOf(t reflect.Type) Decoder {
	dec := decoders[t]
	if dec != nil {
		return dec
	}
	return internal.DecoderOf(t)
}
