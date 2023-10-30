// directive: "form"
// https://ggicci.github.io/httpin/directives/form

package httpin

import "mime/multipart"

type directvieForm struct{}

// Decode implements the "form" executor who extracts values from
// the forms of an HTTP request.
func (*directvieForm) Decode(rtm *DirectiveRuntime) error {
	req := rtm.GetRequest()
	var form multipart.Form
	if req.MultipartForm != nil {
		form = *req.MultipartForm
	} else {
		if req.Form != nil {
			form.Value = req.Form
		}
	}
	extractor := &extractor{
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
func (*directvieForm) Encode(rtm *DirectiveRuntime) error {
	encoder := &formEncoder{
		Setter: rtm.GetRequestBuilder().setForm,
	}
	return encoder.Execute(rtm)
}
