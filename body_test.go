package httpin

import (
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type LanguageLevel struct {
	Language string `json:"lang" xml:"lang"`
	Level    int    `json:"level" xml:"level"`
}

type BodyPayload struct {
	Name      string           `json:"name" xml:"name"`
	Age       int              `json:"age" xml:"age"`
	IsNative  bool             `json:"is_native" xml:"is_native"`
	Hobbies   []string         `json:"hobbies" xml:"hobbies"`
	Languages []*LanguageLevel `json:"languages" xml:"languages"`
}

type JSONBodyPayload struct {
	JSONBody
	BodyPayload
}

type XMLBodyPayload struct {
	XMLBody
	BodyPayload
}

func TestJSONBody(t *testing.T) {
	Convey("json: parse HTTP body in JSON", t, func() {
		resolver, err := buildResolverTree(reflect.TypeOf(JSONBodyPayload{}))
		So(err, ShouldBeNil)
		So(resolver, ShouldNotBeNil)
		r, _ := http.NewRequest("GET", "https://example.com", nil)

		r.Body = io.NopCloser(strings.NewReader(`{
			"name": "Elia",
			"is_native": false,
			"age": 14,
			"hobbies": ["Gaming", "Drawing"],
			"languages": [
				{"lang": "English", "level": 10},
				{"lang": "Japanese", "level": 3}
			]
		}`))
		res, err := resolver.resolve(r)
		So(err, ShouldBeNil)
		So(res.Interface(), ShouldResemble, &JSONBodyPayload{
			BodyPayload: BodyPayload{
				Name:     "Elia",
				Age:      14,
				IsNative: false,
				Hobbies:  []string{"Gaming", "Drawing"},
				Languages: []*LanguageLevel{
					{"English", 10},
					{"Japanese", 3},
				},
			},
		})
	})
}

func TestXMLBody(t *testing.T) {
	Convey("json: parse HTTP body in XML", t, func() {
		resolver, err := buildResolverTree(reflect.TypeOf(XMLBodyPayload{}))
		So(err, ShouldBeNil)
		So(resolver, ShouldNotBeNil)
		r, _ := http.NewRequest("GET", "https://example.com", nil)

		r.Body = io.NopCloser(strings.NewReader(`<BodyPayload>
		<name>Elia</name>
		<age>14</age>
		<is_native>false</is_native>
		<hobbies>Gaming</hobbies>
		<hobbies>Drawing</hobbies>
		<languages>
		   <lang>English</lang>
		   <level>10</level>
		</languages>
		<languages>
		   <lang>Japanese</lang>
		   <level>3</level>
		</languages>
	 </BodyPayload>`))
		res, err := resolver.resolve(r)
		So(err, ShouldBeNil)
		So(res.Interface(), ShouldResemble, &XMLBodyPayload{
			BodyPayload: BodyPayload{
				Name:     "Elia",
				Age:      14,
				IsNative: false,
				Hobbies:  []string{"Gaming", "Drawing"},
				Languages: []*LanguageLevel{
					{"English", 10},
					{"Japanese", 3},
				},
			},
		})
	})
}
