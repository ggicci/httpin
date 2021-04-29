package httpin

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"
)

type Resolver interface {
	Resolve([]byte) interface{}
}

type FieldResolver struct {
	Type       reflect.Type
	Field      reflect.StructField
	Path       []string
	Directives []*Directive
	Fields     []*FieldResolver
}

func (r *FieldResolver) IsRoot() bool {
	return r.Field.Name == ""
}

func (r *FieldResolver) newInstance() reflect.Value {
	return reflect.New(r.Type)
}

func (r *FieldResolver) resolve(req *http.Request) (reflect.Value, error) {
	rv := r.newInstance()

	// Execute directives.
	if len(r.Directives) > 0 {
		directiveContext := &DirectiveContext{
			Request: req,
			Value:   rv,
		}

		for _, dir := range r.Directives {
			if err := dir.Execute(directiveContext); err != nil {
				return rv, fmt.Errorf("execute directive %q failed: %w", dir.Executor, err)
			}
		}
	}

	if len(r.Fields) > 0 {
		for i, fr := range r.Fields {
			field, err := fr.resolve(req)
			if err != nil {
				return rv, err
			}
			rv.Elem().Field(i).Set(field.Elem())
		}
	}

	return rv, nil
}

// buildResolverTree builds a resolver tree for the specified struct type.
// Which helps resolving fields data from input sources.
func buildResolverTree(t reflect.Type) (*FieldResolver, error) {
	root := &FieldResolver{Type: t}
	for i := 0; i < t.NumField(); i++ {
		fieldResolver, err := buildFieldResolver(root, t.Field(i))
		if err != nil {
			return root, err
		}
		root.Fields = append(root.Fields, fieldResolver)
	}

	return root, nil
}

func buildFieldResolver(parent *FieldResolver, field reflect.StructField) (*FieldResolver, error) {
	t := field.Type
	root := &FieldResolver{
		Type:  t,
		Field: field,
		Path:  make([]string, len(parent.Path)+1),
	}
	copy(root.Path, parent.Path)
	root.Path[len(root.Path)-1] = field.Name
	directives, err := parseStructTag(field)
	if err != nil {
		return nil, fmt.Errorf("parse struct tag error: %w", err)
	}
	root.Directives = directives

	if isBasicType(t) {
		return root, nil
	}

	if isTimeType(t) {
		return root, nil
	}

	if isArrayType(t) {
		return root, nil
	}

	if t.Kind() == reflect.Struct {
		for i := 0; i < t.NumField(); i++ {
			fieldResolver, err := buildFieldResolver(root, t.Field(i))
			if err != nil {
				return root, err
			}
			root.Fields = append(root.Fields, fieldResolver)
		}

		return root, nil
	}

	return root, UnsupportedTypeError{Type: t}
}

func parseStructTag(field reflect.StructField) ([]*Directive, error) {
	directives := make([]*Directive, 0)
	// Parse and build resolvers from field struct tag. Tag examples:
	// "query.name"
	// "query.access_token,header.x-api-token"
	inTag := field.Tag.Get("in")
	if inTag == "" {
		return directives, nil // skip
	}
	for _, key := range strings.Split(inTag, ",") {
		directive, err := buildDirective(key)
		if err != nil {
			return nil, err
		}
		directives = append(directives, directive)
	}
	return directives, nil
}
