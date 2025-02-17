package core

import (
	"reflect"

	"github.com/ggicci/httpin/internal"
	"github.com/ggicci/strconvx"
)

type AnyStringConverterAdaptor = strconvx.AnyStringConverterAdaptor

var (
	fileTypes               = make(map[reflect.Type]struct{})
	namedStringableAdaptors = make(map[string]*NamedAnyStringConverterAdaptor)
)

// RegisterCoder registers a custom coder for the given type T. When a field of
// type T is encountered, this coder will be used to convert the value to a
// Stringable, which will be used to convert the value from/to string.
//
// NOTE: this function is designed to override the default Stringable adaptors
// that are registered by this package. For example, if you want to override the
// defualt behaviour of converting a bool value from/to string, you can do this:
//
//	func init() {
//		core.RegisterCoder[bool](func(b *bool) (core.Stringable, error) {
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
func RegisterCoder[T any](adapt func(*T) (Stringable, error)) {
	typ, adaptor := strconvx.ToAnyStringConverterAdaptor[T](adapt)
	strconvxNS.Adapt(typ, adaptor)
}

// RegisterNamedCoder works similar to RegisterCoder, except that it binds the
// coder to a name. This is useful when you only want to override the types in
// a specific struct field. You will be using the "coder" or "decoder" directive
// to specify the name of the coder to use. For example:
//
//	type MyStruct struct {
//		Bool bool // use default bool coder
//		YesNo bool `in:"coder=yesno"` // use YesNo coder
//	}
//
//	func init() {
//		core.RegisterNamedCoder[bool]("yesno", func(b *bool) (core.Stringable, error) {
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
func RegisterNamedCoder[T any](name string, adapt func(*T) (Stringable, error)) {
	typ, adaptor := strconvx.ToAnyStringConverterAdaptor[T](adapt)
	namedStringableAdaptors[name] = &NamedAnyStringConverterAdaptor{
		Name:     name,
		BaseType: typ,
		Adapt:    adaptor,
	}
}

// RegisterFileCoder registers the given type T as a file type. T must implement
// the Fileable interface. Remember if you don't register the type explicitly,
// it won't be recognized as a file type.
func RegisterFileCoder[T Fileable]() error {
	fileTypes[internal.TypeOf[T]()] = struct{}{}
	return nil
}

type NamedAnyStringConverterAdaptor struct {
	Name     string
	BaseType reflect.Type
	Adapt    AnyStringConverterAdaptor
}

func isFileType(typ reflect.Type) bool {
	baseType, _ := BaseTypeOf(typ)
	_, ok := fileTypes[baseType]
	return ok
}
