package httpin

// queryValueExtractor implements the "query" executor who extracts values from
// the querystring of an HTTP request.
func queryValueExtractor(ctx *DirectiveContext) error {
	return extractFromKVS(ctx, ctx.Request.URL.Query(), false)
}
