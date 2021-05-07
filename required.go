package httpin

func RequireField(ctx *DirectiveContext) error {
	if ctx.Context.Value(fieldSet) == nil {
		return ErrMissingField
	}
	return nil
}
