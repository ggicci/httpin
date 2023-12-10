package core

import (
	"fmt"
	"reflect"

	"github.com/ggicci/httpin/internal"
)

type AnyStringableAdaptor = internal.AnyStringableAdaptor

var (
	fileTypes                = make(map[reflect.Type]struct{})
	customStringableAdaptors = make(map[reflect.Type]AnyStringableAdaptor)
	namedStringableAdaptors  = make(map[string]*NamedAnyStringableAdaptor)
)

func RegisterType[T any](adapt func(*T) (Stringable, error)) {
	customStringableAdaptors[internal.TypeOf[T]()] = internal.ToAnyStringableAdaptor[T](adapt)
}

func RegisterNamedType[T any](name string, adapt func(*T) (Stringable, error)) {
	namedStringableAdaptors[name] = &NamedAnyStringableAdaptor{
		Name:     name,
		BaseType: internal.TypeOf[T](),
		Adapt:    internal.ToAnyStringableAdaptor[T](adapt),
	}
}

func RegisterFileType[T Fileable]() error {
	typ := internal.TypeOf[T]()
	if !typ.Implements(fileableType) {
		return fmt.Errorf("file type must implement Fileable interface")
	}
	fileTypes[typ] = struct{}{}
	return nil
}

type NamedAnyStringableAdaptor struct {
	Name     string
	BaseType reflect.Type
	Adapt    AnyStringableAdaptor
}

func isFileType(typ reflect.Type) bool {
	baseType, _ := BaseTypeOf(typ)
	_, ok := fileTypes[baseType]
	return ok
}