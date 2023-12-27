// directive: "query"
// https://ggicci.github.io/httpin/directives/query

package core

import (
	"mime/multipart"
)

type DirectiveQuery struct{}

// Decode implements the "query" executor who extracts values from
// the querystring of an HTTP request.
func (*DirectiveQuery) Decode(rtm *DirectiveRuntime) error {
	req := rtm.GetRequest()
	extractor := &FormExtractor{
		Runtime: rtm,
		Form: multipart.Form{
			Value: req.URL.Query(),
		},
	}
	return extractor.Extract()
}

func (*DirectiveQuery) Encode(rtm *DirectiveRuntime) error {
	encoder := &FormEncoder{
		Setter: rtm.GetRequestBuilder().SetQuery,
	}
	return encoder.Execute(rtm)
}
