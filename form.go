// directive: "form"
// https://ggicci.github.io/httpin/directives/form

package httpin

// formValueExtractor implements the "form" executor who extracts values from
// the forms of an HTTP request.
func formValueExtractor(ctx *DirectiveContext) error {
	return newExtractor(ctx.Request).Execute(ctx)
}
