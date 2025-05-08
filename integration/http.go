package integration

import (
	"mime/multipart"
	"net/http"

	"github.com/ggicci/httpin/core"
)

type HttpMuxVarsFunc func(*http.Request) map[string]string

// UseHttpPathVariable registers a new directive executor which can extract
// values from URL path variables via `http.Request.PathValue` API.
// https://ggicci.github.io/httpin/integrations/http
//
// Usage:
//
//	import httpin_integration "github.com/ggicci/httpin/integration"
//	func init() {
//		httpin_integration.UseHttpPathVariable("path")
//	}
func UseHttpPathVariable(name string) {
	core.RegisterDirective(
		name,
		core.NewDirectivePath((&httpMuxVarsExtractor{}).Execute),
		true,
	)
}

type httpMuxVarsExtractor struct{}

func (mux *httpMuxVarsExtractor) Execute(rtm *core.DirectiveRuntime) error {
	req := rtm.GetRequest()
	kvs := make(map[string][]string)

	for _, key := range rtm.Directive.Argv {
		value := req.PathValue(key)
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
