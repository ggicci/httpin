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

var decoders = make(map[reflect.Type]interface{}) // custom decoders

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
	if decoder == nil {
		panic(fmt.Errorf("httpin: %w: %q", ErrNilTypeDecoder, typ))
	}

	if !isTypeDecoder(decoder) {
		panic(fmt.Errorf("httpin: %w: %q", ErrInvalidTypeDecoder, typ))
	}

	decoders[typ] = decoder
}

// decoderOf retrieves a decoder by type.
func decoderOf(t reflect.Type) interface{} {
	dec := decoders[t]
	if dec != nil {
		return dec
	}
	return internal.DecoderOf(t)
}
