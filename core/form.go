// directive: "form"
// https://ggicci.github.io/httpin/directives/form

package core

import (
	"mime/multipart"
)

type DirectvieForm struct{}

// Decode implements the "form" executor who extracts values from
// the forms of an HTTP request.
func (*DirectvieForm) Decode(rtm *DirectiveRuntime) error {
	req := rtm.GetRequest()
	var form multipart.Form
	if req.MultipartForm != nil {
		form = *req.MultipartForm
	} else {
		if req.Form != nil {
			form.Value = req.Form
		}
	}
	extractor := &FormExtractor{
		Runtime:       rtm,
		Form:          form,
		KeyNormalizer: nil,
	}
	return extractor.Extract()
}

// Encode implements the encoder/request builder for "form" directive.
// It builds the form values of an HTTP request, including:
//   - form data
//   - multipart form data (file upload)
func (*DirectvieForm) Encode(rtm *DirectiveRuntime) error {
	encoder := &FormEncoder{
		Setter: rtm.GetRequestBuilder().SetForm,
	}
	return encoder.Execute(rtm)
}
