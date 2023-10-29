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
func UseGochiURLParam(directive string, fn GochiURLParamFunc) {
	RegisterDirective(directive, &directivePath{
		overrideDecode: (&gochiURLParamExtractor{URLParam: fn}).Execute,
	})
}

type gochiURLParamExtractor struct {
	URLParam GochiURLParamFunc
}

func (chi *gochiURLParamExtractor) Execute(rtm *DirectiveRuntime) error {
	req := rtm.GetRequest()
	kvs := make(map[string][]string)

	for _, key := range rtm.Directive.Argv {
		value := chi.URLParam(req, key)
		if value != "" {
			kvs[key] = []string{value}
		}
	}

	extractor := &extractor{
		Runtime: rtm,
		Form: multipart.Form{
			Value: kvs,
		},
	}
	return extractor.Extract()
}
