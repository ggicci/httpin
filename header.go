package httpin

import (
	"mime/multipart"
	"net/http"
)

// headerValueExtractor implements the "header" executor who extracts values
// from the HTTP headers.
func headerValueExtractor(ctx *DirectiveContext) error {
	extractor := &extractor{
		Form: multipart.Form{
			Value: ctx.Request.Header,
		},
		KeyNormalizer: http.CanonicalHeaderKey,
	}
	return extractor.Execute(ctx)
}
