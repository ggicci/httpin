// Mux vars extension for github.com/gorilla/mux package.

package integration

import (
	"mime/multipart"
	"net/http"

	"github.com/ggicci/httpin"
	"github.com/ggicci/httpin/directive"
	"github.com/ggicci/httpin/internal"
)

// GorillaMuxVarsFunc is mux.Vars
type GorillaMuxVarsFunc func(*http.Request) map[string]string

// UseGorillaMux registers a new directive executor which can extract values
// from `mux.Vars`, i.e. path variables.
// https://ggicci.github.io/httpin/integrations/gorilla
//
// Usage:
//
//	func init() {
//	    httpin.UseGorillaMux("path", mux.Vars)
//	}
func UseGorillaMux(name string, fnVars GorillaMuxVarsFunc) {
	httpin.RegisterDirective(
		name,
		directive.NewDirectivePath((&gorillaMuxVarsExtractor{Vars: fnVars}).Execute),
	)
}

type gorillaMuxVarsExtractor struct {
	Vars GorillaMuxVarsFunc
}

func (mux *gorillaMuxVarsExtractor) Execute(rtm *httpin.DirectiveRuntime) error {
	req := rtm.GetRequest()
	kvs := make(map[string][]string)

	for key, value := range mux.Vars(req) {
		kvs[key] = []string{value}
	}

	Extractor := &internal.Extractor{
		Runtime: rtm,
		Form: multipart.Form{
			Value: kvs,
		},
	}
	return Extractor.Extract()
}
