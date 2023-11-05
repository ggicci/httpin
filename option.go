package httpin

import "errors"

type coreOption func(*core) error

// WithErrorHandler overrides the default error handler.
func WithErrorHandler(custom errorHandler) coreOption {
	return func(c *core) error {
		if err := validateErrorHandler(custom); err != nil {
			return err
		} else {
			c.errorHandler = custom
			return nil
		}
	}
}

// WithMaxMemory overrides the default maximum memory size (32MB) when reading
// the request body. See https://pkg.go.dev/net/http#Request.ParseMultipartForm
// for more details.
func WithMaxMemory(maxMemory int64) coreOption {
	return func(c *core) error {
		if maxMemory < minimumMaxMemory {
			return errors.New("max memory too small")
		}
		c.maxMemory = maxMemory
		return nil
	}
}
