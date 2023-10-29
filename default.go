// directive: "default"
// https://ggicci.github.io/httpin/directives/default

package httpin

import "mime/multipart"

func defaultValueSetter(ctx *DirectiveRuntime) error {
	if ctx.Context.Value(FieldSet) != nil {
		return nil // noop, the field was set by a former executor
	}

	// Transform:
	// 1. ctx.Argv -> input values
	// 2. ["default"] -> keys
	extractor := &extractor{
		Form: multipart.Form{
			Value: map[string][]string{
				"default": ctx.Directive.Argv,
			},
		},
	}
	return extractor.Execute(ctx, "default")
}
