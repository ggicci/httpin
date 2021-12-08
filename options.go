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

func WithMaxMemory(maxMemory int64) Option {
	return func(c *Engine) error {
		if maxMemory < minimumMaxMemory {
			return ErrMaxMemoryTooSmall
		}
		c.maxMemory = maxMemory
		return nil
	}
}
