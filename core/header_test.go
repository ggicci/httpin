package core

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDirectiveHeader_Decode(t *testing.T) {
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

func TestDirectiveHeader_NewRequest(t *testing.T) {
	type ApiQuery struct {
		ApiUid   int     `in:"header=x-api-uid,omitempty"`
		ApiToken *string `in:"header=X-Api-Token,omitempty"`
	}

	t.Run("with all values", func(t *testing.T) {
		tk := "some-secret-token"
		query := &ApiQuery{
			ApiUid:   91241844,
			ApiToken: &tk,
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
	})

	t.Run("with empty value", func(t *testing.T) {
		query := &ApiQuery{
			ApiUid:   0,
			ApiToken: nil,
		}

		co, err := New(ApiQuery{})
		assert.NoError(t, err)
		req, err := co.NewRequest("POST", "/api", query)
		assert.NoError(t, err)

		expected, _ := http.NewRequest("POST", "/api", nil)
		assert.Equal(t, expected, req)

		_, ok := req.Header["X-Api-Uid"]
		assert.False(t, ok)

		_, ok = req.Header["X-Api-Token"]
		assert.False(t, ok)
	})
}
