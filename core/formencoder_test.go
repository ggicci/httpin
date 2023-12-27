package core

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormEncoder_FieldSetByFormerDirectives(t *testing.T) {
	type SearchQuery struct {
		AccessToken string `in:"query=access_token;header=x-api-key"`
	}

	co, err := New(&SearchQuery{})
	assert.NoError(t, err)

	req, err := co.NewRequest("GET", "/search", &SearchQuery{
		AccessToken: "123456",
	})
	assert.NoError(t, err)

	// The AccessToken field should be set by the query directive.
	expected, _ := http.NewRequest("GET", "/search?access_token=123456", nil)
	assert.Equal(t, expected, req)
}
