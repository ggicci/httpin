package httpin

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
)

type Resolver interface {
	Resolve([]byte) interface{}
}

type ResolverContext struct {
	Source   string // e.g. query, header, body
	Key      string // e.g. page, x-api-token
	resolver Resolver
}

type TypeResolver struct {
	Type      reflect.Type
	Field     reflect.StructField
	Resolvers []ResolverContext
	Path      []string
	Fields    []TypeResolver
}

func (r *TypeResolver) IsRoot() bool {
	return r.Field.Name == ""
}

func BuildTypeResolver(t reflect.Type) (TypeResolver, error) {
	return buildTypeResolver(t)
}

func (r *TypeResolver) dump(buffer *bytes.Buffer, indent int) {
	// indent
	buffer.WriteString(fmt.Sprintf("%s-", strings.Repeat(" ", indent)))

	// field name
	if !r.IsRoot() {
		buffer.WriteString(" [" + r.Field.Name + "]")
	}

	// type
	buffer.WriteString(" " + r.Type.String())
	buffer.WriteString(fmt.Sprintf(" (%d)", len(r.Fields)))
	buffer.WriteString("\n")

	for _, field := range r.Fields {
		field.dump(buffer, indent+4)
	}
}

func (t *TypeResolver) DumpTree() string {
	var buffer bytes.Buffer
	buffer.WriteString("\n")
	t.dump(&buffer, 0)
	return string(buffer.Bytes())
}

func buildTypeResolver(t reflect.Type) (TypeResolver, error) {
	root := TypeResolver{
		Type:      t,
		Field:     reflect.StructField{},
		Resolvers: nil,
		Path:      []string{},
		Fields:    make([]TypeResolver, 0),
	}

	if isBasicType(t) {
		return root, nil
	}

	if isTimeType(t) {
		return root, nil
	}

	if t.Kind() == reflect.Struct {
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			child, err := buildTypeResolver(field.Type)
			if err != nil {
				return root, err
			}
			child.Field = field
			root.Fields = append(root.Fields, child)
		}

		return root, nil
	}

	return root, UnsupportedTypeError{Type: t}
}
