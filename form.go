// directive: "form"
// https://ggicci.github.io/httpin/directives/form

package httpin

import "net/http"

// formValueExtractor implements the "form" executor who extracts values from
// the forms of an HTTP request.
func formValueExtractor(ctx *DirectiveRuntime) error {
	req := ctx.Context.Value(RequestValue).(*http.Request)
	return newExtractor(req).Execute(ctx)
}
