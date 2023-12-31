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

func TestDirectivePath_ErrDefaultPathDecodingIsUnimplemented(t *testing.T) {
	type GetProfileRequest struct {
		Username string `in:"path=username"`
	}

	co, err := New(GetProfileRequest{})
	assert.NoError(t, err)
	req, _ := http.NewRequest("GET", "/users/ggicci", nil)
	_, err = co.Decode(req)
	assert.ErrorContains(t, err, "unimplemented path decoding function")
}

func TestDirectivePath_NewRequest_DefalutPathEncodingShouldWork(t *testing.T) {
	assert := assert.New(t)
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
