package core

// DirectiveNoop is a DirectiveExecutor that does nothing, "noop" stands for "no operation".
type DirectiveNoop struct{}

func (*DirectiveNoop) Encode(*DirectiveRuntime) error { return nil }
func (*DirectiveNoop) Decode(*DirectiveRuntime) error { return nil }
