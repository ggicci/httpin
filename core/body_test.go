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

type BodyPayloadInJSON struct {
	Page     int          `in:"form=page"`
	PageSize int          `in:"form=page_size"`
	Body     *BodyPayload `in:"body=json"`
}

type BodyPayloadInXML struct {
	Body *BodyPayload `in:"body=xml"`
}

var sampleBodyPayloadInJSONText = `
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

var sampleBodyPayloadInJSONObject = &BodyPayloadInJSON{
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

var sampleBodyPayloadInXMLText = `
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

var sampleBodyPayloadInXMLObject = &BodyPayloadInXML{
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

func TestBodyDirective_Decode_JSON(t *testing.T) {
	assert := assert.New(t)
	co, err := New(BodyPayloadInJSON{})
	assert.NoError(err)

	r, _ := http.NewRequest("GET", "https://example.com", nil)
	r.Form = make(url.Values)
	r.Form.Set("page", "4")
	r.Form.Set("page_size", "30")
	r.Body = io.NopCloser(strings.NewReader(sampleBodyPayloadInJSONText))
	r.Header.Set("Content-Type", "application/json")
	gotValue, err := co.Decode(r)
	assert.NoError(err)
	assert.Equal(sampleBodyPayloadInJSONObject, gotValue)
}

func TestBodyDirective_Decode_XML(t *testing.T) {
	assert := assert.New(t)
	co, err := New(BodyPayloadInXML{})
	assert.NoError(err)

	r, _ := http.NewRequest("GET", "https://example.com", nil)
	r.Body = io.NopCloser(strings.NewReader(sampleBodyPayloadInXMLText))
	r.Header.Set("Content-Type", "application/xml")

	gotValue, err := co.Decode(r)
	assert.NoError(err)
	assert.Equal(sampleBodyPayloadInXMLObject, gotValue)
}

func TestBodyDirective_Decode_ErrUnknownBodyFormat(t *testing.T) {
	type UnknownBodyFormatPayload struct {
		Body *BodyPayload `in:"body=yaml"`
	}

	co, err := New(UnknownBodyFormatPayload{})
	assert.NoError(t, err)
	req, _ := http.NewRequest("GET", "https://example.com", nil)
	req.Body = io.NopCloser(strings.NewReader(sampleBodyPayloadInJSONText))
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

func TestRegisterBodyFormat(t *testing.T) {
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
	unregisterBodyFormat("yaml")
}

func TestRegisterBodyFormat_ErrNilBodySerializer(t *testing.T) {
	assert.Panics(t, func() {
		RegisterBodyFormat("toml", nil)
	})
}

func TestRegisterBodyFormat_ForceRegister(t *testing.T) {
	assert.NotPanics(t, func() {
		RegisterBodyFormat("yaml", &yamlBody{}, true)
	})
	assert.NotPanics(t, func() {
		RegisterBodyFormat("yaml", &yamlBody{}, true)
	})
	unregisterBodyFormat("yaml")
}

func TestRegisterBodyFormat_ForceRegisterWithEmptyBodyFormat(t *testing.T) {
	assert.PanicsWithError(t, "httpin: body format cannot be empty", func() {
		RegisterBodyFormat("", &yamlBody{}, true)
	})
}

func TestBodyDirective_Encode_JSON(t *testing.T) {
	assert := assert.New(t)
	co, err := New(BodyPayloadInJSON{})
	assert.NoError(err)
	req, err := co.NewRequest("POST", "/data", sampleBodyPayloadInJSONObject)
	expected, _ := http.NewRequest("POST", "/data", nil)
	expected.Form = url.Values{
		"page":      {"4"},
		"page_size": {"30"},
	}
	expected.Header.Set("Content-Type", "application/json")
	assert.NoError(err)
	var body bytes.Buffer
	assert.NoError(json.NewEncoder(&body).Encode(sampleBodyPayloadInJSONObject.Body))
	expected.Body = io.NopCloser(&body)
	assert.Equal(expected, req)

	// On the server side (decode).
	gotValue, err := co.Decode(req)
	assert.NoError(err)
	got, ok := gotValue.(*BodyPayloadInJSON)
	assert.True(ok)
	assert.Equal(sampleBodyPayloadInJSONObject, got)
}

func TestBodyDirective_Encode_XML(t *testing.T) {
	assert := assert.New(t)
	co, err := New(BodyPayloadInXML{})
	assert.NoError(err)
	req, err := co.NewRequest("POST", "/data", sampleBodyPayloadInXMLObject)
	expected, _ := http.NewRequest("POST", "/data", nil)
	expected.Header.Set("Content-Type", "application/xml")
	assert.NoError(err)
	var body bytes.Buffer
	assert.NoError(xml.NewEncoder(&body).Encode(sampleBodyPayloadInXMLObject.Body))
	expected.Body = io.NopCloser(&body)
	assert.Equal(expected, req)

	// On the server side (decode).
	gotValue, err := co.Decode(req)
	assert.NoError(err)
	got, ok := gotValue.(*BodyPayloadInXML)
	assert.True(ok)
	assert.Equal(sampleBodyPayloadInXMLObject, got)
}

func TestBodyDirective_Encode_ErrUnknownBodyFormat(t *testing.T) {
	type UnknownBodyFormatPayload struct {
		Body *BodyPayload `in:"body=yaml"`
	}
	query := &UnknownBodyFormatPayload{
		Body: nil,
	}
	co, err := New(UnknownBodyFormatPayload{})
	assert.NoError(t, err)
	req, err := co.NewRequest("PUT", "/apples/10", query)
	assert.ErrorContains(t, err, "unknown body format: \"yaml\"")
	assert.Nil(t, req)
}

func unregisterBodyFormat(format string) {
	delete(bodyFormats, format)
}
