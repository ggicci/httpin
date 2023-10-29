// directive: "path"
// https://ggicci.github.io/httpin/directives/path

package httpin

type directivePath struct {
	overrideDecode func(*DirectiveRuntime) error
}

func (dir *directivePath) Decode(rtm *DirectiveRuntime) error {
	return dir.overrideDecode(rtm)
}

// Encode replaces the placeholders in URL path with the given value.
func (*directivePath) Encode(rtm *DirectiveRuntime) error {
	encoder := &formEncoder{rtm.GetRequestBuilder().setPath}
	return encoder.Execute(rtm)
}
