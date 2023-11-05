// integration: "gochi"
// https://ggicci.github.io/httpin/integrations/gochi

package integration

import (
	"mime/multipart"
	"net/http"

	"github.com/ggicci/httpin"
	"github.com/ggicci/httpin/directive"
	"github.com/ggicci/httpin/internal"
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
	httpin.Customizer().RegisterDirective(
		name,
		directive.NewDirectivePath((&gochiURLParamExtractor{URLParam: fn}).Execute),
	)
}

type gochiURLParamExtractor struct {
	URLParam GochiURLParamFunc
}

func (chi *gochiURLParamExtractor) Execute(rtm *httpin.DirectiveRuntime) error {
	req := rtm.GetRequest()
	kvs := make(map[string][]string)

	for _, key := range rtm.Directive.Argv {
		value := chi.URLParam(req, key)
		if value != "" {
			kvs[key] = []string{value}
		}
	}

	Extractor := &internal.Extractor{
		Runtime: rtm,
		Form: multipart.Form{
			Value: kvs,
		},
	}
	return Extractor.Extract()
}
