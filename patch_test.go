package httpin

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/ggicci/httpin/patch"
	"github.com/stretchr/testify/assert"
)

type AccountPatch struct {
	Email  patch.Field[string] `in:"form=email"`
	Age    patch.Field[int]    `in:"form=age"`
	Avatar patch.Field[File]   `in:"form=avatar"`
}

func TestPatchField(t *testing.T) {
	fileContent := []byte("hello")
	r := newMultipartFormRequestFromMap(map[string]interface{}{
		"age":    "18",
		"avatar": fileContent,
	})

	core, err := New(AccountPatch{})
	assert.NoError(t, err)
	gotValue, err := core.Decode(r)
	assert.NoError(t, err)
	got := gotValue.(*AccountPatch)

	assert.Equal(t, patch.Field[string]{
		Valid: false,
		Value: "",
	}, got.Email)

	assert.Equal(t, patch.Field[int]{
		Valid: true,
		Value: 18,
	}, got.Age)

	assertFile(t, got.Avatar.Value, "avatar.txt", fileContent)
}

func TestPatchField_DecodeValueFailed(t *testing.T) {
	r, _ := http.NewRequest("GET", "/", nil)
	r.Form = url.Values{
		"email": {"abc@example.com"},
		"age":   {"eighteen"},
	}
	core, err := New(AccountPatch{})
	assert.NoError(t, err)
	gotValue, err := core.Decode(r)
	assert.Error(t, err)
	var ferr *InvalidFieldError
	assert.ErrorAs(t, err, &ferr)
	assert.Equal(t, "Age", ferr.Field)
	assert.Equal(t, "eighteen", ferr.Value)
	assert.Equal(t, "form", ferr.Source)
	assert.Nil(t, gotValue)
}

func TestPatchField_DecodeFileFailed(t *testing.T) {
	body, writer := newMultipartFormWriterFromMap(map[string]interface{}{
		"email":  "abc@example.com",
		"age":    "18",
		"avatar": []byte("hello"),
	})

	// break the boundary to make the file decoder fail
	r, _ := http.NewRequest("POST", "/", breakMultipartFormBoundary(body))
	r.Header.Set("Content-Type", writer.FormDataContentType())

	core, err := New(AccountPatch{})
	assert.NoError(t, err)
	gotValue, err := core.Decode(r)
	assert.Nil(t, gotValue)
	assert.Error(t, err)
}
