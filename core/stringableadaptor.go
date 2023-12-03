package core

import (
	"reflect"

	"github.com/ggicci/httpin/internal"
)

type AnyStringableAdaptor = internal.AnyStringableAdaptor

var (
	customStringableAdaptors = make(map[reflect.Type]AnyStringableAdaptor)
	namedStringableAdaptors  = make(map[string]*NamedAnyStringableAdaptor)
)

// TODO(ggicci): designed to replace RegisterEncoder and RegisterDecoder.
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

type NamedAnyStringableAdaptor struct {
	Name     string
	BaseType reflect.Type
	Adapt    AnyStringableAdaptor
}
