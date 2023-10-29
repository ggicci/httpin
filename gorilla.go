// Mux vars extension for github.com/gorilla/mux package.

package httpin

import (
	"mime/multipart"
	"net/http"
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
	RegisterDirective(name, &directivePath{
		overrideDecode: (&gorillaMuxVarsExtractor{Vars: fnVars}).Execute,
	})
}

type gorillaMuxVarsExtractor struct {
	Vars GorillaMuxVarsFunc
}

func (mux *gorillaMuxVarsExtractor) Execute(rtm *DirectiveRuntime) error {
	req := rtm.GetRequest()
	kvs := make(map[string][]string)

	for key, value := range mux.Vars(req) {
		kvs[key] = []string{value}
	}

	extractor := &extractor{
		Runtime: rtm,
		Form: multipart.Form{
			Value: kvs,
		},
	}
	return extractor.Extract()
}
