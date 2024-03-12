package integration

import (
	"mime/multipart"

	"github.com/ggicci/httpin/core"
	"github.com/labstack/echo/v4"
)

// UseEchoRouter registers a new directive executor which can extract values
// from path variables.
// https://ggicci.github.io/httpin/integrations/echo
//
// Usage:
//
//	import httpin_integration "github.com/ggicci/httpin/integration"
//
//	func init() {
//	    e := echo.New()
//	    httpin_integration.UseEchoRouter("path", e)
//
// // or
//
//	    httpin_integration.UseEchoPathRouter(e)
//	}
func UseEchoRouter(name string, e *echo.Echo) {
	core.RegisterDirective(
		name,
		core.NewDirectivePath((&echoRouterExtractor{e}).Execute),
		true,
	)
}

func UseEchoPathRouter(e *echo.Echo) {
	UseEchoRouter("path", e)
}

// echoRouterExtractor is an extractor for mux.Vars
type echoRouterExtractor struct {
	e *echo.Echo
}

func (mux *echoRouterExtractor) Execute(rtm *core.DirectiveRuntime) error {
	req := rtm.GetRequest()
	kvs := make(map[string][]string)

	c := mux.e.NewContext(req, nil)
	c.SetRequest(req)

	mux.e.Router().Find(req.Method, req.URL.Path, c)

	for _, key := range c.ParamNames() {
		kvs[key] = []string{c.Param(key)}
	}

	extractor := &core.FormExtractor{
		Runtime: rtm,
		Form: multipart.Form{
			Value: kvs,
		},
	}
	return extractor.Extract()
}
