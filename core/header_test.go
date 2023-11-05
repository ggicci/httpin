package core

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDirectiveHeader(t *testing.T) {
	type SearchQuery struct {
		ApiUid   int    `in:"header=x-api-uid"`
		ApiToken string `in:"header=x-api-token"`
	}

	r, _ := http.NewRequest("GET", "/", nil)
	r.Header.Set("X-Api-Token", "some-secret-token")
	r.Header.Set("X-Api-Uid", "91241844")
	expected := &SearchQuery{
		ApiUid:   91241844,
		ApiToken: "some-secret-token",
	}
	co, err := New(SearchQuery{})
	assert.NoError(t, err)
	got, err := co.Decode(r)
	assert.NoError(t, err)
	assert.Equal(t, expected, got.(*SearchQuery))
}

func TestDirectiveHeader_Encode(t *testing.T) {
	type ApiQuery struct {
		ApiUid   int    `in:"header=x-api-uid"`
		ApiToken string `in:"header=X-Api-Token"`
	}

	query := &ApiQuery{
		ApiUid:   91241844,
		ApiToken: "some-secret-token",
	}

	co, err := New(ApiQuery{})
	assert.NoError(t, err)
	req, err := co.NewRequest("POST", "/api", query)
	assert.NoError(t, err)

	expected, _ := http.NewRequest("POST", "/api", nil)
	// NOTE: the key will be canonicalized
	expected.Header.Set("x-api-uid", "91241844")
	expected.Header.Set("X-Api-Token", "some-secret-token")
	assert.Equal(t, expected, req)
}
