package httpin

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/ggicci/httpin/core"
	"github.com/justinas/alice"
	"github.com/stretchr/testify/assert"
)

type Pagination struct {
	Page    int `in:"form=page,page_index,index"`
	PerPage int `in:"form=per_page,page_size"`
}

func TestDecode(t *testing.T) {
	r, _ := http.NewRequest("GET", "/", nil)
	r.Form = url.Values{
		"page":     {"1"},
		"per_page": {"100"},
	}
	expected := &Pagination{
		Page:    1,
		PerPage: 100,
	}

	func() {
		input := &Pagination{}
		err := Decode(r, input) // pointer to a struct instance
		assert.NoError(t, err)
		assert.Equal(t, expected, input)
	}()

	func() {
		input := Pagination{}
		err := Decode(r, &input) // addressable struct instance
		assert.NoError(t, err)
		assert.Equal(t, expected, &input)
	}()

	func() {
		input := &Pagination{}
		err := Decode(r, &input) // pointer to pointer of struct instance
		assert.NoError(t, err)
		assert.Equal(t, expected, input)
	}()

	func() {
		input := Pagination{}
		err := Decode(r, input) // non-pointer struct instance should fail
		assert.ErrorContains(t, err, "input must be a pointer")
	}()
}

func TestDecode_ErrBuildResolverFailed(t *testing.T) {
	r, _ := http.NewRequest("GET", "/", nil)
	r.Form = url.Values{
		"page":     {"1"},
		"per_page": {"100"},
	}

	type Foo struct {
		Name string `in:"nonexistent=foo"`
	}

	assert.Error(t, Decode(r, &Foo{}))
}

func TestDecode_ErrDecodeFailure(t *testing.T) {
	r, _ := http.NewRequest("GET", "/", nil)
	r.Form = url.Values{
		"page":     {"1"},
		"per_page": {"one-hundred"},
	}

	p := &Pagination{}
	assert.Error(t, Decode(r, p))
}

type EchoInput struct {
	Token  string `in:"form=access_token;header=x-api-key;required"`
	Saying string `in:"form=saying"`
}

func EchoHandler(rw http.ResponseWriter, r *http.Request) {
	var input = r.Context().Value(Input).(*EchoInput)
	json.NewEncoder(rw).Encode(input)
}

func TestNewInput_WithNil(t *testing.T) {
	assert.Panics(t, func() {
		NewInput(nil)
	})
}

func TestNewInput_Success(t *testing.T) {
	r, err := http.NewRequest("GET", "/", nil)
	assert.NoError(t, err)

	r.Header.Add("X-Api-Key", "abc")
	var params = url.Values{}
	params.Add("saying", "TO THINE OWE SELF BE TRUE")
	r.URL.RawQuery = params.Encode()

	rw := httptest.NewRecorder()
	handler := alice.New(NewInput(EchoInput{})).ThenFunc(EchoHandler)
	handler.ServeHTTP(rw, r)
	assert.Equal(t, 200, rw.Code)
	expected := `{"Token":"abc","Saying":"TO THINE OWE SELF BE TRUE"}` + "\n"
	assert.Equal(t, expected, rw.Body.String())
}

func TestNewInput_ErrorHandledByDefaultErrorHandler(t *testing.T) {
	r, err := http.NewRequest("GET", "/", nil)
	assert.NoError(t, err)

	var params = url.Values{}
	params.Add("saying", "TO THINE OWE SELF BE TRUE")
	r.URL.RawQuery = params.Encode()

	rw := httptest.NewRecorder()
	handler := alice.New(NewInput(EchoInput{})).ThenFunc(EchoHandler)
	handler.ServeHTTP(rw, r)
	var out map[string]any
	assert.Nil(t, json.NewDecoder(rw.Body).Decode(&out))

	assert.Equal(t, 422, rw.Code)
	assert.Equal(t, "Token", out["field"])
	assert.Equal(t, "required", out["directive"])
	assert.Contains(t, out["error"], "missing required field")
}

func CustomErrorHandler(rw http.ResponseWriter, r *http.Request, err error) {
	var invalidFieldError *core.InvalidFieldError
	if errors.As(err, &invalidFieldError) {
		rw.WriteHeader(http.StatusBadRequest) // status: 400
		io.WriteString(rw, invalidFieldError.Error())
		return
	}
	http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError) // status: 500
}

func TestNewRequest(t *testing.T) {
	req, err := NewRequest("GET", "/products", &Pagination{
		Page:    19,
		PerPage: 50,
	})
	assert.NoError(t, err)

	expected, _ := http.NewRequest("GET", "/products", nil)
	expected.Body = io.NopCloser(strings.NewReader("page=19&per_page=50"))
	expected.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	assert.Equal(t, expected, req)
}

func TestNewRequest_ErrNewFailure(t *testing.T) {
	_, err := NewRequest("GET", "/products", 123)
	assert.Error(t, err)
}
