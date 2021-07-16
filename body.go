package httpin

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
	"reflect"
)

type JSONBody struct{}
type XMLBody struct{}

var (
	typeJSONBody = reflect.TypeOf(JSONBody{})
	typeXMLBody  = reflect.TypeOf(XMLBody{})
)

func decodeJSONBody(req *http.Request, rv reflect.Value) error {
	obj := rv.Interface()
	return json.NewDecoder(req.Body).Decode(&obj)
}

func decodeXMLBody(req *http.Request, rv reflect.Value) error {
	obj := rv.Interface()
	return xml.NewDecoder(req.Body).Decode(&obj)
}
