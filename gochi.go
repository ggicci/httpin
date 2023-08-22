// integration: "gochi"
// https://ggicci.github.io/httpin/integrations/gochi

package httpin

import (
	"mime/multipart"
	"net/http"
)

// GochiURLParamFunc is chi.URLParam
type GochiURLParamFunc func(r *http.Request, key string) string

// UseGochiURLParam registers a directive executor which can extract values
// from `chi.URLParam`, i.e. path variables.
// https://ggicci.github.io/httpin/integrations/gochi
//
// Usage:
//
//	func init() {
//	    httpin.UseGochiURLParam("path", chi.URLParam)
//	}
func UseGochiURLParam(executor string, fn GochiURLParamFunc) {
	RegisterDirectiveExecutor(executor, &gochiURLParamExtractor{URLParam: fn})
}

type gochiURLParamExtractor struct {
	URLParam GochiURLParamFunc
}

func (chi *gochiURLParamExtractor) Execute(ctx *DirectiveRuntime) error {
	req := ctx.Context.Value(RequestValue).(*http.Request)
	kvs := make(map[string][]string)

	for _, key := range ctx.Directive.Argv {
		value := chi.URLParam(req, key)
		if value != "" {
			kvs[key] = []string{value}
		}
	}

	extractor := &extractor{
		Form: multipart.Form{
			Value: kvs,
		},
	}
	return extractor.Execute(ctx)
}
