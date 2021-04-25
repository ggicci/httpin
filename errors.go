package httpin

import "fmt"

type UnsupportedType string

func (e UnsupportedType) Error() string {
	return "unsupported type: " + string(e)
}

type InvalidField struct {
	Name   string
	TagKey string
	Tag    string
	Value  interface{}
	err    error
}

func (f *InvalidField) Error() string {
	return fmt.Sprintf("invalid field: %s(%s, %s): %s", f.Name, f.TagKey, f.Tag, f.err)
}
