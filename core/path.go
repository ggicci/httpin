// directive: "path"
// https://ggicci.github.io/httpin/directives/path

package core

type DirectivePath struct {
	decode func(*DirectiveRuntime) error
}

func NewDirectivePath(decodeFunc func(*DirectiveRuntime) error) *DirectivePath {
	return &DirectivePath{
		decode: decodeFunc,
	}
}

func (dir *DirectivePath) Decode(rtm *DirectiveRuntime) error {
	return dir.decode(rtm)
}

// Encode replaces the placeholders in URL path with the given value.
func (*DirectivePath) Encode(rtm *DirectiveRuntime) error {
	encoder := &FormEncoder{
		Setter: rtm.GetRequestBuilder().SetPath,
	}
	return encoder.Execute(rtm)
}
