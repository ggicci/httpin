package core

import (
	"reflect"

	"github.com/ggicci/httpin/codec"
	"github.com/ggicci/httpin/internal"
	"github.com/ggicci/strconvx"
)

type (
	// StringCodec is implemented by types that support bidirectional conversion
	// between their value and a string representation.
	StringCodec = codec.StringCodec

	// StringSliceCodec is implemented by types that support bidirectional conversion
	// between their value and a []string.
	StringSliceCodec = codec.StringSliceCodec

	// FileCodec is implemented by types that support bidirectional conversion
	// between their value and a file representation. This is used for file uploads.
	FileCodec = codec.FileCodec

	// FileSliceCodec is implemented by types that support bidirectional conversion
	// between their value and a []FileCodec. This is used for file uploads where
	// multiple files can be uploaded at once.
	FileSliceCodec = codec.FileSliceCodec

	StringCodecAdaptor = codec.StringCodecAdaptor
)

var (
	fileTypes                = make(map[reflect.Type]struct{})
	namedStringCodecAdaptors = make(map[string]*NamedStringCodecAdaptor)
)

// RegisterCodec registers a custom codec for the given type T. When a field of
// type T is encountered, this codec will be used to convert the value to a
// StringCodec, which will be used to convert the value from/to string.
//
// NOTE: this function is designed to override the default StringCodec adaptors
// that are registered by this package. For example, if you want to override the
// defualt behaviour of converting a bool value from/to string, you can do this:
//
//	func init() {
//		core.RegisterCodec[bool](func(b *bool) (core.StringCodec, error) {
//			return (*YesNo)(b), nil
//		})
//	}
//
//	type YesNo bool
//
//	func (yn YesNo) String() string {
//		if yn {
//			return "yes"
//		}
//		return "no"
//	}
//
//	func (yn *YesNo) FromString(s string) error {
//		switch s {
//		case "yes":
//			*yn = true
//		case "no":
//			*yn = false
//		default:
//			return fmt.Errorf("invalid YesNo value: %q", s)
//		}
//		return nil
//	}
func RegisterCodec[T any](adaptor func(*T) (StringCodec, error)) {
	internal.StrconvxNS.Adapt(strconvx.ToAnyAdaptor(adaptor))
}

// Deprecated: Use RegisterCodec instead.
func RegisterCoder[T any](adapt func(*T) (StringCodec, error)) {
	RegisterCodec(adapt)
}

// RegisterNamedCodec works similar to RegisterCodec, except that it binds the
// codec to a name. This is useful when you only want to override the types in
// a specific struct field. You will be using the "codec" directive
// to specify the name of the codec to use. For example:
//
//	type MyStruct struct {
//		Bool bool // use default bool codec
//		YesNo bool `in:"codec=yesno"` // use YesNo codec
//	}
//
//	func init() {
//		core.RegisterNamedCodec[bool]("yesno", func(b *bool) (core.StringCodec, error) {
//			return (*YesNo)(b), nil
//		})
//	}
//
//	type YesNo bool
//
//	func (yn YesNo) String() string {
//		if yn {
//			return "yes"
//		}
//		return "no"
//	}
//
//	func (yn *YesNo) FromString(s string) error {
//		switch s {
//		case "yes":
//			*yn = true
//		case "no":
//			*yn = false
//		default:
//			return fmt.Errorf("invalid YesNo value: %q", s)
//		}
//		return nil
//	}
func RegisterNamedCodec[T any](name string, adapt func(*T) (StringCodec, error)) {
	typ, adaptor := strconvx.ToAnyAdaptor(adapt)
	namedStringCodecAdaptors[name] = &NamedStringCodecAdaptor{
		Name:     name,
		BaseType: typ,
		Adaptor:  adaptor,
	}
}

// Deprecated: Use RegisterNamedCodec instead.
func RegisterNamedCoder[T any](name string, adapt func(*T) (StringCodec, error)) {
	RegisterNamedCodec(name, adapt)
}

// RegisterFileCodec registers the given type T as a file type. T must implement
// the FileCodec interface. Remember if you don't register the type explicitly,
// it won't be recognized as a file type.
func RegisterFileCodec[T FileCodec]() {
	fileTypes[internal.TypeOf[T]()] = struct{}{}
}

// Deprecated: Use RegisterFileCodec instead.
func RegisterFileCoder[T FileCodec]() {
	RegisterFileCodec[T]()
}

type NamedStringCodecAdaptor struct {
	Name     string
	BaseType reflect.Type
	Adaptor  StringCodecAdaptor
}

func isFileType(typ reflect.Type) bool {
	baseType, _ := internal.BaseTypeOf(typ)
	_, ok := fileTypes[baseType]
	return ok
}
