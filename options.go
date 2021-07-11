package httpin

type option func(*Engine)

// WithErrorStatusCode configures the HTTP status code sent to the client when
// decoding a request failed. Which is used in the `NewInput` middleware.
// The default value is 422.
func WithErrorStatusCode(code int) option {
	return func(c *Engine) {
		c.errorStatusCode = code
	}
}
