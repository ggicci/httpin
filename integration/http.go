package integration

import (
	"github.com/ggicci/httpin/core"
	"mime/multipart"
	"net/http"
)

type HttpMuxVarsFunc func(*http.Request) map[string]string

func UseHttpMux(name string) {
	core.RegisterDirective(
		name,
		core.NewDirectivePath((&httpMuxVarsExtractor{}).Execute),
		true,
	)
}

func UseHttpPathMux() {
	UseHttpMux("path")
}

type httpMuxVarsExtractor struct {
	Vars http.ServeMux
}

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
