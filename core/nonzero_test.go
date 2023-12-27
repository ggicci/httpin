package core

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

type NonzeroQuery struct {
	Name     string `in:"query=name;nonzero"`
	AgeRange []int  `in:"query=age;nonzero"`
}

func TestDirectiveNonzero_Decode(t *testing.T) {
	co, err := New(&NonzeroQuery{})
	assert.NoError(t, err)

	r, _ := http.NewRequest("GET", "/users", nil)
	r.URL.RawQuery = url.Values{
		"name": {"ggicci"},
		"age":  {"18", "999"},
	}.Encode()

	got, err := co.Decode(r)
	assert.NoError(t, err)
	assert.Equal(t, &NonzeroQuery{
		Name:     "ggicci",
		AgeRange: []int{18, 999},
	}, got.(*NonzeroQuery))
}

func TestDirectiveNonzero_Decode_ErrZeroValue(t *testing.T) {
	co, err := New(&NonzeroQuery{})
	assert.NoError(t, err)

	r, _ := http.NewRequest("GET", "/users", nil)
	r.URL.RawQuery = url.Values{
		"name": {"ggicci"},
	}.Encode()

	_, err = co.Decode(r)
	assert.Error(t, err)
	var invalidField *InvalidFieldError
	assert.ErrorAs(t, err, &invalidField)
	assert.Equal(t, "AgeRange", invalidField.Field)
	assert.Equal(t, "nonzero", invalidField.Directive)
	assert.Empty(t, invalidField.Key)
	assert.Nil(t, invalidField.Value)
}

func TestDirectiveNonzero_Decode_InNestedJSONBody_Issue49(t *testing.T) {
	type UpdateUserInput struct {
		Payload struct {
			Display string `json:"display" in:"nonzero"`
		} `in:"body=json"`
	}

	// NOTE: WithNestedDirectivesEnabled(true) is required to enable nested directives.
	co, err := New(&UpdateUserInput{}, WithNestedDirectivesEnabled(true))
	assert.NoError(t, err)

	r, _ := http.NewRequest("POST", "/users/1", nil)
	r.Header.Set("Content-Type", "application/json")
	r.Body = makeBodyReader(`{"display": ""}`)
	got, err := co.Decode(r)
	assert.Nil(t, got)
	assert.ErrorContains(t, err, "nonzero")
	var invalidField *InvalidFieldError
	assert.ErrorAs(t, err, &invalidField)
	assert.Equal(t, "Payload", invalidField.Field)
	assert.Equal(t, "nonzero", invalidField.Directive)
}

func TestDirectiveNonzero_NewRequest(t *testing.T) {
	co, err := New(&NonzeroQuery{})
	assert.NoError(t, err)

	expected, _ := http.NewRequest("GET", "/users", nil)
	expected.URL.RawQuery = url.Values{
		"name": {"ggicci"},
		"age":  {"18", "999"},
	}.Encode()

	req, err := co.NewRequest("GET", "/users", &NonzeroQuery{
		Name:     "ggicci",
		AgeRange: []int{18, 999},
	})
	assert.NoError(t, err)
	assert.Equal(t, expected, req)
}

func TestDirectiveNonzero_NewRequest_ErrZeroValue(t *testing.T) {
	co, err := New(&NonzeroQuery{})
	assert.NoError(t, err)

	_, err = co.NewRequest("GET", "/users", &NonzeroQuery{})
	assert.ErrorContains(t, err, "zero value")
	assert.ErrorContains(t, err, "Name")
	assert.ErrorContains(t, err, "AgeRange")
}
