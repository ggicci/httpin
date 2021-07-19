// URLParam extension for github.com/go-chi/chi.

package httpin

import "net/http"

// GochiURLParamFunc is chi.URLParam
type GochiURLParamFunc func(r *http.Request, key string) string

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

	return extractFromKVS(ctx, kvs, false)
}
