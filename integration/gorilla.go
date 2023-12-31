// Mux vars extension for github.com/gorilla/mux package.

package integration

import (
	"mime/multipart"
	"net/http"

	"github.com/ggicci/httpin/core"
)

// GorillaMuxVarsFunc is mux.Vars
type GorillaMuxVarsFunc func(*http.Request) map[string]string

// UseGorillaMux registers a new directive executor which can extract values
// from `mux.Vars`, i.e. path variables.
// https://ggicci.github.io/httpin/integrations/gorilla
//
// Usage:
//
//	import httpin_integration "github.com/ggicci/httpin/integration"
//
//	func init() {
//	    httpin_integration.UseGorillaMux("path", mux.Vars)
//	}
func UseGorillaMux(name string, fnVars GorillaMuxVarsFunc) {
	core.RegisterDirective(
		name,
		core.NewDirectivePath((&gorillaMuxVarsExtractor{Vars: fnVars}).Execute),
		true,
	)
}

type gorillaMuxVarsExtractor struct {
	Vars GorillaMuxVarsFunc
}

func (mux *gorillaMuxVarsExtractor) Execute(rtm *core.DirectiveRuntime) error {
	req := rtm.GetRequest()
	kvs := make(map[string][]string)

	for key, value := range mux.Vars(req) {
		kvs[key] = []string{value}
	}

	extractor := &core.FormExtractor{
		Runtime: rtm,
		Form: multipart.Form{
			Value: kvs,
		},
	}
	return extractor.Extract()
}
