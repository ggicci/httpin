package httpin

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/justinas/alice"
	"github.com/stretchr/testify/assert"
)

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

func TestNewInput_Error_byDefaultErrorHandler(t *testing.T) {
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
	assert.Equal(t, "required", out["source"])
	assert.Contains(t, out["error"], "missing required field")
}

func CustomErrorHandler(rw http.ResponseWriter, r *http.Request, err error) {
	var invalidFieldError *InvalidFieldError
	if errors.As(err, &invalidFieldError) {
		rw.WriteHeader(http.StatusBadRequest) // status: 400
		io.WriteString(rw, invalidFieldError.Error())
		return
	}
	http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError) // status: 500
}

func TestNewInput_Error_byCustomErrorHandler(t *testing.T) {
	r, err := http.NewRequest("GET", "/", nil)
	assert.NoError(t, err)

	var params = url.Values{}
	params.Add("saying", "TO THINE OWE SELF BE TRUE")
	r.URL.RawQuery = params.Encode()

	rw := httptest.NewRecorder()
	handler := alice.New(
		NewInput(EchoInput{}, WithErrorHandler(CustomErrorHandler)),
	).ThenFunc(EchoHandler)
	handler.ServeHTTP(rw, r)
	assert.Equal(t, 400, rw.Code)
	assert.Contains(t, rw.Body.String(), `invalid field "Token":`)
}

func TestReplaceDefaultErrorHandler(t *testing.T) {
	// Nil handler should panic.
	assert.PanicsWithError(t, "httpin: nil error handler", func() {
		Customizer().RegisterErrorHandler(nil)
	})

	r, err := http.NewRequest("GET", "/", nil)
	assert.NoError(t, err)

	var params = url.Values{}
	params.Add("saying", "TO THINE OWE SELF BE TRUE")
	r.URL.RawQuery = params.Encode()
	rw := httptest.NewRecorder()
	handler := alice.New(NewInput(EchoInput{})).ThenFunc(EchoHandler)
	// NOTE: replace global error handler after NewInput should work
	Customizer().RegisterErrorHandler(CustomErrorHandler)

	handler.ServeHTTP(rw, r)
	assert.Equal(t, 400, rw.Code)
}
