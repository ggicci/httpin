package httpin

type Option func(*Engine) error

// WithErrorHandler overrides the default error handler.
func WithErrorHandler(h ErrorHandler) Option {
	return func(c *Engine) error {
		if h == nil {
			return ErrNilErrorHandler
		}
		c.errorHandler = h
		return nil
	}
}
