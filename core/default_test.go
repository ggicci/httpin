package core

import (
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/ggicci/httpin/internal"
	"github.com/ggicci/httpin/patch"
	"github.com/stretchr/testify/assert"
)

func TestDirectiveDefault_Decode(t *testing.T) {
	type ThingWithDefaultValues struct {
		Page           int                   `in:"form=page;default=1"`
		PointerPage    *int                  `in:"form=pointer_page;default=1"`
		PatchPage      patch.Field[int]      `in:"form=patch_page;default=1"`
		PerPage        int                   `in:"form=per_page;default=20"`
		StateList      []string              `in:"form=state;default=pending,in_progress,failed"`
		PatchStateList patch.Field[[]string] `in:"form=patch_state;default=a,b,c"`
	}

	r, _ := http.NewRequest("GET", "/", nil)
	r.Form = url.Values{
		"page":         {"7"},
		"pointer_page": {"9"},
		"patch_page":   {"11"},
		"state":        {},
		"patch_state":  {},
	}
	expected := &ThingWithDefaultValues{
		Page:           7,
		PointerPage:    internal.Pointerize[int](9),
		PatchPage:      patch.Field[int]{Value: 11, Valid: true},
		PerPage:        20,
		StateList:      []string{"pending", "in_progress", "failed"},
		PatchStateList: patch.Field[[]string]{Value: []string{"a", "b", "c"}, Valid: true},
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

	// Second decode generates error
	xxx, err = co.Decode(r)
	assert.NoError(t, err)
	assert.Equal(t, expected, xxx)
}

func TestDirectiveDefault_NewRequest(t *testing.T) {
	type ListTicketRequest struct {
		Page    int      `in:"query=page;default=1"`
		PerPage int      `in:"query=per_page;default=20"`
		States  []string `in:"query=state;default=assigned,in_progress"`
	}

	co, err := New(ListTicketRequest{})
	assert.NoError(t, err)

	payload := &ListTicketRequest{
		Page: 2,
	}
	expected, _ := http.NewRequest("GET", "/tickets", nil)
	expected.URL.RawQuery = url.Values{
		"page":     {"2"},
		"per_page": {"20"},
		"state":    {"assigned", "in_progress"},
	}.Encode()
	req, err := co.NewRequest("GET", "/tickets", payload)
	assert.NoError(t, err)
	assert.Equal(t, expected, req)
}

func TestDirectiveDefault_NewRequest_WithNamedCoder(t *testing.T) {
	registerMyDate()
	type ListUsersRequest struct {
		Page             int       `in:"query=page;default=1"`
		PerPage          int       `in:"query=per_page;default=20"`
		RegistrationDate time.Time `in:"query=registration_date;default=2020-01-01;coder=mydate"`
	}

	co, err := New(ListUsersRequest{})
	assert.NoError(t, err)

	payload := &ListUsersRequest{
		Page:    2,
		PerPage: 10,
	}
	expected, _ := http.NewRequest("GET", "/users", nil)
	expected.URL.RawQuery = url.Values{
		"page":              {"2"},
		"per_page":          {"10"},
		"registration_date": {"2020-01-01"},
	}.Encode()
	req, err := co.NewRequest("GET", "/users", payload)
	assert.NoError(t, err)
	assert.Equal(t, expected, req)
	unregisterMyDate()
}
