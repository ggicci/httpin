package core

import (
	"errors"
)

const minimumMaxMemory = int64(1 << 10)  // 1KB
const defaultMaxMemory = int64(32 << 20) // 32 MB

type Option func(*Core) error

var globalNestedDirectivesEnabled bool = false

// EnableNestedDirectives sets the global flag to enable nested directives.
// Nested directives are disabled by default.
func EnableNestedDirectives(on bool) {
	globalNestedDirectivesEnabled = on
}

// WithErrorHandler overrides the default error handler.
func WithErrorHandler(custom ErrorHandler) Option {
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

// WithNestedDirectivesEnabled enables/disables nested directives.
func WithNestedDirectivesEnabled(enable bool) Option {
	return func(c *Core) error {
		c.enableNestedDirectives = enable
		return nil
	}
}
