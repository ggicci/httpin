package httpin

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/ggicci/patch"
	"github.com/stretchr/testify/assert"
)

func TestPatchField_Form(t *testing.T) {
	type AccountPatch struct {
		Email patch.Field[string] `in:"form=email"`
		Age   patch.Field[int]    `in:"form=age"`
	}

	r, _ := http.NewRequest("GET", "/", nil)
	r.Form = url.Values{
		"email": {},
		"age":   {"18"},
	}
	expected := &AccountPatch{
		Email: patch.Field[string]{
			Valid: false,
			Value: "",
		},
		Age: patch.Field[int]{
			Valid: true,
			Value: 18,
		},
	}
	core, err := New(AccountPatch{})
	assert.NoError(t, err)
	got, err := core.Decode(r)
	assert.NoError(t, err)
	assert.Equal(t, expected, got)
}
