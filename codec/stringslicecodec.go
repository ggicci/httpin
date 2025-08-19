package codec

import (
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
