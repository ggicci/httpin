// directive: "header"
// https://ggicci.github.io/httpin/directives/header

package directive

import (
	"mime/multipart"
	"net/http"
)

type DirectiveHeader struct{}

// Decode implements the "header" executor who extracts values
// from the HTTP headers.
func (*DirectiveHeader) Decode(rtm *DirectiveRuntime) error {
	req := rtm.GetRequest()
	Extractor := &Extractor{
		Runtime: rtm,
		Form: multipart.Form{
			Value: req.Header,
		},
		KeyNormalizer: http.CanonicalHeaderKey,
	}
	return Extractor.Extract()
}

func (*DirectiveHeader) Encode(rtm *DirectiveRuntime) error {
	encoder := &FormEncoder{
		Setter: rtm.GetRequestBuilder().SetHeader,
	}
	return encoder.Execute(rtm)
}
