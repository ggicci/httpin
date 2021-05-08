package httpin

import (
	"reflect"

	"github.com/ggicci/httpin/internal"
)

type Decoder internal.Decoder

var decoders = map[reflect.Type]Decoder{} // custom decoders

func decoderOf(t reflect.Type) Decoder {
	dec := decoders[t]
	if dec != nil {
		return dec
	}
	return internal.DecoderOf(t)
}
