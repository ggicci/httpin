// directive: "body"
// https://ggicci.github.io/httpin/directives/body

package directive

import (
	"io"

	"github.com/ggicci/httpin/internal"
)

// DirectiveBody is the implementation of the "body" directive.
type DirectiveBody struct{}

func (*DirectiveBody) Decode(rtm *DirectiveRuntime) error {
	req := rtm.GetRequest()
	bodyFormat := rtm.Directive.Argv[0]
	bodyDecoder := internal.DefaultRegistry.GetBodyDecoder(bodyFormat)
	if err := bodyDecoder.Decode(req.Body, rtm.Value.Elem().Addr().Interface()); err != nil {
		return err
	}
	return nil
}

func (*DirectiveBody) Encode(rtm *DirectiveRuntime) error {
	bodyFormat := rtm.Directive.Argv[0]
	bodyEncoder := internal.DefaultRegistry.GetBodyDecoder(bodyFormat)
	if bodyReader, err := bodyEncoder.Encode(rtm.Value.Interface()); err != nil {
		return err
	} else {
		rtm.GetRequestBuilder().SetBody(bodyFormat, io.NopCloser(bodyReader))
		rtm.MarkFieldSet(true)
		return nil
	}
}
