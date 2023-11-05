// directive: "path"
// https://ggicci.github.io/httpin/directives/path

package core

type DirectivePath struct {
	overrideDecode func(*DirectiveRuntime) error
}

func NewDirectivePath(decoder func(*DirectiveRuntime) error) *DirectivePath {
	return &DirectivePath{
		overrideDecode: decoder,
	}
}

func (dir *DirectivePath) Decode(rtm *DirectiveRuntime) error {
	return dir.overrideDecode(rtm)
}

// Encode replaces the placeholders in URL path with the given value.
func (*DirectivePath) Encode(rtm *DirectiveRuntime) error {
	encoder := &FormEncoder{
		Setter: rtm.GetRequestBuilder().SetPath,
	}
	return encoder.Execute(rtm)
}
