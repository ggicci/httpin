// directive: "body"
// https://ggicci.github.io/httpin/directives/body

package httpin

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"reflect"
	"strings"
)

const (
	bodyTypeJSON = "json"
	bodyTypeXML  = "xml"
)

type (
	JSONBody struct{}
	XMLBody  struct{}

	BodyDecoder interface {
		Decode(src io.Reader, dst interface{}) error
	}
)

var (
	bodyTypeAnnotationJSON = reflect.TypeOf(JSONBody{})
	bodyTypeAnnotationXML  = reflect.TypeOf(XMLBody{})

	bodyDecoders = map[string]BodyDecoder{
		bodyTypeJSON: &defaultJSONBodyDecoder{},
		bodyTypeXML:  &defaultXMLBodyDecoder{},
	}
)

func RegisterBodyDecoder(bodyType string, decoder BodyDecoder) {
	if _, ok := bodyDecoders[bodyType]; ok {
		panic(fmt.Errorf("httpin: %w: %q", ErrDuplicateBodyDecoder, bodyType))
	}
	ReplaceBodyDecoder(bodyType, decoder)
}

func ReplaceBodyDecoder(bodyType string, decoder BodyDecoder) {
	if bodyType == "" {
		panic("httpin: body type cannot be empty")
	}
	bodyDecoders[bodyType] = decoder
}

func bodyDirectiveNormalizer(dir *Directive) error {
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

func bodyTypeString(bodyType reflect.Type) string {
	switch bodyType {
	case bodyTypeAnnotationJSON:
		return bodyTypeJSON
	case bodyTypeAnnotationXML:
		return bodyTypeXML
	default:
		panic(fmt.Errorf("httpin: %w: %q", ErrUnknownBodyType, bodyType))
	}
}

func bodyDecoder(ctx *DirectiveContext) error {
	var (
		bodyType = ctx.Argv[0]
		decoder  = bodyDecoders[bodyType]
	)

	if decoder == nil {
		return ErrUnknownBodyType
	}

	obj := ctx.Value.Interface()
	if err := decoder.Decode(ctx.Request.Body, &obj); err != nil {
		return err
	}

	ctx.DeliverContextValue(StopRecursion, true)
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
