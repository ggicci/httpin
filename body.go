// directive: "body"
// https://ggicci.github.io/httpin/directives/body

package httpin

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/ggicci/owl"
)

func init() {
	RegisterBodyFormat("json", &defaultJSONBody{})
	RegisterBodyFormat("xml", &defaultXMLBody{})
}

var bodyFormats = make(map[string]BodyEncodeDecoder)

// BodyEncodeDecoder is the interface for encoding and decoding the request body.
// Common body formats are: json, xml, yaml, etc.
type BodyEncodeDecoder interface {
	// Decode decodes the request body into the specified object.
	Decode(src io.Reader, dst any) error
	// Encode encodes the specified object into a reader for the request body.
	Encode(src any) (io.Reader, error)
}

// RegisterBodyFormat registers a new data formatter for the body request, which has the
// BodyEncoderDecoder interface implemented. Panics on taken name, empty name or nil
// decoder. Pass parameter replace (true) to ignore the name conflict.
//
// The BodyEncoderDecoder is used by the body directive to decode and encode the data in
// the given format (body format).
//
// It is also useful when you want to override the default registered
// BodyEncoderDecoder. For example, the default JSON decoder is borrowed from
// encoding/json. You can replace it with your own implementation, e.g.
// json-iterator/go. For example:
//
//	func init() {
//	    RegisterBodyFormat("json", &myJSONBody{}, true) // force register, replace the old one
//	    RegisterBodyFormat("yaml", &myYAMLBody{}) // register a new body format "yaml"
//	}
func RegisterBodyFormat(format string, body BodyEncodeDecoder, replace ...bool) {
	format = strings.ToLower(format)
	force := len(replace) > 0 && replace[0]
	if _, ok := bodyFormats[format]; ok && !force {
		panicOnError(fmt.Errorf("duplicate body format: %q", format))
	}
	if format == "" {
		panicOnError(errors.New("body format cannot be empty"))
	}
	if body == nil {
		panicOnError(errors.New("body encoder and decoder cannot be nil"))
	}
	bodyFormats[format] = body
}

// normalizeBodyDirective normalizes the body directive of the resolver.
// If no body format specified, the default type is "json".
func normalizeBodyDirective(r *owl.Resolver) error {
	dir := r.GetDirective("body")
	if dir == nil || dir.Name != "body" {
		return nil
	}
	if len(dir.Argv) == 0 {
		dir.Argv = []string{"json"} // use json as default when no body format specified
	}
	dir.Argv[0] = strings.ToLower(dir.Argv[0])

	var bodyFormat = dir.Argv[0]
	if _, ok := bodyFormats[bodyFormat]; !ok {
		return fmt.Errorf("unknown body format: %q", bodyFormat)
	}
	return nil
}

type defaultJSONBody struct{}

func (de *defaultJSONBody) Decode(src io.Reader, dst any) error {
	return json.NewDecoder(src).Decode(dst)
}

func (en *defaultJSONBody) Encode(src any) (io.Reader, error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(src); err != nil {
		return nil, err
	}
	return &buf, nil
}

type defaultXMLBody struct{}

func (de *defaultXMLBody) Decode(src io.Reader, dst any) error {
	return xml.NewDecoder(src).Decode(dst)
}

func (en *defaultXMLBody) Encode(src any) (io.Reader, error) {
	var buf bytes.Buffer
	if err := xml.NewEncoder(&buf).Encode(src); err != nil {
		return nil, err
	}
	return &buf, nil
}

// directiveBody is the implementation of the "body" directive.
type directiveBody struct{}

func (*directiveBody) Decode(rtm *DirectiveRuntime) error {
	var (
		req        = rtm.GetRequest()
		bodyFormat = rtm.Directive.Argv[0]
		decoder    = bodyFormats[bodyFormat]
	)
	if err := decoder.Decode(req.Body, rtm.Value.Elem().Addr().Interface()); err != nil {
		return err
	}
	return nil
}

func (*directiveBody) Encode(rtm *DirectiveRuntime) error {
	bodyFormat := rtm.Directive.Argv[0]
	bodyEncoder := bodyFormats[bodyFormat]
	if bodyReader, err := bodyEncoder.Encode(rtm.Value.Interface()); err != nil {
		return err
	} else {
		rtm.GetRequestBuilder().setBody(bodyFormat, io.NopCloser(bodyReader))
		rtm.MarkFieldSet(true)
		return nil
	}
}
