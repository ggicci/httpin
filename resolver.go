package httpin

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"strings"
)

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

func (r *FieldResolver) resolve(req *http.Request) (reflect.Value, error) {
	rv := reflect.New(r.Type)
	debug("resolve: %s (of %s)\n", r.Field.Name, r.Type)

	// Execute directives.
	if len(r.Directives) > 0 {
		inheritableContext := context.Background()
		for _, dir := range r.Directives {
			directiveContext := &DirectiveContext{
				Directive: *dir,
				Request:   req,
				ValueType: r.Type,
				Value:     rv,
				Context:   inheritableContext,
			}
			debug("  > execute directive: %s with %v\n", dir.Executor, dir.Argv)
			if err := dir.Execute(directiveContext); err != nil {
				return rv, &InvalidField{
					Field:         r.Field.Name,
					Source:        dir.Executor,
					Value:         nil, // FIXME(ggicci): add source data
					InternalError: err,
				}
			}
			inheritableContext = directiveContext.Context
		}
	}

	if len(r.Fields) > 0 { // struct
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

// buildResolverTree builds a tree of resolvers for the specified struct type.
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
	directives, err := parseStructTag(field)
	if err != nil {
		return nil, fmt.Errorf("parse struct tag failed: %w", err)
	}
	t := field.Type
	path := make([]string, len(parent.Path)+1)
	copy(path, parent.Path)
	path[len(path)-1] = field.Name
	root := &FieldResolver{
		Type:       t,
		Field:      field,
		Path:       path,
		Directives: directives,
	}

	if field.Anonymous && t.Kind() == reflect.Struct && len(directives) == 0 {
		for i := 0; i < t.NumField(); i++ {
			fieldResolver, err := buildFieldResolver(root, t.Field(i))
			if err != nil {
				return nil, err
			}
			root.Fields = append(root.Fields, fieldResolver)
		}
	}

	return root, nil
}

// parseStructTag parses and builds a resolver for a field by inspecting the struct tag.
// The tag named "in" will be extracted by httpin.
// Example contents of the `in` tag:
//   - `in:"query=name"`
//   - `in:"query=access_token,token;header=x-api-token"`
// Which should conform to the format:
//
//    <intag> := "<direction_1>[;<direction_2>...[;<direction_N>]]"
//    <direction> := <executor>[=<arg_1>[,<arg_2>...[,<arg_N>]]]
//
// For short, use `;` as directions' delimiter, use `,` as arguments' delimiter.
func parseStructTag(field reflect.StructField) ([]*Directive, error) {
	directives := make([]*Directive, 0)
	inTag := field.Tag.Get("in")
	if inTag == "" {
		return directives, nil // skip
	}
	for _, key := range strings.Split(inTag, ";") {
		directive, err := BuildDirective(key)
		if err != nil {
			return nil, err
		}
		directives = append(directives, directive)
	}

	return directives, nil
}
