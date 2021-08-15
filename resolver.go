package httpin

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"
)

type fieldResolver struct {
	Type       reflect.Type
	Field      reflect.StructField
	Path       []string
	Directives []*Directive
	Fields     []*fieldResolver
}

func (r *fieldResolver) isBodyDecoderAnnotation() bool {
	return r.Type == bodyTypeAnnotationJSON || r.Type == bodyTypeAnnotationXML
}

func (r *fieldResolver) resolve(req *http.Request) (reflect.Value, error) {
	rv := reflect.New(r.Type)

	// Then execute directives.
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
			if err := dir.Execute(directiveContext); err != nil {
				var (
					fe       fieldError
					gotValue interface{}
				)

				if errors.As(err, &fe) {
					gotValue = fe.Value
				}

				return rv, &InvalidFieldError{
					Field:         r.Field.Name,
					Source:        dir.Executor,
					Value:         gotValue,
					ErrorMessage:  err.Error(),
					internalError: err,
				}
			}
			inheritableContext = directiveContext.Context
		}

		// When all directives got executed, check context value of "StopRecursion"
		// to determine whether we should resolve the "children fields" further.
		if stopRecusrion, ok := inheritableContext.Value(StopRecursion).(bool); ok && stopRecusrion {
			return rv, nil
		}
	}

	if len(r.Fields) > 0 { // struct
		for i, fr := range r.Fields {
			if fr.Field.PkgPath != "" {
				continue // skip unexported field
			}

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
func buildResolverTree(t reflect.Type) (*fieldResolver, error) {
	root := &fieldResolver{Type: t}

	var typeOfBody reflect.Type
	for i := 0; i < t.NumField(); i++ {
		fieldResolver, err := buildFieldResolver(root, t.Field(i))
		if err != nil {
			return nil, err
		}

		// Check if there's a body decoder annotation field.
		if fieldResolver.isBodyDecoderAnnotation() {
			if typeOfBody != nil {
				return nil, fmt.Errorf("%w: %s", ErrDuplicateAnnotationField, fieldResolver.Field.Name)
			}

			// Inject a "body" directive to the root.
			typeOfBody = fieldResolver.Type
			dir, _ := buildDirective(fmt.Sprintf("body=%s", bodyTypeString(typeOfBody)))
			root.Directives = []*Directive{dir}
		}

		root.Fields = append(root.Fields, fieldResolver)
	}

	return root, nil
}

func buildFieldResolver(parent *fieldResolver, field reflect.StructField) (*fieldResolver, error) {
	t := field.Type
	path := make([]string, len(parent.Path)+1)
	copy(path, parent.Path)
	path[len(path)-1] = field.Name

	root := &fieldResolver{
		Type:       t,
		Field:      field,
		Path:       path,
		Directives: make([]*Directive, 0),
	}

	// Skip parsing struct tags if met body resolver annotation field.
	if root.isBodyDecoderAnnotation() {
		return root, nil
	}

	// Parse the struct tag and build the directives.
	if directives, err := parseStructTag(field); err != nil {
		return nil, fmt.Errorf("parse struct tag failed: %w", err)
	} else {
		root.Directives = directives
	}

	if field.Anonymous && t.Kind() == reflect.Struct && len(root.Directives) == 0 {
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
		directive, err := buildDirective(key)
		if err != nil {
			return nil, err
		}
		directives = append(directives, directive)
	}

	return directives, nil
}
