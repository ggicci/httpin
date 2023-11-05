package internal

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

var DefaultRegistry = NewRegistry()
var fileEncoderInterface = reflect.TypeOf((*FileEncoder)(nil)).Elem()

type Registry struct {
	builtinTypeEncoders priorityPair // builtin encoders, always registered
	builtinTypeDecoders priorityPair // builtin decoders, always registered

	typedEncoders priorityPair                 // encoders (by type)
	namedEncoders map[string]*NamedEncoderInfo // encoders (by name)
	typedDecoders priorityPair                 // decoders (by type)
	namedDecoders map[string]*NamedDecoderInfo // decoders (by name)

	fileTypes   map[reflect.Type]FileDecoderAdaptor
	bodyFormats map[string]BodyEncodeDecoder
}

type NamedEncoderInfo struct {
	Name     string
	Original Encoder
}

type NamedDecoderInfo struct {
	Name     string
	Original any
	Adapted  ValueDecoderAdaptor
}

func NewRegistry() *Registry {
	r := &Registry{
		builtinTypeEncoders: newPriorityPair(),
		builtinTypeDecoders: newPriorityPair(),

		typedEncoders: newPriorityPair(),
		namedEncoders: make(map[string]*NamedEncoderInfo),
		typedDecoders: newPriorityPair(),
		namedDecoders: make(map[string]*NamedDecoderInfo),

		fileTypes:   make(map[reflect.Type]FileDecoderAdaptor),
		bodyFormats: make(map[string]BodyEncodeDecoder),
	}

	// Always register builtin stuffs.
	r.registerBuiltinTypeEncoders()
	r.registerBuiltinTypeDecoders()
	r.registerBuiltinFileTypes()
	r.registerBuiltinBodyFormats()
	return r
}

func (r *Registry) RegisterEncoder(typ reflect.Type, encoder Encoder, force ...bool) error {
	return r.registerTypedEncoderTo(r.typedEncoders, typ, encoder, len(force) > 0 && force[0])
}

func (r *Registry) RegisterNamedEncoder(name string, encoder Encoder, force ...bool) error {
	ignoreConflict := len(force) > 0 && force[0]
	if _, ok := r.namedEncoders[name]; ok && !ignoreConflict {
		return fmt.Errorf("duplicate name: %q", name)
	}
	if err := validateEncoder(encoder); err != nil {
		return err
	}

	r.namedEncoders[name] = &NamedEncoderInfo{
		Name:     name,
		Original: encoder,
	}
	return nil
}

func (r *Registry) GetEncoder(typ reflect.Type) Encoder {
	if e := r.typedEncoders.GetOne(typ); e != nil {
		return e.(Encoder)
	}
	if e := r.builtinTypeEncoders.GetOne(typ); e != nil {
		return e.(Encoder)
	}
	return nil
}

func (r *Registry) GetNamedEncoder(name string) *NamedEncoderInfo {
	return r.namedEncoders[name]
}

func (r *Registry) RemoveEncoder(typ reflect.Type) {
	delete(r.typedEncoders, typ)
}

func (r *Registry) RemoveNamedEncoder(name string) {
	delete(r.namedEncoders, name)
}

func (r *Registry) RegisterDecoder(typ reflect.Type, decoder Decoder[any], force ...bool) error {
	return r.registerTypedDecoderTo(r.typedDecoders, typ, decoder, len(force) > 0 && force[0])
}

func (r *Registry) RegisterNamedDecoder(name string, typ reflect.Type, decoder Decoder[any], force ...bool) error {
	ignoreConflict := len(force) > 0 && force[0]
	if _, ok := r.namedDecoders[name]; ok && !ignoreConflict {
		return fmt.Errorf("duplicate name: %q", name)
	}
	if err := validateDecoder(decoder); err != nil {
		return err
	}
	r.namedDecoders[name] = &NamedDecoderInfo{
		Name:     name,
		Original: decoder,
		Adapted:  AdaptDecoder(typ, NewSmartDecoder(typ, ToAnyDecoder(decoder))).(ValueDecoderAdaptor),
	}
	return nil
}

func (r *Registry) GetDecoder(typ reflect.Type) ValueDecoderAdaptor {
	if d := r.typedDecoders.GetOne(typ); d != nil {
		return d.(ValueDecoderAdaptor)
	}
	if d := r.builtinTypeDecoders.GetOne(typ); d != nil {
		return d.(ValueDecoderAdaptor)
	}
	return nil
}

func (r *Registry) GetNamedDecoder(name string) *NamedDecoderInfo {
	return r.namedDecoders[name]
}

func (r *Registry) RemoveDecoder(typ reflect.Type) {
	delete(r.typedDecoders, typ)
}

func (r *Registry) RemoveNamedDecoder(name string) {
	delete(r.namedDecoders, name)
}

func (r *Registry) RegisterFileType(typ reflect.Type, fd FileDecoder[any]) error {
	if r.IsFileType(typ) {
		return fmt.Errorf("duplicate file type: %v", typ)
	}
	if !typ.Implements(fileEncoderInterface) {
		return fmt.Errorf("file type must implement FileEncoder interface")
	}
	if fd == nil {
		return errors.New("file decoder cannot be nil")
	}
	r.fileTypes[typ] = AdaptDecoder(typ, fd).(FileDecoderAdaptor) // FIXME: check fd type
	return nil
}

func (r *Registry) RegisterBodyFormat(format string, body BodyEncodeDecoder, force ...bool) error {
	ignoreConflict := len(force) > 0 && force[0]
	format = strings.ToLower(format)
	if _, ok := r.bodyFormats[format]; ok && !ignoreConflict {
		return fmt.Errorf("duplicate body format: %q", format)
	}
	if format == "" {
		return errors.New("body format cannot be empty")
	}
	if body == nil {
		return errors.New("body encoder/decoder cannot be nil")
	}
	r.bodyFormats[format] = body
	return nil
}

func (r *Registry) GetBodyDecoder(format string) BodyEncodeDecoder {
	return r.bodyFormats[format]
}

func (r *Registry) GetFileDecoder(typ reflect.Type) FileDecoderAdaptor {
	return r.fileTypes[typ]
}

func (r *Registry) IsFileType(typ reflect.Type) bool {
	return r.GetFileDecoder(typ) != nil
}

func (r *Registry) RemoveFileType(typ reflect.Type) {
	delete(r.fileTypes, typ)
}

func (r *Registry) registerBuiltinTypeEncoders() {
	for typ, encoder := range theBuiltinEncoders {
		r.registerTypedEncoderTo(r.builtinTypeEncoders, typ, encoder.(Encoder), false)
	}
}

func (r *Registry) registerBuiltinTypeDecoders() {
	for typ, decoder := range theBuiltinDecoders {
		r.registerTypedDecoderTo(r.builtinTypeDecoders, typ, decoder, false)
	}
}

func (r *Registry) registerBuiltinFileTypes() {
	r.RegisterFileType(
		reflect.TypeOf((*File)(nil)),
		ToAnyFileDecoder[*File](FileDecoderFunc[*File](decodeFile)),
	)
}

func (r *Registry) registerBuiltinBodyFormats() {
	r.RegisterBodyFormat("json", &JSONBody{})
	r.RegisterBodyFormat("xml", &XMLBody{})
}

func (r *Registry) registerTypedEncoderTo(p priorityPair, typ reflect.Type, encoder Encoder, force bool) error {
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

func (r *Registry) registerTypedDecoderTo(p priorityPair, typ reflect.Type, decoder Decoder[any], force bool) error {
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
