package httpin

// RequiredField implements the "required" executor who indicates that the field
// must be set. If the field value were not set by former executors, error
// `ErrMissingField` will be returned.
func RequireField(ctx *DirectiveContext) error {
	if ctx.Context.Value(fieldSet) == nil {
		return ErrMissingField
	}
	return nil
}
