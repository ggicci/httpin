package httpin

import (
	"reflect"
	"strings"
)

func isPatchField(t reflect.Type) bool {
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
