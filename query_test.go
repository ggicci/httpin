package httpin

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDirectiveQuery(t *testing.T) {
	type SearchQuery struct {
		Query      string `in:"query=q;required"`
		PageNumber int    `in:"query=p"`
		PageSize   int    `in:"query=page_size"`
	}

	r, _ := http.NewRequest("GET", "/?q=doggy&p=2&page_size=5", nil)
	expected := &SearchQuery{
		Query:      "doggy",
		PageNumber: 2,
		PageSize:   5,
	}

	core, err := New(SearchQuery{})
	assert.NoError(t, err)
	got, err := core.Decode(r)
	assert.NoError(t, err)
	assert.Equal(t, expected, got.(*SearchQuery))
}
