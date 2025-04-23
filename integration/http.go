package integration

import (
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/ggicci/httpin/core"
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

	for _, value := range strings.Split(req.Pattern, "/") {
		if strings.HasPrefix(value, "{") && strings.HasSuffix(value, "}") {
			value = strings.TrimSuffix(strings.TrimPrefix(value, "{"), "}")
			kvs[value] = []string{req.PathValue(value)}
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
