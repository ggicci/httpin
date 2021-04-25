package httpin

type EngineOption func(*Engine)

func WithQueryTag(tag string) EngineOption {
	return func(e *Engine) {
		e.queryTag = tag
	}
}

func WithHeaderTag(tag string) EngineOption {
	return func(e *Engine) {
		e.headerTag = tag
	}
}

func WithBodyTag(tag string) EngineOption {
	return func(e *Engine) {
		e.bodyTag = tag
	}
}
