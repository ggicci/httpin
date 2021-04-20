package httpin

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"
)

var basicKinds = map[reflect.Kind]struct{}{
	reflect.Bool:       struct{}{},
	reflect.Int:        struct{}{},
	reflect.Int8:       struct{}{},
	reflect.Int16:      struct{}{},
	reflect.Int32:      struct{}{},
	reflect.Int64:      struct{}{},
	reflect.Uint:       struct{}{},
	reflect.Uint8:      struct{}{},
	reflect.Uint16:     struct{}{},
	reflect.Uint32:     struct{}{},
	reflect.Uint64:     struct{}{},
	reflect.Float32:    struct{}{},
	reflect.Float64:    struct{}{},
	reflect.Complex64:  struct{}{},
	reflect.Complex128: struct{}{},
	reflect.String:     struct{}{},
}

func isBasicType(typ reflect.Type) bool {
	_, ok := basicKinds[typ.Kind()]
	return ok
}

// readForm
func readForm(inputType reflect.Type, form url.Values) (reflect.Value, error) {
	rv := reflect.New(inputType)

	for i := 0; i < inputType.NumField(); i++ {
		field := inputType.Field(i)
		if name := field.Tag.Get("query"); name != "" {
			formValue, _ := form[name]
			fmt.Printf("query: %v, formValue: %v\n", name, formValue)
			if err := setField(rv.Elem().Field(i), formValue); err != nil {
				return rv, err
			}
		}
	}

	return rv, nil
}

func setField(fv reflect.Value, formValue []string) error {
	if len(formValue) == 0 {
		// TODO(ggicci): throw an error if decorator like "required" set?
		// I think the validation can be handled by other libraries.
		return nil
	}

	ft := fv.Type()
	if isBasicType(ft) {
		setBasicValue(fv, ft, formValue[0])
		return nil
	}

	switch ft.Kind() {
	case reflect.Slice:
		// Create new slice.
		elemType := ft.Elem()
		if !isBasicType(elemType) {
			return fmt.Errorf("%s: unsupported element type of slice", ft.Name())
		}
		rSlice := reflect.MakeSlice(ft, len(formValue), len(formValue))
		for i, strValue := range formValue {
			setBasicValue(rSlice.Index(i), elemType, strValue)
		}
		return nil
		// TODO(ggicci): support custom types
	}

	return fmt.Errorf("%s: unsupported type", ft.Name())
}

func setBasicValue(fv reflect.Value, ft reflect.Type, strValue string) error {
	switch ft.Kind() {
	case reflect.Bool:
		if v, err := strconv.ParseBool(strValue); err != nil {
			return fmt.Errorf("%s: %v", ft.Name(), err)
		} else {
			fv.SetBool(v)
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if v, err := strconv.ParseInt(strValue, 10, 64); err != nil {
			return fmt.Errorf("%s: %v", ft.Name(), err)
		} else {
			fv.SetInt(v)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if v, err := strconv.ParseUint(strValue, 10, 64); err != nil {
			return fmt.Errorf("%s: %v", ft.Name(), err)
		} else {
			fv.SetUint(v)
		}
	case reflect.Float32, reflect.Float64:
		if v, err := strconv.ParseFloat(strValue, 10); err != nil {
			return fmt.Errorf("%s: %v", ft.Name(), err)
		} else {
			fv.SetFloat(v)
		}
	case reflect.Complex64, reflect.Complex128:
		if v, err := strconv.ParseComplex(strValue, 128); err != nil {
			return fmt.Errorf("%s: %v", ft.Name(), err)
		} else {
			fv.SetComplex(v)
		}
	case reflect.String:
		fv.SetString(strValue)
	}
	return nil
}
