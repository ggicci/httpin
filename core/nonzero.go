// directive: "nonzero"
// https://ggicci.github.io/httpin/directives/nonzero

package core

import "errors"

// DirectiveNonzero implements the "nonzero" executor who indicates that the field must not be a "zero value".
// In golang, the "zero value" means:
//   - nil
//   - false
//   - 0
//   - ""
//   - etc.
//
// Unlike the "required" executor, the "nonzero" executor checks the value of the field.
type DirectiveNonzero struct{}

func (*DirectiveNonzero) Decode(rtm *DirectiveRuntime) error {
	if rtm.Value.Elem().IsZero() {
		return errors.New("zero value")
	}
	return nil
}

func (*DirectiveNonzero) Encode(rtm *DirectiveRuntime) error {
	if rtm.Value.IsZero() {
		return errors.New("zero value")
	}
	return nil
}
