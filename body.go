// directive: "body"
// https://ggicci.github.io/httpin/directives/body

package httpin

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/ggicci/owl"
)

const (
	bodyTypeJSON = "json"
	bodyTypeXML  = "xml"
)

type (
	// BodyDecoder decodes the request body into the specified object. Common body types are:
	// json, xml, yaml, and others.
	BodyDecoder interface {
		Decode(src io.Reader, dst interface{}) error
	}
)

var (
	bodyDecoders = map[string]BodyDecoder{
		bodyTypeJSON: &defaultJSONBodyDecoder{},
		bodyTypeXML:  &defaultXMLBodyDecoder{},
	}
)

// RegisterBodyDecoder registers a new body decoder. Panic if the body type is already registered.
//
//	func init() {
//	    RegisterBodyDecoder("yaml", &myYAMLBodyDecoder{})
//	}
func RegisterBodyDecoder(bodyType string, decoder BodyDecoder) {
	if _, ok := bodyDecoders[bodyType]; ok {
		panic(fmt.Errorf("httpin: %w: %q", ErrDuplicateBodyDecoder, bodyType))
	}
	ReplaceBodyDecoder(bodyType, decoder)
}

// ReplaceBodyDecoder replaces or add the body decoder of the specified type.
// Which is useful when you want to override the default body decoder. For example,
// the default JSON decoder is borrowed from encoding/json. You can replace it with
// your own implementation, e.g. json-iterator/go.
//
//	func init() {
//	    ReplaceBodyDecoder("json", &myJSONBodyDecoder{})
//	}
func ReplaceBodyDecoder(bodyType string, decoder BodyDecoder) {
	if bodyType == "" {
		panic("httpin: body type cannot be empty")
	}
	bodyDecoders[bodyType] = decoder
}

// normalizeBodyDirective normalizes the body directive of the resolver.
// If no body type specified, the default type is "json".
func normalizeBodyDirective(r *owl.Resolver) error {
	dir := r.GetDirective("body")
	if dir == nil || dir.Name != "body" {
		return nil
	}
	if len(dir.Argv) == 0 {
		dir.Argv = []string{bodyTypeJSON}
	}
	dir.Argv[0] = strings.ToLower(dir.Argv[0])

	var bodyType = dir.Argv[0]
	if _, ok := bodyDecoders[bodyType]; !ok {
		return fmt.Errorf("%w: %q", ErrUnknownBodyType, bodyType)
	}
	return nil
}

func bodyDecoder(rtm *DirectiveRuntime) error {
	var (
		req      = rtm.Context.Value(RequestValue).(*http.Request)
		bodyType = rtm.Directive.Argv[0]
		decoder  = bodyDecoders[bodyType]
	)
	if err := decoder.Decode(req.Body, rtm.Value.Elem().Addr().Interface()); err != nil {
		return err
	}
	return nil
}

type defaultJSONBodyDecoder struct{}

func (de *defaultJSONBodyDecoder) Decode(src io.Reader, dst interface{}) error {
	return json.NewDecoder(src).Decode(dst)
}

type defaultXMLBodyDecoder struct{}

func (de *defaultXMLBodyDecoder) Decode(src io.Reader, dst interface{}) error {
	return xml.NewDecoder(src).Decode(dst)
}
