// directive: "body"
// https://ggicci.github.io/httpin/directives/body

package httpin

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"reflect"
	"strings"
)

type JSONBody struct{}
type XMLBody struct{}

var (
	bodyTypeAnnotationJSON = reflect.TypeOf(JSONBody{})
	bodyTypeAnnotationXML  = reflect.TypeOf(XMLBody{})
)

const (
	bodyTypeJSON = "json"
	bodyTypeXML  = "xml"
)

func bodyDirectiveNormalizer(dir *Directive) error {
	if len(dir.Argv) == 0 {
		dir.Argv = []string{bodyTypeJSON}
	}
	dir.Argv[0] = strings.ToLower(dir.Argv[0])

	var bodyType = dir.Argv[0]
	if bodyType != bodyTypeJSON && bodyType != bodyTypeXML {
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
		panic(ErrUnknownBodyType)
	}
}

func bodyDecoder(ctx *DirectiveContext) error {
	var err = ErrUnknownBodyType
	switch ctx.Argv[0] { // body type
	case bodyTypeJSON:
		err = decodeJSONBody(ctx.Request, ctx.Value)
	case bodyTypeXML:
		err = decodeXMLBody(ctx.Request, ctx.Value)
	}

	if err != nil {
		return err
	}

	ctx.DeliverContextValue(StopRecursion, true)
	return nil
}

func decodeJSONBody(req *http.Request, rv reflect.Value) error {
	obj := rv.Interface()
	return json.NewDecoder(req.Body).Decode(&obj)
}

func decodeXMLBody(req *http.Request, rv reflect.Value) error {
	obj := rv.Interface()
	return xml.NewDecoder(req.Body).Decode(&obj)
}
