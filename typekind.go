package httpin

import "reflect"

type typeKind int

const (
	typeT           typeKind = iota // T
	typeTSlice                      // []T
	typePatchT                      // patch.Field[T]
	typePatchTSlice                 // patch.Field[[]T]
)

// baseTypeOf returns the base type of the given type its kind. The kind
// represents how the given type is constructed from the base type.
//   - T -> T, typeT
//   - []T -> T, typeTSlice
//   - patch.Field[T] -> T, typePatchT
//   - patch.Field[[]T] -> T, typePatchTSlice
func baseTypeOf(valueType reflect.Type) (reflect.Type, typeKind) {
	if valueType.Kind() == reflect.Slice {
		return valueType.Elem(), typeTSlice
	}
	if isPatchField(valueType) {
		subElemType, isMulti := patchFieldElemType(valueType)
		if isMulti {
			return subElemType, typePatchTSlice
		} else {
			return subElemType, typePatchT
		}
	}
	return valueType, typeT
}

// typeOf returns the reflect.Type of a given type.
// e.g. typeOf[int]() returns reflect.TypeOf(0)
func typeOf[T any]() reflect.Type {
	var zero [0]T
	return reflect.TypeOf(zero).Elem()
}
