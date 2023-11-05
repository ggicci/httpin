package core

import "errors"

const minimumMaxMemory = int64(1 << 10)  // 1KB
const defaultMaxMemory = int64(32 << 20) // 32 MB

type Option func(*Core) error

// WithErrorHandler overrides the default error handler.
func WithErrorHandler(custom errorHandler) Option {
	return func(c *Core) error {
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
func WithMaxMemory(maxMemory int64) Option {
	return func(c *Core) error {
		if maxMemory < minimumMaxMemory {
			return errors.New("max memory too small")
		}
		c.maxMemory = maxMemory
		return nil
	}
}
