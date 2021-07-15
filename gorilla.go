// Mux vars extension for github.com/gorilla/mux package.

package httpin

import "net/http"

// GorillaMuxVarsFunc is mux.Vars
type GorillaMuxVarsFunc func(*http.Request) map[string]string

// UseGorillaMux registers a new directive executor which can extract path
// variables from the URL. Which works as an accompany to gorilla's mux package.
//
// Example:
//
//    UseGorillaMux("path", mux.Vars)
//
//    type GetUserInput struct {
//       UserID `httpin:"path=user_id"`
//    }
func UseGorillaMux(executor string, fnVars GorillaMuxVarsFunc) {
	RegisterDirectiveExecutor(executor, &gorillaMuxVarsExtractor{Vars: fnVars})
}

type gorillaMuxVarsExtractor struct {
	Vars GorillaMuxVarsFunc
}

func (mux *gorillaMuxVarsExtractor) Execute(ctx *DirectiveContext) error {
	var kvs = make(map[string][]string)
	for key, value := range mux.Vars(ctx.Request) {
		kvs[key] = []string{value}
	}

	return extractFromKVS(ctx, kvs, false)
}
