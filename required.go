// directive: "required"
// https://ggicci.github.io/httpin/directives/required

package httpin

// required implements the "required" executor who indicates that the field
// must be set. If the field value were not set by former executors, error
// `ErrMissingField` will be returned.
func required(ctx *DirectiveContext) error {
	if ctx.Context.Value(FieldSet) == nil {
		return ErrMissingField
	}
	return nil
}
