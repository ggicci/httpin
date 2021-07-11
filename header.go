package httpin

// headerValueExtractor implements the "header" executor who extracts values
// from the HTTP headers.
func headerValueExtractor(ctx *DirectiveContext) error {
	return extractFromKVS(ctx, ctx.Request.Header, true)
}
