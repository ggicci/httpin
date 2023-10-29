// directive: "query"
// https://ggicci.github.io/httpin/directives/query

package httpin

import (
	"mime/multipart"
)

type directiveQuery struct{}

// Decode implements the "query" executor who extracts values from
// the querystring of an HTTP request.
func (*directiveQuery) Decode(rtm *DirectiveRuntime) error {
	req := rtm.GetRequest()
	extractor := &extractor{
		Runtime: rtm,
		Form: multipart.Form{
			Value: req.URL.Query(),
		},
	}
	return extractor.Extract()
}

func (*directiveQuery) Encode(rtm *DirectiveRuntime) error {
	encoder := &formEncoder{rtm.GetRequestBuilder().setQuery}
	return encoder.Execute(rtm)
}
