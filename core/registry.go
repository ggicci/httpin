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
	builtinTypeEncoders internal.PriorityPair // builtin encoders, always registered
	builtinTypeDecoders internal.PriorityPair // builtin decoders, always registered

	typedEncoders internal.PriorityPair        // encoders (by type)
	namedEncoders map[string]*namedEncoderInfo // encoders (by name)
	typedDecoders internal.PriorityPair        // decoders (by type)
	namedDecoders map[string]*namedDecoderInfo // decoders (by name)

	fileTypes map[reflect.Type]FileDecoderAdaptor
}

type namedEncoderInfo struct {
	Name     string
	Original Encoder
}

type namedDecoderInfo struct {
	Name     string
	Original any
	Adapted  ValueDecoderAdaptor
}

func newRegistry() *registry {
	r := &registry{
		builtinTypeEncoders: internal.NewPriorityPair(),
		builtinTypeDecoders: internal.NewPriorityPair(),

		typedEncoders: internal.NewPriorityPair(),
		namedEncoders: make(map[string]*namedEncoderInfo),
		typedDecoders: internal.NewPriorityPair(),
		namedDecoders: make(map[string]*namedDecoderInfo),

		fileTypes: make(map[reflect.Type]FileDecoderAdaptor),
	}

	// Always register builtin stuffs.
	r.registerBuiltinTypeEncoders()
	r.registerBuiltinTypeDecoders()
	return r
}

func (r *registry) RegisterEncoder(typ reflect.Type, encoder Encoder, force ...bool) error {
	return r.registerTypedEncoderTo(r.typedEncoders, typ, encoder, len(force) > 0 && force[0])
}

func (r *registry) RegisterNamedEncoder(name string, encoder Encoder, force ...bool) error {
	ignoreConflict := len(force) > 0 && force[0]
	if _, ok := r.namedEncoders[name]; ok && !ignoreConflict {
		return fmt.Errorf("duplicate name: %q", name)
	}
	if err := validateEncoder(encoder); err != nil {
		return err
	}

	r.namedEncoders[name] = &namedEncoderInfo{
		Name:     name,
		Original: encoder,
	}
	return nil
}

func (r *registry) GetEncoder(typ reflect.Type) Encoder {
	if e := r.typedEncoders.GetOne(typ); e != nil {
		return e.(Encoder)
	}
	if e := r.builtinTypeEncoders.GetOne(typ); e != nil {
		return e.(Encoder)
	}
	return nil
}

func (r *registry) GetNamedEncoder(name string) *namedEncoderInfo {
	return r.namedEncoders[name]
}

func (r *registry) RemoveEncoder(typ reflect.Type) {
	delete(r.typedEncoders, typ)
}

func (r *registry) RemoveNamedEncoder(name string) {
	delete(r.namedEncoders, name)
}

func (r *registry) RegisterDecoder(typ reflect.Type, decoder Decoder[any], force ...bool) error {
	return r.registerTypedDecoderTo(r.typedDecoders, typ, decoder, len(force) > 0 && force[0])
}

func (r *registry) RegisterNamedDecoder(name string, typ reflect.Type, decoder Decoder[any], force ...bool) error {
	ignoreConflict := len(force) > 0 && force[0]
	if _, ok := r.namedDecoders[name]; ok && !ignoreConflict {
		return fmt.Errorf("duplicate name: %q", name)
	}
	if err := validateDecoder(decoder); err != nil {
		return err
	}
	r.namedDecoders[name] = &namedDecoderInfo{
		Name:     name,
		Original: decoder,
		Adapted:  AdaptDecoder(typ, NewSmartDecoder(typ, ToAnyDecoder(decoder))).(ValueDecoderAdaptor),
	}
	return nil
}

func (r *registry) GetDecoder(typ reflect.Type) ValueDecoderAdaptor {
	if d := r.typedDecoders.GetOne(typ); d != nil {
		return d.(ValueDecoderAdaptor)
	}
	if d := r.builtinTypeDecoders.GetOne(typ); d != nil {
		return d.(ValueDecoderAdaptor)
	}
	return nil
}

func (r *registry) GetNamedDecoder(name string) *namedDecoderInfo {
	return r.namedDecoders[name]
}

func (r *registry) RemoveDecoder(typ reflect.Type) {
	delete(r.typedDecoders, typ)
}

func (r *registry) RemoveNamedDecoder(name string) {
	delete(r.namedDecoders, name)
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

func (r *registry) registerBuiltinTypeEncoders() {
	for typ, encoder := range theBuiltinEncoders {
		r.registerTypedEncoderTo(r.builtinTypeEncoders, typ, encoder.(Encoder), false)
	}
}

func (r *registry) registerBuiltinTypeDecoders() {
	for typ, decoder := range theBuiltinDecoders {
		r.registerTypedDecoderTo(r.builtinTypeDecoders, typ, decoder, false)
	}
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
