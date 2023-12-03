package core

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func myCustomPathDecode(rtm *DirectiveRuntime) error {
	return assert.AnError
}

func TestDirectivePath_Decode(t *testing.T) {
	pathDirective := NewDirectivePath(myCustomPathDecode)
	assert.ErrorIs(t, pathDirective.Decode(nil), assert.AnError)
}

func TestDirectivePath_Encode(t *testing.T) {
	assert := assert.New(t)
	RegisterDirective("path", NewDirectivePath(myCustomPathDecode))
	type Repository struct {
		Name       string `json:"name"`
		Visibility string `json:"visibility"` // public, private, internal
		License    string `json:"license"`
	}
	type CreateRepositoryRequest struct {
		Owner   string      `in:"path=owner"`
		Payload *Repository `in:"body=json"`
	}

	query := &CreateRepositoryRequest{
		Owner: "ggicci",
		Payload: &Repository{
			Name:       "httpin",
			Visibility: "public",
			License:    "MIT",
		},
	}

	co, err := New(query)
	assert.NoError(err)
	req, err := co.NewRequest("POST", "/users/{owner}/repos", query)
	assert.NoError(err)

	expected, _ := http.NewRequest("POST", "/users/ggicci/repos", nil)
	expected.Header.Set("Content-Type", "application/json")
	var body bytes.Buffer
	assert.NoError(json.NewEncoder(&body).Encode(query.Payload))
	expected.Body = io.NopCloser(&body)
	assert.Equal(expected, req)
}
