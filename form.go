package httpin

// formValueExtractor implements the "form" executor who extracts values from
// the forms of an HTTP request.
func formValueExtractor(ctx *DirectiveContext) error {
	return extractFromKVS(ctx, ctx.Request.Form, false)
}
