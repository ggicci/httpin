// directive: "body"
// https://ggicci.github.io/httpin/directives/body

package core

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/ggicci/httpin/internal"
)

// ErrUnknownBodyFormat is returned when a codec for the specified body format has not been specified.
var ErrUnknownBodyFormat = errors.New("unknown body format")

// DirectiveBody is the implementation of the "body" directive.
type DirectiveBody struct{}

func (db *DirectiveBody) Decode(rtm *DirectiveRuntime) error {
	req := rtm.GetRequest()
	bodyFormat, bodyCodec := db.getCodec(rtm)
	if bodyCodec == nil {
		return fmt.Errorf("%w: %q", ErrUnknownBodyFormat, bodyFormat)
	}
	if err := bodyCodec.Decode(req.Body, rtm.Value.Elem().Addr().Interface()); err != nil {
		return err
	}
	return nil
}

func (db *DirectiveBody) Encode(rtm *DirectiveRuntime) error {
	bodyFormat, bodyCodec := db.getCodec(rtm)
	if bodyCodec == nil {
		return fmt.Errorf("%w: %q", ErrUnknownBodyFormat, bodyFormat)
	}
	if bodyReader, err := bodyCodec.Encode(rtm.Value.Interface()); err != nil {
		return err
	} else {
		rtm.GetRequestBuilder().SetBody(bodyFormat, io.NopCloser(bodyReader))
		rtm.MarkFieldSet(true)
		return nil
	}
}

func (*DirectiveBody) getCodec(rtm *DirectiveRuntime) (bodyFormat string, bodyCodec BodyCodec) {
	bodyFormat = "json"
	if len(rtm.Directive.Argv) > 0 {
		bodyFormat = strings.ToLower(rtm.Directive.Argv[0])
	}
	bodyCodec = getBodyCodec(bodyFormat)
	return
}

var bodyFormats = map[string]BodyCodec{
	"json": &JSONBody{},
	"xml":  &XMLBody{},
}

// BodyCodec is the interface for encoding and decoding the HTTP request body.
// Common body formats are: json, xml, yaml, etc.
type BodyCodec interface {
	// Decode decodes the request body into the specified object.
	Decode(src io.Reader, dst any) error
	// Encode encodes the specified object into a reader for the request body.
	Encode(src any) (io.Reader, error)
}

// RegisterBodyFormat registers a new data formatter for the body request, which has the
// BodyCodec interface implemented. Panics on taken name, empty name or nil
// decoder. Pass parameter force (true) to ignore the name conflict.
//
// The BodyCodec is used by the body directive to decode and encode the data in
// the given format (body format).
//
// It is also useful when you want to override the default registered
// BodyCodec. For example, the default JSON decoder is borrowed from
// encoding/json. You can replace it with your own implementation, e.g.
// json-iterator/go. For example:
//
//	func init() {
//	    RegisterBodyFormat("json", &myJSONBody{}, true) // force register, replace the old one
//	    RegisterBodyFormat("yaml", &myYAMLBody{}) // register a new body format "yaml"
//	}
func RegisterBodyFormat(format string, codec BodyCodec, force ...bool) {
	internal.PanicOnError(
		registerBodyFormat(format, codec, force...),
	)
}

func getBodyCodec(bodyFormat string) BodyCodec {
	return bodyFormats[bodyFormat]
}

type JSONBody struct{}

func (de *JSONBody) Decode(src io.Reader, dst any) error {
	return json.NewDecoder(src).Decode(dst)
}

func (en *JSONBody) Encode(src any) (io.Reader, error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(src); err != nil {
		return nil, err
	}
	return &buf, nil
}

type XMLBody struct{}

func (de *XMLBody) Decode(src io.Reader, dst any) error {
	return xml.NewDecoder(src).Decode(dst)
}

func (en *XMLBody) Encode(src any) (io.Reader, error) {
	var buf bytes.Buffer
	if err := xml.NewEncoder(&buf).Encode(src); err != nil {
		return nil, err
	}
	return &buf, nil
}

func registerBodyFormat(format string, codec BodyCodec, force ...bool) error {
	ignoreConflict := len(force) > 0 && force[0]
	format = strings.ToLower(format)
	if !ignoreConflict && getBodyCodec(format) != nil {
		return fmt.Errorf("duplicate body format: %q", format)
	}
	if format == "" {
		return errors.New("body format cannot be empty")
	}
	if codec == nil {
		return errors.New("body codec cannot be nil")
	}
	bodyFormats[format] = codec
	return nil
}
