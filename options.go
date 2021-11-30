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
