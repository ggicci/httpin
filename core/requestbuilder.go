package core

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
)

type RequestBuilder struct {
	Query      url.Values
	Form       url.Values
	Attachment map[string][]FileEncoder
	Header     http.Header
	Cookie     []*http.Cookie
	Path       map[string]string // placeholder: value
	BodyType   string            // json, xml, etc.
	Body       io.ReadCloser
}

func (rb *RequestBuilder) Populate(req *http.Request) error {
	if err := rb.validate(); err != nil {
		return err
	}

	// Populate the querystring.
	req.URL.RawQuery = rb.Query.Encode()

	// Populate forms.
	if rb.hasForm() {
		if rb.hasAttachment() { // multipart form
			if err := rb.populateMultipartForm(req); err != nil {
				return err
			}
		} else { // urlencoded form
			req.Form = rb.Form
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		}
	}

	// Populate body.
	if rb.hasBody() {
		req.Body = rb.Body
		req.Header.Set("Content-Type", rb.bodyContentType())
	}

	// Populate path.
	if rb.hasPath() {
		newPath := req.URL.Path
		for key, value := range rb.Path {
			newPath = strings.Replace(newPath, "{"+key+"}", value, -1)
		}
		req.URL.Path = newPath
		req.URL.RawPath = ""
	}

	// Populate the headers.
	if rb.Header != nil {
		req.Header = rb.Header
	}

	// Populate the cookies.
	for _, cookie := range rb.Cookie {
		req.AddCookie(cookie)
	}

	return nil
}

func (rb *RequestBuilder) SetQuery(key string, value []string) {
	if rb.Query == nil {
		rb.Query = make(url.Values)
	}
	rb.Query[key] = value
}

func (rb *RequestBuilder) SetForm(key string, value []string) {
	if rb.Form == nil {
		rb.Form = make(url.Values)
	}
	rb.Form[key] = value
}

func (rb *RequestBuilder) SetHeader(key string, value []string) {
	if rb.Header == nil {
		rb.Header = make(http.Header)
	}
	rb.Header[http.CanonicalHeaderKey(key)] = value
}

func (rb *RequestBuilder) SetPath(key string, value []string) {
	if rb.Path == nil {
		rb.Path = make(map[string]string)
	}
	if len(value) > 0 {
		rb.Path[key] = value[0]
	}
}

func (rb *RequestBuilder) SetBody(bodyType string, bodyReader io.ReadCloser) {
	rb.BodyType = bodyType
	rb.Body = bodyReader
}

func (rb *RequestBuilder) SetAttachment(key string, files []FileEncoder) {
	if rb.Attachment == nil {
		rb.Attachment = make(map[string][]FileEncoder)
	}
	rb.Attachment[key] = files
}

func (rb *RequestBuilder) bodyContentType() string {
	switch rb.BodyType {
	case "json":
		return "application/json"
	case "xml":
		return "application/xml"
	}
	return ""
}

func (rb *RequestBuilder) validate() error {
	if rb.hasAttachment() && rb.hasBody() {
		return errors.New("cannot use body directive and file upload at the same time")
	}
	return nil
}

func (rb *RequestBuilder) hasPath() bool {
	return len(rb.Path) > 0
}

func (rb *RequestBuilder) hasForm() bool {
	return len(rb.Form) > 0 || rb.hasAttachment()
}

func (rb *RequestBuilder) hasAttachment() bool {
	return len(rb.Attachment) > 0
}

func (rb *RequestBuilder) hasBody() bool {
	return rb.Body != nil && rb.BodyType != ""
}

func (rb *RequestBuilder) populateMultipartForm(req *http.Request) error {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	// Populate the form fields.
	for k, v := range rb.Form {
		for _, sv := range v {
			fieldWriter, _ := writer.CreateFormField(k)
			fieldWriter.Write([]byte(sv))
		}
	}

	// Populate the attachments.
	for key, files := range rb.Attachment {
		for i, file := range files {
			filename, contentReader, err := file.Encode()
			filename = normalizeUploadFilename(key, filename, i)

			if err != nil {
				return fmt.Errorf("upload %s %q: %w", key, filename, err)
			}

			fileWriter, _ := writer.CreateFormFile(key, filename)
			if _, err = io.Copy(fileWriter, contentReader); err != nil {
				return fmt.Errorf("upload %s %q: %w", key, filename, err)
			}
		}
	}

	if err := writer.Close(); err != nil {
		return fmt.Errorf("close multipart writer: %w", err)
	}

	// Set the body and content type.
	req.Body = io.NopCloser(body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return nil
}

func normalizeUploadFilename(key, filename string, index int) string {
	if filename == "" {
		return fmt.Sprintf("%s_%d", key, index)
	}
	return filepath.Base(filename)
}
