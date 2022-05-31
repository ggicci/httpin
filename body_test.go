package httpin

import (
	"errors"
	"io"
	"net/http"
	"net/url"
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

type JSONBodyPayloadWithAnnotation struct {
	JSONBody
	BodyPayload
}

type XMLBodyPayloadWithAnnotation struct {
	XMLBody
	BodyPayload
}

type JSONBodyPayloadWithBodyDirective struct {
	Page     int          `in:"form=page"`
	PageSize int          `in:"form=page_size"`
	Body     *BodyPayload `in:"body=json"`
}

type ThingWithDuplicateAnnotations struct {
	JSONBody
	XMLBody
	Page int `in:"form=page"`
}

type ThingWithInvalidBodyType struct {
	Username string `in:"form=username"`

	Patch map[string]interface{} `in:"body=yaml"`
}

type ThingWithEmptyBodyType struct {
	Username string `in:"form=username"`

	Patch map[string]interface{} `in:"body"`
}

func TestAnnotationField(t *testing.T) {
	Convey("annotate: duplicate annotations", t, func() {
		resolver, err := buildResolverTree(reflect.TypeOf(ThingWithDuplicateAnnotations{}))
		So(resolver, ShouldBeNil)
		So(err, ShouldBeError)
		So(errors.Is(err, ErrDuplicateAnnotationField), ShouldBeTrue)
	})
}

func TestNormalizeBodyDirective(t *testing.T) {
	Convey("body directive: empty body type defaults to json", t, func() {
		resolver, err := buildResolverTree(reflect.TypeOf(ThingWithEmptyBodyType{}))
		So(err, ShouldBeNil)
		So(resolver.Fields[1].Directives[0].Argv[0], ShouldEqual, "json")
	})

	Convey("body directive: unknown body type", t, func() {
		resolver, err := buildResolverTree(reflect.TypeOf(ThingWithInvalidBodyType{}))
		So(resolver, ShouldBeNil)
		So(err, ShouldBeError)
		So(errors.Is(err, ErrUnknownBodyType), ShouldBeTrue)
	})
}

func TestJSONBody(t *testing.T) {
	Convey("body: parse HTTP body (in JSON) into a field of the input struct", t, func() {
		resolver, err := buildResolverTree(reflect.TypeOf(JSONBodyPayloadWithBodyDirective{}))
		So(err, ShouldBeNil)
		So(resolver, ShouldNotBeNil)
		r, _ := http.NewRequest("GET", "https://example.com", nil)

		r.Form = make(url.Values)
		r.Form.Set("page", "4")
		r.Form.Set("page_size", "30")
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
		So(res.Interface(), ShouldResemble, &JSONBodyPayloadWithBodyDirective{
			Page:     4,
			PageSize: 30,
			Body: &BodyPayload{
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

	Convey("body: parse HTTP body (in JSON) to the input struct, use annotation", t, func() {
		resolver, err := buildResolverTree(reflect.TypeOf(JSONBodyPayloadWithAnnotation{}))
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
		So(res.Interface(), ShouldResemble, &JSONBodyPayloadWithAnnotation{
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
	Convey("body: parse HTTP body (in XML) to the input struct, use annotation", t, func() {
		resolver, err := buildResolverTree(reflect.TypeOf(XMLBodyPayloadWithAnnotation{}))
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
		So(res.Interface(), ShouldResemble, &XMLBodyPayloadWithAnnotation{
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

func TestBodyDecoderDecodeFailed(t *testing.T) {
	Convey("body: parse request body in JSON failed", t, func() {
		resolver, err := buildResolverTree(reflect.TypeOf(JSONBodyPayloadWithAnnotation{}))
		So(err, ShouldBeNil)
		So(resolver, ShouldNotBeNil)
		r, _ := http.NewRequest("GET", "https://example.com", nil)

		r.Body = io.NopCloser(strings.NewReader(`{"name": "Elia"`))
		_, err = resolver.resolve(r)
		So(err, ShouldBeError)
	})
}

func Test_bodyTypeString(t *testing.T) {
	Convey("bodyTypeString should panic on unknown body type", t, func() {
		type yamlBody struct{}
		So(func() {
			bodyTypeString(reflect.TypeOf(yamlBody{}))
		}, ShouldPanic)
	})
}

type yamlBodyDecoder struct{}

func (de *yamlBodyDecoder) Decode(src io.Reader, dst interface{}) error {
	// for test only
	(*(*(dst.(*interface{}))).(*map[string]interface{})) = map[string]interface{}{
		"version": 3,
	}
	return nil
}

type YamlInput struct {
	Body map[string]interface{} `in:"body=yaml"`
}

type ThingWithUnknownBodyDecoder struct {
	Body map[string]interface{} `in:"body=yml"`
}

func TestCustomBodyDecoder(t *testing.T) {
	Convey("body: register new body decoder", t, func() {
		So(func() { RegisterBodyDecoder("yaml", &yamlBodyDecoder{}) }, ShouldNotPanic)

		resolver, err := buildResolverTree(reflect.TypeOf(YamlInput{}))
		So(err, ShouldBeNil)
		So(resolver, ShouldNotBeNil)
		r, _ := http.NewRequest("GET", "https://example.com", nil)

		r.Body = io.NopCloser(strings.NewReader(`version: "3"`))
		res, err := resolver.resolve(r)
		So(err, ShouldBeNil)
		So(res.Interface(), ShouldResemble, &YamlInput{
			Body: map[string]interface{}{
				"version": 3,
			},
		})
	})

	Convey("body: panic on duplicate body decoder", t, func() {
		So(func() { RegisterBodyDecoder("json", &yamlBodyDecoder{}) }, ShouldPanic)
		So(func() { RegisterBodyDecoder("", &yamlBodyDecoder{}) }, ShouldPanic)
	})

	Convey("body: unknown body decoder", t, func() {
		_, err := buildResolverTree(reflect.TypeOf(ThingWithUnknownBodyDecoder{}))
		So(errors.Is(err, ErrUnknownBodyType), ShouldBeTrue)
	})
}
