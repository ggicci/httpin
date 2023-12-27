package core

import (
	"reflect"
	"strings"
)

type TypeKind int

const (
	TypeKindT           TypeKind = iota // T
	TypeKindTSlice                      // []T
	TypeKindPatchT                      // patch.Field[T]
	TypeKindPatchTSlice                 // patch.Field[[]T]
)

// BaseTypeOf returns the base type of the given type its kind. The kind
// represents how the given type is constructed from the base type.
//   - T -> T, TypeKindT
//   - []T -> T, TypeKindTSlice
//   - patch.Field[T] -> T, TypeKindPatchT
//   - patch.Field[[]T] -> T, TypeKindPatchTSlice
func BaseTypeOf(valueType reflect.Type) (reflect.Type, TypeKind) {
	if valueType.Kind() == reflect.Slice {
		return valueType.Elem(), TypeKindTSlice
	}
	if IsPatchField(valueType) {
		subElemType, isMulti := patchFieldElemType(valueType)
		if isMulti {
			return subElemType, TypeKindPatchTSlice
		} else {
			return subElemType, TypeKindPatchT
		}
	}
	return valueType, TypeKindT
}

func IsPatchField(t reflect.Type) bool {
	return t.Kind() == reflect.Struct &&
		t.PkgPath() == "github.com/ggicci/httpin/patch" &&
		strings.HasPrefix(t.Name(), "Field[")
}

func patchFieldElemType(t reflect.Type) (reflect.Type, bool) {
	fv, _ := t.FieldByName("Value")
	if fv.Type.Kind() == reflect.Slice {
		return fv.Type.Elem(), true
	}
	return fv.Type, false
}
