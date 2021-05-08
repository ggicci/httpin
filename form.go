package httpin

import (
	"fmt"
	"net/http"
	"reflect"
)

// FormValueExtractor implements the "form" executor who extracts values from
// the forms of an HTTP request.
func FormValueExtractor(ctx *DirectiveContext) error {
	return extractFromKVS(ctx, ctx.Request.Form, false)
}

// HeaderValueExtractor implements the "header" executor who extracts values
// from the HTTP headers.
func HeaderValueExtractor(ctx *DirectiveContext) error {
	return extractFromKVS(ctx, ctx.Request.Header, true)
}

func extractFromKVS(ctx *DirectiveContext, kvs map[string][]string, headerKey bool) error {
	for _, key := range ctx.Directive.Argv {
		debug("    > execute directive %q with key %q\n", ctx.Directive.Executor, key)
		if headerKey {
			key = http.CanonicalHeaderKey(key)
		}
		if err := extractFromKVSWithKey(ctx, kvs, key); err != nil {
			return err
		}
	}
	return nil
}

func extractFromKVSWithKey(ctx *DirectiveContext, kvs map[string][]string, key string) error {
	if ctx.Context.Value(fieldSet) == true {
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
		return err
	} else {
		ctx.Value.Elem().Set(reflect.ValueOf(interfaceValue))
	}

	ctx.DeliverContextValue(fieldSet, true)
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
			return fmt.Errorf("at index %d: %w", i, err)
		} else {
			theSlice.Index(i).Set(reflect.ValueOf(interfaceValue))
		}
	}

	ctx.Value.Elem().Set(theSlice)
	ctx.DeliverContextValue(fieldSet, true)
	return nil
}
