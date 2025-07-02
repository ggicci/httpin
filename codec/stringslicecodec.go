package codec

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/ggicci/httpin/internal"
)

// StringSliceCodec is implemented by types that can be converted to and from a slice of strings.
// It defines methods for parsing from []string and serializing back to []string.
type StringSliceCodec interface {
	ToStringSlice() ([]string, error)
	FromStringSlice([]string) error
}

// NewStringSliceCodec creates a StringSliceCodec instance, it allows adapting
// the underlying StringCodec by passing through a custom StringCodecAdaptor.
func (ns *Namespace) NewStringSliceCodec(rv reflect.Value, adaptor StringCodecAdaptor) (StringSliceCodec, error) {
	if rv.Type().Implements(stringSliceCodecType) && rv.CanInterface() {
		return rv.Interface().(StringSliceCodec), nil
	}

	// When rv is of type patch.Field[T].
	if IsPatchField(rv.Type()) {
		return ns.NewStringSliceCodec4PatchField(rv, adaptor)
	}

	if isSliceType(rv.Type()) && !isByteSliceType(rv.Type()) {
		// When rv is of []T or [N]T but no []byte or [N]byte.
		return ns.NewStringSliceCodec4Slice(rv, adaptor)
	} else {
		// When rv is of []byte or [N]byte.
		return ns.NewStringSliceCodec4SingleStringCodec(rv, adaptor)
	}
}

// StringSliceCodec4PatchField makes patch.Field[T] implement StringSliceCodec
// as long as T implements StringCodec. It is used to eliminate the effort of
// implementing StringSliceCodec for patch.Field[T] for every type T.
type StringSliceCodec4PatchField struct {
	Value reflect.Value // of patch.Field[T]
	codec StringSliceCodec
}

func (ns *Namespace) NewStringSliceCodec4PatchField(rv reflect.Value, adaptor StringCodecAdaptor) (*StringSliceCodec4PatchField, error) {
	codec, err := ns.NewStringSliceCodec(rv.FieldByName("Value"), adaptor)
	if err != nil {
		return nil, err
	} else {
		return &StringSliceCodec4PatchField{
			Value: rv,
			codec: codec,
		}, nil
	}
}

func (ssc *StringSliceCodec4PatchField) ToStringSlice() ([]string, error) {
	if ssc.Value.FieldByName("Valid").Bool() {
		return ssc.codec.ToStringSlice()
	} else {
		return []string{}, nil
	}
}

func (ssc *StringSliceCodec4PatchField) FromStringSlice(values []string) error {
	if err := ssc.codec.FromStringSlice(values); err != nil {
		return err
	} else {
		ssc.Value.FieldByName("Valid").SetBool(true)
		return nil
	}
}

// StringCodecs is a slice of StringCodec, which implements StringSliceCodec.
type StringCodecs []StringCodec

func (sc StringCodecs) ToStringSlice() ([]string, error) {
	values := make([]string, len(sc))
	for i, s := range sc {
		if value, err := s.ToString(); err != nil {
			return nil, fmt.Errorf("cannot stringify %v at index %d: %w", s, i, err)
		} else {
			values[i] = value
		}
	}
	return values, nil
}

func (sc StringCodecs) FromStringSlice(values []string) error {
	for i, s := range values {
		if err := sc[i].FromString(s); err != nil {
			return fmt.Errorf("cannot convert from string %q at index %d: %w", s, i, err)
		}
	}
	return nil
}

// StringSliceCodec4Slice makes []T or [T] implement StringSliceCodec as long as
// T implements StringCodec. It is used to eliminate the effort of implementing
// StringSliceCodec for []T or [T] for every type T.
type StringSliceCodec4Slice struct {
	Value   reflect.Value // of []T or [T]
	Adaptor StringCodecAdaptor
	ns      *Namespace // Namespace to create StringCodec instances
}

func (ns *Namespace) NewStringSliceCodec4Slice(rv reflect.Value, adaptor StringCodecAdaptor) (*StringSliceCodec4Slice, error) {
	if !rv.CanAddr() {
		return nil, errors.New("unaddressable value")
	}
	return &StringSliceCodec4Slice{Value: rv, Adaptor: adaptor, ns: ns}, nil
}

func (ssc *StringSliceCodec4Slice) ToStringSlice() ([]string, error) {
	var codecs = make(StringCodecs, ssc.Value.Len())
	for i := 0; i < ssc.Value.Len(); i++ {
		if codec, err := ssc.ns.NewStringCodec(ssc.Value.Index(i), ssc.Adaptor); err != nil {
			return nil, fmt.Errorf("cannot create StringCodec from %q at index %d: %w", ssc.Value.Index(i), i, err)
		} else {
			codecs[i] = codec
		}
	}
	return codecs.ToStringSlice()
}

func (ssc *StringSliceCodec4Slice) FromStringSlice(ss []string) error {
	var codecs = make(StringCodecs, len(ss))
	ssc.Value.Set(reflect.MakeSlice(ssc.Value.Type(), len(ss), len(ss)))
	for i := range ss {
		if codec, err := ssc.ns.NewStringCodec(ssc.Value.Index(i), ssc.Adaptor); err != nil {
			return fmt.Errorf("cannot create StringCodec at index %d: %w", i, err)
		} else {
			codecs[i] = codec
		}
	}
	return codecs.FromStringSlice(ss)
}

// StringSliceCodec4SingleStringCodec turns a single StringCodec into a StringSliceCodec.
type StringSliceCodec4SingleStringCodec struct{ StringCodec }

func (ns *Namespace) NewStringSliceCodec4SingleStringCodec(
	rv reflect.Value, adaptor StringCodecAdaptor) (*StringSliceCodec4SingleStringCodec, error) {
	if codec, err := ns.NewStringCodec(rv, adaptor); err != nil {
		return nil, err
	} else {
		return &StringSliceCodec4SingleStringCodec{codec}, nil
	}
}

func (ssc *StringSliceCodec4SingleStringCodec) ToStringSlice() ([]string, error) {
	if value, err := ssc.ToString(); err != nil {
		return nil, err
	} else {
		return []string{value}, nil
	}
}

func (ssc *StringSliceCodec4SingleStringCodec) FromStringSlice(values []string) error {
	if len(values) > 0 {
		return ssc.FromString(values[0])
	}
	return nil
}

var (
	stringSliceCodecType = internal.TypeOf[StringSliceCodec]()
	byteType             = internal.TypeOf[byte]()
)

// isSliceType checks if the given type is a slice or an array (i.e., []T or [N]T).
func isSliceType(t reflect.Type) bool {
	return t.Kind() == reflect.Slice || t.Kind() == reflect.Array
}

// isByteSliceType checks if the given type is a slice of bytes (i.e., []byte or [N]byte).
func isByteSliceType(t reflect.Type) bool {
	if isSliceType(t) && t.Elem() == byteType {
		return true
	}
	return false
}
