package httpin

import (
	"net/http"
)

// Middleware is a constructor for making a chain, which acts as a list of
// http.Handler constructors. We recommend using
// https://github.com/justinas/alice to chain your HTTP middleware functions
// and the app handler.
type Middleware func(http.Handler) http.Handler
