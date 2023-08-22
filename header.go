// directive: "header"
// https://ggicci.github.io/httpin/directives/header

package httpin

import (
	"mime/multipart"
	"net/http"
)

// headerValueExtractor implements the "header" executor who extracts values
// from the HTTP headers.
func headerValueExtractor(ctx *DirectiveRuntime) error {
	req := ctx.Context.Value(RequestValue).(*http.Request)
	extractor := &extractor{
		Form: multipart.Form{
			Value: req.Header,
		},
		KeyNormalizer: http.CanonicalHeaderKey,
	}
	return extractor.Execute(ctx)
}
