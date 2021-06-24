// Mux vars extension for github.com/gorilla/mux package.

package httpin

import "net/http"

type MuxVarsFunc func(*http.Request) map[string]string

// UseGorillaMux registers a new directive executor which can extract path
// variables from the URL.
//
// Example: UseGorillaMux("path", mux.Vars)
//
//    type GetUserInput struct {
//       UserID `httpin:"path=user_id"`
//    }
func UseGorillaMux(executor string, fnVars MuxVarsFunc) {
	RegisterDirectiveExecutor(executor, &gorillaMuxVarsExtractor{Vars: fnVars})
}

type gorillaMuxVarsExtractor struct {
	Vars MuxVarsFunc
}

func (mux *gorillaMuxVarsExtractor) Execute(ctx *DirectiveContext) error {
	var kvs = make(map[string][]string)
	for key, value := range mux.Vars(ctx.Request) {
		kvs[key] = []string{value}
	}

	return ExtractFromKVS(ctx, kvs, false)
}
