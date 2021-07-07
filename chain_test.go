package httpin

import (
	"encoding/json"
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

func TestChain(t *testing.T) {
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

	Convey("Decode request failed", t, func() {
		r, err := http.NewRequest("GET", "/", nil)
		So(err, ShouldBeNil)

		var params = url.Values{}
		params.Add("saying", "TO THINE OWE SELF BE TRUE")
		r.URL.RawQuery = params.Encode()

		rw := httptest.NewRecorder()
		handler := alice.New(
			NewInput(EchoInput{}, WithErrorStatusCode(400)),
		).ThenFunc(EchoHandler)
		handler.ServeHTTP(rw, r)
		So(rw.Code, ShouldEqual, 400)
		var out map[string]interface{}
		So(json.NewDecoder(rw.Body).Decode(&out), ShouldBeNil)
		So(out["field"], ShouldEqual, "Token")
		So(out["source"], ShouldEqual, "required")
		So(out["error"], ShouldEqual, ErrMissingField.Error())
	})
}
