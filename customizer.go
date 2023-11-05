package httpin

import "github.com/ggicci/httpin/internal"

func Customizer() customizer {
	return customizerImpl{}
}

type customizer interface {
	// RegisterDirective registers a DirectiveExecutor with the given directive name. The
	// directive should be able to both extract the value from the HTTP request and build
	// the HTTP request from the value. The Decode API is used to decode data from the HTTP
	// request to a field of the input struct, and Encode API is used to encode the field of
	// the input struct to the HTTP request.
	//
	// Will panic if the name were taken or given executor is nil. Pass parameter force
	// (true) to ignore the name conflict.
	RegisterDirective(name string, executor DirectiveExecutor, force ...bool)

	// RegisterErrorHandler replaces the default error handler with the given
	// custom error handler. The default error handler will be used in the http.Handler
	// that decoreated by the middleware created by NewInput().
	RegisterErrorHandler(handler errorHandler)
}

type customizerImpl struct{}

func (customizerImpl) RegisterDirective(name string, executor DirectiveExecutor, force ...bool) {
	registerDirectiveExecutorToNamespace(decoderNamespace, name, executor, force...)
	registerDirectiveExecutorToNamespace(encoderNamespace, name, executor, force...)
}

func (customizerImpl) RegisterErrorHandler(handler errorHandler) {
	internal.PanicOnError(validateErrorHandler(handler))
	globalCustomErrorHandler = handler
}
