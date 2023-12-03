package core

import (
	"errors"
	"fmt"
	"reflect"
)

var defaultRegistry = newRegistry()
var fileEncoderInterface = reflect.TypeOf((*FileEncoder)(nil)).Elem()

// registry is just a place to gather all encoders and decoders together.
type registry struct {
	fileTypes map[reflect.Type]FileDecoderAdaptor
}

func newRegistry() *registry {
	r := &registry{
		fileTypes: make(map[reflect.Type]FileDecoderAdaptor),
	}
	return r
}

func (r *registry) RegisterFileType(typ reflect.Type, fd FileDecoder[any]) error {
	if r.IsFileType(typ) {
		return fmt.Errorf("duplicate file type: %v", typ)
	}
	if !typ.Implements(fileEncoderInterface) {
		return fmt.Errorf("file type must implement FileEncoder interface")
	}
	if fd == nil {
		return errors.New("file decoder cannot be nil")
	}
	r.fileTypes[typ] = AdaptDecoder(typ, fd).(FileDecoderAdaptor)
	return nil
}

func (r *registry) GetFileDecoder(typ reflect.Type) FileDecoderAdaptor {
	return r.fileTypes[typ]
}

func (r *registry) IsFileType(typ reflect.Type) bool {
	_, ok := r.fileTypes[typ]
	return ok
}

func (r *registry) RemoveFileType(typ reflect.Type) {
	delete(r.fileTypes, typ)
}

// ToPointerEncoder makes an encoder for a type (T) be able to used as an
// encoder for a T's pointer type (*T).
type ToPointerEncoder struct {
	Encoder
}

func (pe ToPointerEncoder) Encode(value reflect.Value) (string, error) {
	return pe.Encoder.Encode(value.Elem())
}

// Encoder is a type that can encode a value of type T to a string. It is
// used by the "form", "query", and "header" directives to encode a value.
type Encoder interface {
	Encode(value reflect.Value) (string, error)
}
