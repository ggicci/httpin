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
	. "github.com/smartystreets/goconvey/convey"
)

type EchoInput struct {
	Token  string `in:"form=access_token;header=x-api-key;required"`
	Saying string `in:"form=saying"`
}

func EchoHandler(rw http.ResponseWriter, r *http.Request) {
	var input = r.Context().Value(Input).(*EchoInput)
	json.NewEncoder(rw).Encode(input)
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

func TestMiddleware(t *testing.T) {
	Convey("Should panic on invalid input", t, func() {
		So(func() { NewInput(nil) }, ShouldPanic)
	})

	Convey("Decode request successfully", t, func() {
		r, err := http.NewRequest("GET", "/", nil)
		So(err, ShouldBeNil)

		r.Header.Add("X-Api-Key", "abc")
		var params = url.Values{}
		params.Add("saying", "TO THINE OWE SELF BE TRUE")
		r.URL.RawQuery = params.Encode()

		rw := httptest.NewRecorder()
		handler := alice.New(NewInput(EchoInput{})).ThenFunc(EchoHandler)
		handler.ServeHTTP(rw, r)
		So(rw.Code, ShouldEqual, 200)
		expected := `{"Token":"abc","Saying":"TO THINE OWE SELF BE TRUE"}` + "\n"
		So(rw.Body.String(), ShouldEqual, expected)
	})

	Convey("Decode request failed with default error handler", t, func() {
		r, err := http.NewRequest("GET", "/", nil)
		So(err, ShouldBeNil)

		var params = url.Values{}
		params.Add("saying", "TO THINE OWE SELF BE TRUE")
		r.URL.RawQuery = params.Encode()

		rw := httptest.NewRecorder()
		handler := alice.New(NewInput(EchoInput{})).ThenFunc(EchoHandler)
		handler.ServeHTTP(rw, r)
		var out map[string]interface{}
		So(json.NewDecoder(rw.Body).Decode(&out), ShouldBeNil)
		So(out["field"], ShouldEqual, "Token")
		So(out["source"], ShouldEqual, "required")
		So(out["error"], ShouldEqual, ErrMissingField.Error())
	})

	Convey("Decode request failed with custom error handler", t, func() {
		r, err := http.NewRequest("GET", "/", nil)
		So(err, ShouldBeNil)

		var params = url.Values{}
		params.Add("saying", "TO THINE OWE SELF BE TRUE")
		r.URL.RawQuery = params.Encode()

		rw := httptest.NewRecorder()
		handler := alice.New(
			NewInput(EchoInput{}, WithErrorHandler(CustomErrorHandler)),
		).ThenFunc(EchoHandler)
		handler.ServeHTTP(rw, r)
		So(rw.Code, ShouldEqual, 400)
		So(rw.Body.String(), ShouldContainSubstring, `invalid field "Token":`)
	})
}

func TestReplaceDefaultErrorHandler(t *testing.T) {
	Convey("Given a nil handler should panic", t, func() {
		So(func() { ReplaceDefaultErrorHandler(nil) }, ShouldPanic)
	})

	Convey("Replace default error handler globally", t, func() {
		r, err := http.NewRequest("GET", "/", nil)
		So(err, ShouldBeNil)

		var params = url.Values{}
		params.Add("saying", "TO THINE OWE SELF BE TRUE")
		r.URL.RawQuery = params.Encode()
		rw := httptest.NewRecorder()
		handler := alice.New(NewInput(EchoInput{})).ThenFunc(EchoHandler)
		// NOTE: replace global error handler after NewInput should work
		ReplaceDefaultErrorHandler(CustomErrorHandler)

		handler.ServeHTTP(rw, r)
		So(rw.Code, ShouldEqual, 400)
	})
}
