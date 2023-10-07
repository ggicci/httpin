// directive: "required"
// https://ggicci.github.io/httpin/directives/required

package httpin

// required implements the "required" executor who indicates that the field
// must be set. If the field value were not set by former executors, error
// `ErrMissingField` will be returned.
// If the required field is a member of a child struct that is nil then no
// error will be returned.
func required(ctx *DirectiveRuntime) error {
	// check that the containing struct is at root, or is non nil.
	if ctx.Resolver.Parent.IsRoot() || ctx.Resolver.Parent.Context.Value(FieldSet) != nil {
		if ctx.Context.Value(FieldSet) == nil {
			return ErrMissingField
		}
	}

	return nil
}
