package core

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/ggicci/httpin/patch"
	"github.com/stretchr/testify/assert"
)

func TestDirectiveDefault_Decode(t *testing.T) {
	type ThingWithDefaultValues struct {
		Page      int      `in:"form=page;default=1"`
		PerPage   int      `in:"form=per_page;default=20"`
		StateList []string `in:"form=state;default=pending,in_progress,failed"`
	}

	r, _ := http.NewRequest("GET", "/", nil)
	r.Form = url.Values{
		"page":  {"7"},
		"state": {},
	}
	expected := &ThingWithDefaultValues{
		Page:      7,
		PerPage:   20,
		StateList: []string{"pending", "in_progress", "failed"},
	}
	co, err := New(ThingWithDefaultValues{})
	assert.NoError(t, err)
	got, err := co.Decode(r)
	assert.NoError(t, err)
	assert.Equal(t, expected, got)
}

func TestDirectiveDefault_Decode_PointerTypeFields(t *testing.T) {
	assert := assert.New(t)
	type Input struct {
		Page      *int     `in:"form=page;default=1"`
		PerPage   *int     `in:"form=per_page;default=20"`
		StateList []string `in:"form=state;default=pending,in_progress,failed"`
	}
	co, err := New(Input{})
	assert.NoError(err)

	r := newMultipartFormRequestFromMap(map[string]any{})
	gotValue, err := co.Decode(r)
	assert.NoError(err)
	got := gotValue.(*Input)
	assert.Equal(1, *got.Page)
	assert.Equal(20, *got.PerPage)
	assert.Equal([]string{"pending", "in_progress", "failed"}, got.StateList)
}

func TestDirectiveDefault_Decode_PatchField(t *testing.T) {
	type ThingWithDefaultValues struct {
		Page      patch.Field[int]      `in:"form=page;default=1"`
		PerPage   patch.Field[int]      `in:"form=per_page;default=20"`
		StateList patch.Field[[]string] `in:"form=state;default=pending,in_progress,failed"`
	}

	r := newMultipartFormRequestFromMap(map[string]any{
		"page": "7",
	})
	expected := &ThingWithDefaultValues{
		Page:      patch.Field[int]{Value: 7, Valid: true},
		PerPage:   patch.Field[int]{Value: 20, Valid: true},
		StateList: patch.Field[[]string]{Value: []string{"pending", "in_progress", "failed"}, Valid: true},
	}
	co, err := New(ThingWithDefaultValues{})
	assert.NoError(t, err)
	got, err := co.Decode(r)
	assert.NoError(t, err)
	assert.Equal(t, expected, got)
}

// FIX: https://github.com/ggicci/httpin/issues/77
// Decode parameter struct with default values only works the first time
func TestDirectiveDeafult_Decode_DecodeTwice(t *testing.T) {
	type ThingWithDefaultValues struct {
		Id      uint `in:"query=id;required"`
		Page    int  `in:"query=page;default=1"`
		PerPage int  `in:"query=page_size;default=127"`
	}

	r, _ := http.NewRequest("GET", "/?id=123", nil)
	expected := &ThingWithDefaultValues{
		Id:      123,
		Page:    1,
		PerPage: 127,
	}

	co, err := New(ThingWithDefaultValues{})
	assert.NoError(t, err)

	// First decode works as expected
	xxx, err := co.Decode(r)
	assert.NoError(t, err)
	assert.Equal(t, expected, xxx)

	// Second decode generates eror
	xxx, err = co.Decode(r)
	assert.NoError(t, err)
	assert.Equal(t, expected, xxx)
}
