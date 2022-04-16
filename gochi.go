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
//     func init() {
//         httpin.UseGochiURLParam("path", chi.URLParam)
//     }
func UseGochiURLParam(executor string, fn GochiURLParamFunc) {
	RegisterDirectiveExecutor(executor, &gochiURLParamExtractor{URLParam: fn}, nil)
}

type gochiURLParamExtractor struct {
	URLParam GochiURLParamFunc
}

func (chi *gochiURLParamExtractor) Execute(ctx *DirectiveContext) error {
	var kvs = make(map[string][]string)

	for _, key := range ctx.Argv {
		value := chi.URLParam(ctx.Request, key)
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
