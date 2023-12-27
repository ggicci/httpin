package internal

import "fmt"

type StringableAdaptor[T any] func(*T) (Stringable, error)
type AnyStringableAdaptor func(any) (Stringable, error)

func NewAnyStringableAdaptor[T any](adapt StringableAdaptor[T]) AnyStringableAdaptor {
	return func(v any) (Stringable, error) {
		if cv, ok := v.(*T); ok {
			return adapt(cv)
		} else {
			return nil, fmt.Errorf("%w: cannot convert %T to %s", ErrTypeMismatch, v, TypeOf[*T]())
		}
	}
}
