package httpin

import (
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

func TestBodyDecoder_JSON(t *testing.T) {
	type JSONBodyPayloadWithBodyDirective struct {
		Page     int          `in:"form=page"`
		PageSize int          `in:"form=page_size"`
		Body     *BodyPayload `in:"body=json"`
	}

	assert := assert.New(t)
	core, err := New(JSONBodyPayloadWithBodyDirective{})
	assert.NoError(err)

	r, _ := http.NewRequest("GET", "https://example.com", nil)
	r.Form = make(url.Values)
	r.Form.Set("page", "4")
	r.Form.Set("page_size", "30")
	r.Body = io.NopCloser(strings.NewReader(`
	{
		"name": "Elia",
		"is_native": false,
		"age": 14,
		"hobbies": ["Gaming", "Drawing"],
		"languages": [
			{"lang": "English", "level": 10},
			{"lang": "Japanese", "level": 3}
		]
	}`))
	r.Header.Set("Content-Type", "application/json")

	expected := &JSONBodyPayloadWithBodyDirective{
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

	gotValue, err := core.Decode(r)
	assert.NoError(err)
	assert.Equal(expected, gotValue)
}

func TestBodyDecoder_XML(t *testing.T) {
	type XMLBodyPayloadWithBodyDirective struct {
		Body *BodyPayload `in:"body=xml"`
	}

	assert := assert.New(t)
	core, err := New(XMLBodyPayloadWithBodyDirective{})
	assert.NoError(err)

	r, _ := http.NewRequest("GET", "https://example.com", nil)
	r.Body = io.NopCloser(strings.NewReader(`
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
	</BodyPayload>`))
	r.Header.Set("Content-Type", "application/xml")

	expected := &XMLBodyPayloadWithBodyDirective{
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
	gotValue, err := core.Decode(r)
	assert.NoError(err)
	assert.Equal(expected, gotValue)
}

func TestBodyDecoder_DefaultsToJSON(t *testing.T) {
	type Payload struct {
		Body *BodyPayload `in:"body"`
	}

	core, err := New(Payload{})
	assert.NoError(t, err)
	d := core.resolver.Lookup("Body").GetDirective("body")
	assert.Equal(t, "json", d.Argv[0])
}

func TestBodyDecoder_ErrUnknownBodyType(t *testing.T) {
	type UnknownBodyTypePayload struct {
		Body *BodyPayload `in:"body=yaml"`
	}

	core, err := New(UnknownBodyTypePayload{})
	assert.ErrorIs(t, err, ErrUnknownBodyType)
	assert.Nil(t, core)
}

type yamlBodyDecoder struct{}

var errYamlNotImplemented = errors.New("yaml not implemented")

func (de *yamlBodyDecoder) Decode(src io.Reader, dst interface{}) error {
	return errYamlNotImplemented // for test only
}

type YamlInput struct {
	Body map[string]interface{} `in:"body=yaml"`
}

func TestRegisterBodyDecoder(t *testing.T) {
	assert.NotPanics(t, func() {
		RegisterBodyDecoder("yaml", &yamlBodyDecoder{})
	})
	assert.Panics(t, func() {
		RegisterBodyDecoder("yaml", &yamlBodyDecoder{})
	})

	core, err := New(YamlInput{})
	assert.NoError(t, err)

	r, _ := http.NewRequest("GET", "https://example.com", nil)
	r.Body = io.NopCloser(strings.NewReader(`version: "3"`))

	gotValue, err := core.Decode(r)
	assert.ErrorIs(t, err, errYamlNotImplemented)
	assert.Nil(t, gotValue)
}

func TestRegisterBodyDecoder_forceReplace(t *testing.T) {
	assert.NotPanics(t, func() {
		RegisterBodyDecoder("yaml", &yamlBodyDecoder{}, true)
	})
	assert.NotPanics(t, func() {
		RegisterBodyDecoder("yaml", &yamlBodyDecoder{}, true)
	})
}

func TestRegisterBodyDecoder_forceReplace_withEmptyBodyType(t *testing.T) {
	assert.PanicsWithValue(t, "httpin: body type cannot be empty", func() {
		RegisterBodyDecoder("", &yamlBodyDecoder{}, true)
	})
}
