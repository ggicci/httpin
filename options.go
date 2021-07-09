package httpin

type option func(*Engine)

func WithErrorStatusCode(code int) option {
	return func(c *Engine) {
		c.errorStatusCode = code
	}
}
