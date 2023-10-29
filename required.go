// directive: "required"
// https://ggicci.github.io/httpin/directives/required

package httpin

// directiveRequired implements the "required" executor who indicates that the field must be set.
// If the field value were not set by former executors, errMissingField will be
// returned.
//
// NOTE: the "required" executor does not check the value of the field, it only checks
// if the field is set. In realcases, it's used to require that the key is present in
// the input data, e.g. form, header, etc. But it allows the value to be empty.
type directiveRequired struct{}

func (*directiveRequired) Decode(rtm *DirectiveRuntime) error {
	if rtm.IsFieldSet() {
		return nil
	}
	return errMissingField
}

func (*directiveRequired) Encode(rtm *DirectiveRuntime) error {
	return nil // noop
}
