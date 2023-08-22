// directive: "query"
// https://ggicci.github.io/httpin/directives/query

package httpin

import (
	"mime/multipart"
	"net/http"
)

// queryValueExtractor implements the "query" executor who extracts values from
// the querystring of an HTTP request.
func queryValueExtractor(ctx *DirectiveRuntime) error {
	req := ctx.Context.Value(RequestValue).(*http.Request)
	extractor := &extractor{
		Form: multipart.Form{
			Value: req.URL.Query(),
		},
	}
	return extractor.Execute(ctx)
}
