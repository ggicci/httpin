package core

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/ggicci/httpin/internal"
)

var defaultRegistry = newRegistry()
var fileEncoderInterface = reflect.TypeOf((*FileEncoder)(nil)).Elem()

// registry is just a place to gather all encoders and decoders together.
type registry struct {
	fileTypes map[reflect.Type]FileDecoderAdaptor
}

type namedDecoderInfo struct {
	Name     string
	Original any
	Adapted  ValueDecoderAdaptor
	Adapt    AnyStringableAdaptor
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

func (r *registry) registerTypedEncoderTo(p internal.PriorityPair, typ reflect.Type, encoder Encoder, force bool) error {
	if err := validateEncoder(encoder); err != nil {
		return err
	}

	if err := p.SetPair(typ, encoder, nil, force); err != nil {
		return err
	}

	if typ.Kind() != reflect.Pointer {
		// When we have a non-pointer type (T), we also register the encoder for its
		// pointer type (*T). The encoder for the pointer type (*T) will be registered
		// as the secondary encoder.
		if err := p.SetPair(reflect.PtrTo(typ), nil, ToPointerEncoder{encoder}, force); err != nil {
			return err
		}
	}

	return nil
}

func (r *registry) registerTypedDecoderTo(p internal.PriorityPair, typ reflect.Type, decoder Decoder[any], force bool) error {
	if err := validateDecoder(decoder); err != nil {
		return err
	}

	primaryDecoder := AdaptDecoder(typ, NewSmartDecoder(typ, decoder))
	if err := p.SetPair(typ, primaryDecoder, nil, force); err != nil {
		return err
	}

	if typ.Kind() == reflect.Pointer {
		// When we have a pointer type (*T), we also register the decoder for its base
		// type (T). The decoder for the base type (T) will be registered as the
		// secondary decoder.
		baseType := typ.Elem()
		secondaryDecoder := AdaptDecoder(baseType, NewSmartDecoder(baseType, decoder))
		return p.SetPair(baseType, nil, secondaryDecoder, force)
	} else {
		// When we have a non-pointer type (T), we also register the decoder for its
		// pointer type (*T). The decoder for the pointer type (*T) will be registered
		// as the secondary decoder.
		pointerType := reflect.PtrTo(typ)
		secondaryDecoder := AdaptDecoder(pointerType, NewSmartDecoder(pointerType, decoder))
		return p.SetPair(pointerType, nil, secondaryDecoder, force)
	}
}

func validateEncoder(encoder any) error {
	if encoder == nil || internal.IsNil(reflect.ValueOf(encoder)) {
		return errors.New("nil encoder")
	}
	return nil
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
