// directive: "header"
// https://ggicci.github.io/httpin/directives/header

package httpin

import (
	"mime/multipart"
	"net/http"
)

type directiveHeader struct{}

// Decode implements the "header" executor who extracts values
// from the HTTP headers.
func (*directiveHeader) Decode(rtm *DirectiveRuntime) error {
	req := rtm.GetRequest()
	extractor := &extractor{
		Runtime: rtm,
		Form: multipart.Form{
			Value: req.Header,
		},
		KeyNormalizer: http.CanonicalHeaderKey,
	}
	return extractor.Extract()
}

func (*directiveHeader) Encode(rtm *DirectiveRuntime) error {
	encoder := &formEncoder{
		Setter: rtm.GetRequestBuilder().setHeader,
	}
	return encoder.Execute(rtm)
}
