package core

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
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

type JSONBodyPayloadWithBodyDirective struct {
	Page     int          `in:"form=page"`
	PageSize int          `in:"form=page_size"`
	Body     *BodyPayload `in:"body=json"`
}

type XMLBodyPayloadWithBodyDirective struct {
	Body *BodyPayload `in:"body=xml"`
}

var sampleJSON_JSONBodyPayloadWithBodyDirective = `
{
	"name": "Elia",
	"is_native": false,
	"age": 14,
	"hobbies": ["Gaming", "Drawing"],
	"languages": [
		{"lang": "English", "level": 10},
		{"lang": "Japanese", "level": 3}
	]
}`

var sampleObject_JSONBodyPayloadWithBodyDirective = &JSONBodyPayloadWithBodyDirective{
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
}

var sampleXML_XMLBodyPayloadWithBodyDirective = `
<BodyPayload>
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
</BodyPayload>`

var sampleObject_XMLBodyPayloadWithBodyDirective = &XMLBodyPayloadWithBodyDirective{
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
}

func TestBodyDecoder_JSON(t *testing.T) {
	assert := assert.New(t)
	co, err := New(JSONBodyPayloadWithBodyDirective{})
	assert.NoError(err)

	r, _ := http.NewRequest("GET", "https://example.com", nil)
	r.Form = make(url.Values)
	r.Form.Set("page", "4")
	r.Form.Set("page_size", "30")
	r.Body = io.NopCloser(strings.NewReader(sampleJSON_JSONBodyPayloadWithBodyDirective))
	r.Header.Set("Content-Type", "application/json")
	gotValue, err := co.Decode(r)
	assert.NoError(err)
	assert.Equal(sampleObject_JSONBodyPayloadWithBodyDirective, gotValue)
}

func TestBodyDecoder_XML(t *testing.T) {

	assert := assert.New(t)
	co, err := New(XMLBodyPayloadWithBodyDirective{})
	assert.NoError(err)

	r, _ := http.NewRequest("GET", "https://example.com", nil)
	r.Body = io.NopCloser(strings.NewReader(sampleXML_XMLBodyPayloadWithBodyDirective))
	r.Header.Set("Content-Type", "application/xml")

	gotValue, err := co.Decode(r)
	assert.NoError(err)
	assert.Equal(sampleObject_XMLBodyPayloadWithBodyDirective, gotValue)
}

// func TestBodyDecoder_DefaultsToJSON(t *testing.T) {
// 	type Payload struct {
// 		Body *BodyPayload `in:"body"`
// 	}

// 	co, err :=New(Payload{})
// 	assert.NoError(t, err)
// 	d := core.resolver.Lookup("Body").GetDirective("body")
// 	assert.Equal(t, "json", d.Argv[0])
// }

func TestBodyDecoder_ErrUnknownBodyFormat(t *testing.T) {
	type UnknownBodyFormatPayload struct {
		Body *BodyPayload `in:"body=yaml"`
	}

	co, err := New(UnknownBodyFormatPayload{})
	assert.NoError(t, err)
	req, _ := http.NewRequest("GET", "https://example.com", nil)
	req.Body = io.NopCloser(strings.NewReader(sampleJSON_JSONBodyPayloadWithBodyDirective))
	_, err = co.Decode(req)
	assert.ErrorContains(t, err, "unknown body format: \"yaml\"")
}

type yamlBody struct{}

var errYamlNotImplemented = errors.New("yaml not implemented")

func (de *yamlBody) Decode(src io.Reader, dst any) error {
	return errYamlNotImplemented // for test only
}

func (en *yamlBody) Encode(src any) (io.Reader, error) {
	return nil, errYamlNotImplemented // for test only
}

type YamlInput struct {
	Body map[string]any `in:"body=yaml"`
}

func TestRegisterBody(t *testing.T) {
	assert.NotPanics(t, func() {
		RegisterBodyFormat("yaml", &yamlBody{})
	})
	assert.Panics(t, func() {
		RegisterBodyFormat("yaml", &yamlBody{})
	})

	co, err := New(YamlInput{})
	assert.NoError(t, err)

	r, _ := http.NewRequest("GET", "https://example.com", nil)
	r.Body = io.NopCloser(strings.NewReader(`version: "3"`))

	gotValue, err := co.Decode(r)
	assert.ErrorIs(t, err, errYamlNotImplemented)
	assert.Nil(t, gotValue)
}

func TestRegisterBody_nil(t *testing.T) {
	assert.Panics(t, func() {
		RegisterBodyFormat("toml", nil)
	})
}

func TestRegisterBody_forceReplace(t *testing.T) {
	assert.NotPanics(t, func() {
		RegisterBodyFormat("yaml", &yamlBody{}, true)
	})
	assert.NotPanics(t, func() {
		RegisterBodyFormat("yaml", &yamlBody{}, true)
	})
}

func TestRegisterBody_forceReplace_withEmptyBodyFormat(t *testing.T) {
	assert.PanicsWithError(t, "httpin: body format cannot be empty", func() {
		RegisterBodyFormat("", &yamlBody{}, true)
	})
}

func TestBodyEncoder_JSON(t *testing.T) {
	assert := assert.New(t)
	co, err := New(JSONBodyPayloadWithBodyDirective{})
	assert.NoError(err)
	req, err := co.NewRequest("POST", "/data", sampleObject_JSONBodyPayloadWithBodyDirective)
	expected, _ := http.NewRequest("POST", "/data", nil)
	expected.Form = url.Values{
		"page":      {"4"},
		"page_size": {"30"},
	}
	expected.Header.Set("Content-Type", "application/json")
	assert.NoError(err)
	var body bytes.Buffer
	assert.NoError(json.NewEncoder(&body).Encode(sampleObject_JSONBodyPayloadWithBodyDirective.Body))
	expected.Body = io.NopCloser(&body)
	assert.NoError(err)
	assertRequest(t, expected, req)

	// On the server side (decode).
	gotValue, err := co.Decode(req)
	assert.NoError(err)
	got, ok := gotValue.(*JSONBodyPayloadWithBodyDirective)
	assert.True(ok)
	assert.Equal(sampleObject_JSONBodyPayloadWithBodyDirective, got)
}

func TestBodyEncoder_XML(t *testing.T) {
	assert := assert.New(t)
	co, err := New(XMLBodyPayloadWithBodyDirective{})
	assert.NoError(err)
	req, err := co.NewRequest("POST", "/data", sampleObject_XMLBodyPayloadWithBodyDirective)
	expected, _ := http.NewRequest("POST", "/data", nil)
	expected.Header.Set("Content-Type", "application/xml")
	assert.NoError(err)
	var body bytes.Buffer
	assert.NoError(xml.NewEncoder(&body).Encode(sampleObject_XMLBodyPayloadWithBodyDirective.Body))
	expected.Body = io.NopCloser(&body)
	assert.NoError(err)
	assertRequest(t, expected, req)

	// On the server side (decode).
	gotValue, err := co.Decode(req)
	assert.NoError(err)
	got, ok := gotValue.(*XMLBodyPayloadWithBodyDirective)
	assert.True(ok)
	assert.Equal(sampleObject_XMLBodyPayloadWithBodyDirective, got)
}

func assertRequest(t *testing.T, expected, actual *http.Request) {
	assert := assert.New(t)
	assert.Equal(expected.Method, actual.Method)
	assert.Equal(expected.URL.Path, actual.URL.Path)
	assert.Equal(expected.URL.RawQuery, actual.URL.RawQuery)
	assert.Equal(expected.Header, actual.Header)
	assert.Equal(expected.Form, actual.Form)
	assert.Equal(expected.MultipartForm, actual.MultipartForm)
	assert.Equal(expected.PostForm, actual.PostForm)
	assert.Equal(expected.Cookies(), actual.Cookies())
	assert.Equal(expected.ContentLength, actual.ContentLength)

	if expected.Body == nil {
		assert.Nil(actual.Body)
	} else {
		expectedContent, err := io.ReadAll(expected.Body)
		assert.NoError(err)

		// Make a copy. The actual request may be used later to send request for an integration test.
		var bodyCopy bytes.Buffer
		actualContent, err := io.ReadAll(io.TeeReader(actual.Body, &bodyCopy))
		actual.Body = io.NopCloser(&bodyCopy) // replace the body that has been consumed
		assert.NoError(err)
		assert.Equal(expectedContent, actualContent)
	}
}
