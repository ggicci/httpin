// directive: "path"
// https://ggicci.github.io/httpin/directives/path

package core

import "errors"

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

// defaultPathDirective is the default path directive, which only supports encoding,
// while the decoding function is not implmented. Because the path decoding depends on the
// routing framework, it should be implemented in the integration package.
// See integration/gochi.go and integration/gorilla.go for examples.
var defaultPathDirective = NewDirectivePath(func(rtm *DirectiveRuntime) error {
	return errors.New("unimplemented path decoding function")
})
