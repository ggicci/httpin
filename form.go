package httpin

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"time"
)

var basicKinds = map[reflect.Kind]struct{}{
	reflect.Bool:       {},
	reflect.Int:        {},
	reflect.Int8:       {},
	reflect.Int16:      {},
	reflect.Int32:      {},
	reflect.Int64:      {},
	reflect.Uint:       {},
	reflect.Uint8:      {},
	reflect.Uint16:     {},
	reflect.Uint32:     {},
	reflect.Uint64:     {},
	reflect.Float32:    {},
	reflect.Float64:    {},
	reflect.Complex64:  {},
	reflect.Complex128: {},
	reflect.String:     {},
}

var timeType = reflect.TypeOf(time.Time{})

func isBasicType(typ reflect.Type) bool {
	_, ok := basicKinds[typ.Kind()]
	return ok
}

func isTimeType(typ reflect.Type) bool {
	return typ == timeType
}

// readForm
func readForm(inputType reflect.Type, form url.Values) (reflect.Value, error) {
	rv := reflect.New(inputType)

	for i := 0; i < inputType.NumField(); i++ {
		field := inputType.Field(i)
		if name := field.Tag.Get("query"); name != "" {
			formValue, _ := form[name]
			// fmt.Printf("query: %v, formValue: %v\n", name, formValue)
			if err := setField(rv.Elem().Field(i), formValue); err != nil {
				return rv, fmt.Errorf("parse field %s: %w", name, err)
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
		if err := setBasicValue(fv, ft, formValue[0]); err != nil {
			return err
		}
		return nil
	}

	if isTimeType(ft) {
		if err := setTimeValue(fv, ft, formValue[0]); err != nil {
			return err
		}
		return nil
	}

	// TODO(ggicci): hook custom parsers

	switch ft.Kind() {
	case reflect.Slice:
		if err := setSliceValue(fv, ft, formValue); err != nil {
			return err
		}
		return nil
	}

	return fmt.Errorf("%s: unsupported type", ft.Name())
}

func setBasicValue(fv reflect.Value, ft reflect.Type, strValue string) error {
	switch ft.Kind() {
	case reflect.Bool:
		if v, err := strconv.ParseBool(strValue); err != nil {
			return err
		} else {
			fv.SetBool(v)
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if v, err := strconv.ParseInt(strValue, 10, 64); err != nil {
			return err
		} else {
			fv.SetInt(v)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if v, err := strconv.ParseUint(strValue, 10, 64); err != nil {
			return err
		} else {
			fv.SetUint(v)
		}
	case reflect.Float32, reflect.Float64:
		if v, err := strconv.ParseFloat(strValue, 10); err != nil {
			return err
		} else {
			fv.SetFloat(v)
		}
	case reflect.Complex64, reflect.Complex128:
		if v, err := strconv.ParseComplex(strValue, 128); err != nil {
			return err
		} else {
			fv.SetComplex(v)
		}
	case reflect.String:
		fv.SetString(strValue)
	}
	return nil
}

func tryParsingTime(value string) (time.Time, error) {
	// Try parsing value as RFC3339 format.
	if t, err := time.Parse(time.RFC3339Nano, value); err == nil {
		return t, nil
	}

	// Try parsing value as int64 (timestamp).
	// TODO(ggicci): can support float timestamp, e.g. 1618974933.284368
	if timestamp, err := strconv.ParseInt(value, 10, 64); err == nil {
		return time.Unix(timestamp, 0), nil
	}

	return time.Time{}, fmt.Errorf("invalid time value, use time.RFC3339Nano format or timestamp")
}

func setTimeValue(fv reflect.Value, ft reflect.Type, strValue string) error {
	// Try parsing strValue as time.Time in following formats.
	timeValue, err := tryParsingTime(strValue)
	if err != nil {
		return err
	}
	fv.Set(reflect.ValueOf(timeValue))
	return nil
}

func setSliceValue(fv reflect.Value, ft reflect.Type, formValue []string) error {
	elemType := ft.Elem()

	if isBasicType(elemType) {
		rSlice := reflect.MakeSlice(ft, len(formValue), len(formValue))
		for i, strValue := range formValue {
			if err := setBasicValue(rSlice.Index(i), elemType, strValue); err != nil {
				return fmt.Errorf("at index %d: %w", i, err)
			}
		}
		fv.Set(rSlice)
		return nil
	}

	if isTimeType(elemType) {
		rSlice := reflect.MakeSlice(ft, len(formValue), len(formValue))
		for i, strValue := range formValue {
			if err := setTimeValue(rSlice.Index(i), elemType, strValue); err != nil {
				return fmt.Errorf("at index %d: %w", i, err)
			}
		}
		fv.Set(rSlice)
		return nil
	}

	// TODO(ggicci): hook custom parsers
	return fmt.Errorf("%s: unsupported element type of slice", ft.Name())
}
