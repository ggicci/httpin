package httpin

type option func(*core)

func WithErrorStatusCode(code int) option {
	return func(c *core) {
		c.errorStatusCode = code
	}
}
