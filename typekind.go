package httpin

import "reflect"

type typeKind int

const (
	typeKindScalar     typeKind = iota // T
	typeKindMulti                      // []T
	typeKindPatch                      // patch.Field[T]
	typeKindPatchMulti                 // patch.Field[[]T]
)

// baseTypeOf returns the scalar element type of a given type.
//   - T -> T, typeKindScalar
//   - []T -> T, typeKindMulti
//   - patch.Field[T] -> T, typeKindPatch
//   - patch.Field[[]T] -> T, typeKindPatchMulti
//
// The given type is gonna use the decoder of the scalar element type to decode
// the input values.
func baseTypeOf(valueType reflect.Type) (reflect.Type, typeKind) {
	if valueType.Kind() == reflect.Slice {
		return valueType.Elem(), typeKindMulti
	}
	if isPatchField(valueType) {
		subElemType, isMulti := patchFieldElemType(valueType)
		if isMulti {
			return subElemType, typeKindPatchMulti
		} else {
			return subElemType, typeKindPatch
		}
	}
	return valueType, typeKindScalar
}

// typeOf returns the reflect.Type of a given type.
// e.g. typeOf[int]() returns reflect.TypeOf(0)
func typeOf[T any]() reflect.Type {
	var zero [0]T
	return reflect.TypeOf(zero).Elem()
}
