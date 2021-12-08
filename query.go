package httpin

import "mime/multipart"

// queryValueExtractor implements the "query" executor who extracts values from
// the querystring of an HTTP request.
func queryValueExtractor(ctx *DirectiveContext) error {
	extractor := &extractor{
		Form: multipart.Form{
			Value: ctx.Request.URL.Query(),
		},
	}
	return extractor.Execute(ctx)
}
