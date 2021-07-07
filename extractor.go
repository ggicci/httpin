package httpin

import (
	"fmt"
	"net/http"
	"reflect"
)

func ExtractFromKVS(ctx *DirectiveContext, kvs map[string][]string, isHeaderKey bool) error {
	for _, key := range ctx.directive.Argv {
		debug("    > execute directive %q with key %q\n", ctx.directive.Executor, key)
		if isHeaderKey {
			key = http.CanonicalHeaderKey(key)
		}
		if err := extractFromKVSWithKey(ctx, kvs, key); err != nil {
			return err
		}
	}
	return nil
}

func extractFromKVSWithKey(ctx *DirectiveContext, kvs map[string][]string, key string) error {
	if ctx.Context.Value(FieldSet) == true {
		debug("    > field already set, skip\n")
		return nil
	}

	// NOTE(ggicci): Array?
	if ctx.ValueType.Kind() == reflect.Slice {
		return extractFromKVSWithKeyForSlice(ctx, kvs, key)
	}

	decoder := decoderOf(ctx.ValueType)
	if decoder == nil {
		return UnsupportedTypeError{ctx.ValueType}
	}

	formValues, exists := kvs[key]
	if !exists {
		debug("    > key %q not found in %s\n", key, ctx.Executor)
		return nil
	}
	var got string
	if len(formValues) > 0 {
		got = formValues[0]
	}
	if interfaceValue, err := decoder.Decode([]byte(got)); err != nil {
		return fieldError{key, got, err}
	} else {
		ctx.Value.Elem().Set(reflect.ValueOf(interfaceValue))
	}

	ctx.DeliverContextValue(FieldSet, true)
	return nil
}

func extractFromKVSWithKeyForSlice(ctx *DirectiveContext, kvs map[string][]string, key string) error {
	elemType := ctx.ValueType.Elem()

	decoder := decoderOf(elemType)
	if decoder == nil {
		return UnsupportedTypeError{ctx.ValueType}
	}

	formValues, exists := kvs[key]
	if !exists {
		debug("    > key %q not found in %s\n", key, ctx.Executor)
		return nil
	}

	theSlice := reflect.MakeSlice(ctx.ValueType, len(formValues), len(formValues))
	for i, formValue := range formValues {
		if interfaceValue, err := decoder.Decode([]byte(formValue)); err != nil {
			return fieldError{key, formValues, fmt.Errorf("at index %d: %w", i, err)}
		} else {
			theSlice.Index(i).Set(reflect.ValueOf(interfaceValue))
		}
	}

	ctx.Value.Elem().Set(theSlice)
	ctx.DeliverContextValue(FieldSet, true)
	return nil
}
