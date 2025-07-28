package core

type DirectiveExecutor interface {
	// Encode encodes the field of the input struct to the HTTP request.
	Encode(*DirectiveRuntime) error

	// Decode decodes the field of the input struct from the HTTP request.
	Decode(*DirectiveRuntime) error
}
