package core

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/ggicci/httpin/internal"
)

type FileMarshaler = internal.FileMarshaler

type RequestBuilder struct {
	Query      url.Values
	Form       url.Values
	Attachment map[string][]FileMarshaler
	Header     http.Header
	Cookie     []*http.Cookie
	Path       map[string]string // placeholder: value
	BodyType   string            // json, xml, etc.
	Body       io.ReadCloser
	ctx        context.Context
}

func NewRequestBuilder(ctx context.Context) *RequestBuilder {
	return &RequestBuilder{
		Query:      make(url.Values),
		Form:       make(url.Values),
		Attachment: make(map[string][]FileMarshaler),
		Header:     make(http.Header),
		Cookie:     make([]*http.Cookie, 0),
		Path:       make(map[string]string),
		ctx:        ctx,
	}
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
			rb.populateForm(req)
		}
	}

	// Populate body.
	if rb.hasBody() {
		req.Body = rb.Body
		rb.Header.Set("Content-Type", rb.bodyContentType())
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
	rb.Query[key] = value
}

func (rb *RequestBuilder) SetForm(key string, value []string) {
	rb.Form[key] = value
}

func (rb *RequestBuilder) SetHeader(key string, value []string) {
	rb.Header[http.CanonicalHeaderKey(key)] = value
}

func (rb *RequestBuilder) SetPath(key string, value []string) {
	if len(value) > 0 {
		rb.Path[key] = value[0]
	}
}

func (rb *RequestBuilder) SetBody(bodyType string, bodyReader io.ReadCloser) {
	rb.BodyType = bodyType
	rb.Body = bodyReader
}

func (rb *RequestBuilder) SetAttachment(key string, files []FileMarshaler) {
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
	if rb.hasForm() && rb.hasBody() {
		return errors.New("cannot use both form and body directive at the same time")
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

func (rb *RequestBuilder) populateForm(req *http.Request) {
	rb.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	formData := rb.Form.Encode()
	req.Body = io.NopCloser(strings.NewReader(formData))
}

func (rb *RequestBuilder) populateMultipartForm(req *http.Request) error {
	// Create a pipe and a multipart writer.
	pr, pw := io.Pipe()
	writer := multipart.NewWriter(pw)

	// Write the multipart form data to the pipe in a separate goroutine.
	go func() {
		defer pw.Close()
		defer writer.Close()

		// Populate the form fields.
		for k, v := range rb.Form {
			for _, sv := range v {
				select {
				case <-rb.ctx.Done():
					pw.CloseWithError(rb.ctx.Err())
					return
				default:
					fieldWriter, _ := writer.CreateFormField(k)
					fieldWriter.Write([]byte(sv))
				}
			}
		}

		// Populate the attachments.
		for key, files := range rb.Attachment {
			for i, file := range files {
				select {
				case <-rb.ctx.Done():
					pw.CloseWithError(rb.ctx.Err())
					return
				default:
					filename := file.Filename()
					contentReader, err := file.MarshalFile()
					filename = normalizeUploadFilename(key, filename, i)

					if err != nil {
						pw.CloseWithError(fmt.Errorf("upload %s %q failed: %w", key, filename, err))
						return
					}

					fileWriter, _ := writer.CreateFormFile(key, filename)
					if _, err = io.Copy(fileWriter, contentReader); err != nil {
						pw.CloseWithError(fmt.Errorf("upload %s %q failed: %w", key, filename, err))
						return
					}
				}
			}
		}
	}()

	// Set the body to the read end of the pipe and the content type.
	req.Body = io.NopCloser(pr)
	rb.Header.Set("Content-Type", writer.FormDataContentType())
	return nil
}

func normalizeUploadFilename(key, filename string, index int) string {
	if filename == "" {
		return fmt.Sprintf("%s_%d", key, index)
	}
	return filepath.Base(filename)
}
