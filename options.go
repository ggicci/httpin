package httpin

type Option func(*Engine) error

// WithErrorHandler overrides the default error handler.
func WithErrorHandler(custom ErrorHandler) Option {
	return func(c *Engine) error {
		if custom == nil {
			return ErrNilErrorHandler
		}
		c.errorHandler = custom
		return nil
	}
}

// WithMaxMemory overrides the default maximum memory size (32MB) when reading
// the request body. See https://pkg.go.dev/net/http#Request.ParseMultipartForm
// for more details.
func WithMaxMemory(maxMemory int64) Option {
	return func(c *Engine) error {
		if maxMemory < minimumMaxMemory {
			return ErrMaxMemoryTooSmall
		}
		c.maxMemory = maxMemory
		return nil
	}
}
