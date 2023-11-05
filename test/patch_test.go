package test_test

import (
	"net/http"
	"net/url"
	"path/filepath"
	"testing"

	"github.com/ggicci/httpin"
	"github.com/ggicci/httpin/patch"
	"github.com/stretchr/testify/assert"
)

type AccountPatch struct {
	Email    patch.Field[string]         `in:"form=email"`
	Age      patch.Field[int]            `in:"form=age"`
	Avatar   patch.Field[*httpin.File]   `in:"form=avatar"`
	Hobbies  patch.Field[[]string]       `in:"form=hobbies"`
	Pictures patch.Field[[]*httpin.File] `in:"form=pictures"`
}

func TestPatchField(t *testing.T) {
	fileContent := []byte("hello")
	r := newMultipartFormRequestFromMap(map[string]any{
		"age":    "18",
		"avatar": fileContent,
		"hobbies": []string{
			"reading",
			"swimming",
		},
	})

	core, err := httpin.New(AccountPatch{})
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

	assert.Equal(t, patch.Field[[]string]{
		Valid: true,
		Value: []string{"reading", "swimming"},
	}, got.Hobbies)

	assert.Equal(t, patch.Field[[]*httpin.File]{
		Valid: false,
		Value: nil,
	}, got.Pictures)

	assertDecodedFile(t, got.Avatar.Value, "avatar.txt", fileContent)
}

func TestPatchField_DecodeValueFailed(t *testing.T) {
	r, _ := http.NewRequest("GET", "/", nil)
	r.Form = url.Values{
		"email": {"abc@example.com"},
		"age":   {"eighteen"},
	}
	core, err := httpin.New(AccountPatch{})
	assert.NoError(t, err)
	gotValue, err := core.Decode(r)
	assert.Error(t, err)
	var ferr *httpin.InvalidFieldError
	assert.ErrorAs(t, err, &ferr)
	assert.Equal(t, "Age", ferr.Field)
	assert.Equal(t, []string{"eighteen"}, ferr.Value)
	assert.Equal(t, "form", ferr.Source)
	assert.Nil(t, gotValue)
}

func TestPatchField_DecodeFileFailed(t *testing.T) {
	body, writer := newMultipartFormWriterFromMap(map[string]any{
		"email":  "abc@example.com",
		"age":    "18",
		"avatar": []byte("hello"),
	})

	// break the boundary to make the file decoder fail
	r, _ := http.NewRequest("POST", "/", breakMultipartFormBoundary(body))
	r.Header.Set("Content-Type", writer.FormDataContentType())

	core, err := httpin.New(AccountPatch{})
	assert.NoError(t, err)
	gotValue, err := core.Decode(r)
	assert.Nil(t, gotValue)
	assert.Error(t, err)
}

func TestPatchField_Encode(t *testing.T) {
	type ListQuery struct {
		Username patch.Field[string]   `in:"query=username"`
		Age      patch.Field[int]      `in:"query=age"`
		State    patch.Field[[]string] `in:"query=state[]"`
	}

	core, err := httpin.New(ListQuery{})
	assert.NoError(t, err)

	testcases := []struct {
		Query    *ListQuery
		Expected url.Values
	}{
		{&ListQuery{
			Username: patch.Field[string]{Value: "ggicci", Valid: true},
			Age:      patch.Field[int]{Value: 18, Valid: false},
		}, url.Values{"username": {"ggicci"}}},
		{&ListQuery{
			Age: patch.Field[int]{Value: 18, Valid: false},
		}, url.Values{}},
		{&ListQuery{
			Age: patch.Field[int]{Value: 18, Valid: true},
		}, url.Values{"age": {"18"}}},
		{&ListQuery{
			Username: patch.Field[string]{Value: "ggicci", Valid: true},
			Age:      patch.Field[int]{Value: 18, Valid: true},
			State: patch.Field[[]string]{
				Value: []string{"reading", "swimming"},
				Valid: true,
			},
		}, url.Values{
			"username": {"ggicci"},
			"age":      {"18"},
			"state[]":  {"reading", "swimming"},
		}},
	}

	for _, c := range testcases {
		req, err := core.NewRequest("GET", "/list", c.Query)
		assert.NoError(t, err)

		expected, _ := http.NewRequest("GET", "/list", nil)
		expected.URL.RawQuery = c.Expected.Encode()
		assert.Equal(t, expected, req)
	}
}

func TestPatchField_Encode_noFiles(t *testing.T) {
	assert := assert.New(t)
	payload := &AccountPatch{
		Email:  patch.Field[string]{Value: "abc@example.com", Valid: true},
		Age:    patch.Field[int]{Value: 18, Valid: true},
		Avatar: patch.Field[*httpin.File]{Value: nil, Valid: false},
		Hobbies: patch.Field[[]string]{
			Value: []string{"reading", "swimming"},
			Valid: true,
		},
		Pictures: patch.Field[[]*httpin.File]{Value: nil, Valid: false},
	}

	expected, _ := http.NewRequest("POST", "/patchAccount", nil)
	expected.Form = url.Values{
		"email":   {"abc@example.com"},
		"age":     {"18"},
		"hobbies": {"reading", "swimming"},
	}
	expected.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	core, err := httpin.New(AccountPatch{})
	assert.NoError(err)
	req, err := core.NewRequest("POST", "/patchAccount", payload)
	assert.NoError(err)
	assert.Equal(expected, req)
}

func TestPatchField_Encode_withFiles(t *testing.T) {
	assert := assert.New(t)
	avatarFile := createTempFile(t, []byte("handsome avatar image"))
	pic1Filename := createTempFile(t, []byte("pic1 content"))
	pic2Filename := createTempFile(t, []byte("pic2 content"))

	payload := &AccountPatch{
		Email:  patch.Field[string]{Value: "abc@example.com", Valid: true},
		Age:    patch.Field[int]{Value: 18, Valid: true},
		Avatar: patch.Field[*httpin.File]{Value: httpin.UploadWithFilename(avatarFile), Valid: true},
		Hobbies: patch.Field[[]string]{
			Value: []string{"reading", "swimming"},
			Valid: true,
		},
		Pictures: patch.Field[[]*httpin.File]{
			Value: []*httpin.File{
				httpin.UploadWithFilename(pic1Filename),
				httpin.UploadWithFilename(pic2Filename),
			},
			Valid: true,
		},
	}

	// See TestMultipartFormEncode_UploadFilename for more details.
	core, err := httpin.New(AccountPatch{})
	assert.NoError(err)
	req, err := core.NewRequest("POST", "/patchAccount", payload)
	assert.NoError(err)

	// Server side: receive files (decode).
	gotValue, err := core.Decode(req)
	assert.NoError(err)
	got, ok := gotValue.(*AccountPatch)
	assert.True(ok)
	assert.True(got.Email.Valid)
	assert.Equal("abc@example.com", got.Email.Value)
	assert.True(got.Age.Valid)
	assert.Equal(18, got.Age.Value)
	assert.True(got.Hobbies.Valid)
	assert.Equal([]string{"reading", "swimming"}, got.Hobbies.Value)
	assert.True(got.Avatar.Valid)
	assertDecodedFile(t, got.Avatar.Value, filepath.Base(avatarFile), []byte("handsome avatar image"))
	assert.True(got.Pictures.Valid)
	assert.Len(got.Pictures.Value, 2)
	assertDecodedFile(t, got.Pictures.Value[0], filepath.Base(pic1Filename), []byte("pic1 content"))
	assertDecodedFile(t, got.Pictures.Value[1], filepath.Base(pic2Filename), []byte("pic2 content"))
}
