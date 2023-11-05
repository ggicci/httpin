package core

import (
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type RequiredQuery struct {
	CreatedAt time.Time `in:"form=created_at;required"`
	Color     string    `in:"form=colour,color"`
}

func TestDirectiveRequired_RequiredFieldMissing(t *testing.T) {
	r, _ := http.NewRequest("GET", "/", nil)
	r.Form = url.Values{
		"color": {"red"},
	}
	co, err := New(&RequiredQuery{})
	assert.NoError(t, err)
	_, err = co.Decode(r)
	assert.ErrorContains(t, err, "missing required field")
	var invalidField *InvalidFieldError
	assert.ErrorAs(t, err, &invalidField)
	assert.Equal(t, "required", invalidField.Source)
	assert.Empty(t, invalidField.Key)
	assert.Nil(t, invalidField.Value)
}

func TestDirectiveRequired_NonRequiredFieldAbsent(t *testing.T) {
	r, _ := http.NewRequest("GET", "/", nil)
	r.Form = url.Values{
		"created_at": {"1991-11-10T08:00:00+08:00"},
		"is_soldout": {"true"},
		"page":       {"1"},
		"per_page":   {"20"},
	}
	expected := &RequiredQuery{
		CreatedAt: time.Date(1991, 11, 10, 0, 0, 0, 0, time.UTC),
		Color:     "",
	}
	co, err := New(RequiredQuery{})
	assert.NoError(t, err)
	got, err := co.Decode(r)
	assert.NoError(t, err)
	assert.Equal(t, expected, got.(*RequiredQuery))
}
