// integration: "gochi"
// https://ggicci.github.io/httpin/integrations/gochi

package integration

import (
	"mime/multipart"
	"net/http"

	"github.com/ggicci/httpin/core"
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
func UseGochiURLParam(name string, fn GochiURLParamFunc) {
	core.RegisterDirective(
		name,
		core.NewDirectivePath((&gochiURLParamExtractor{URLParam: fn}).Execute),
	)
}

type gochiURLParamExtractor struct {
	URLParam GochiURLParamFunc
}

func (chi *gochiURLParamExtractor) Execute(rtm *core.DirectiveRuntime) error {
	req := rtm.GetRequest()
	kvs := make(map[string][]string)

	for _, key := range rtm.Directive.Argv {
		value := chi.URLParam(req, key)
		if value != "" {
			kvs[key] = []string{value}
		}
	}

	extractor := &core.FormExtractor{
		Runtime: rtm,
		Form: multipart.Form{
			Value: kvs,
		},
	}
	return extractor.Extract()
}
