// directive: "omitempty"
// https://ggicci.github.io/httpin/directives/omitempty

package core

// DirectiveOmitEmpty is used with the DirectiveQuery, DirectiveForm, and DirectiveHeader to indicate that the field
// should be omitted when the value is empty.
// It does not have any affect when used by itself
type DirectiveOmitEmpty struct{}

func (*DirectiveOmitEmpty) Decode(_ *DirectiveRuntime) error {
	return nil
}

func (*DirectiveOmitEmpty) Encode(_ *DirectiveRuntime) error {
	return nil
}
