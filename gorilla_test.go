package httpin

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	. "github.com/smartystreets/goconvey/convey"
)

type GetPostOfUserInput struct {
	Username string `in:"path=username"`
	PostID   int64  `in:"path=pid"`
}

func GetPostOfUserHandler(rw http.ResponseWriter, r *http.Request) {
	var input = r.Context().Value(Input).(*GetPostOfUserInput)
	json.NewEncoder(rw).Encode(input)
}

func TestGorillaMuxVars(t *testing.T) {
	UseGorillaMux("path", mux.Vars) // register the "path" executor

	Convey("Gorilla: can extract mux vars", t, func() {
		rw := httptest.NewRecorder()
		r, err := http.NewRequest("GET", "/ggicci/posts/1024", nil)
		So(err, ShouldBeNil)

		router := mux.NewRouter()
		router.Handle("/{username}/posts/{pid}", alice.New(
			NewInput(GetPostOfUserInput{}),
		).ThenFunc(GetPostOfUserHandler)).Methods("GET")
		router.ServeHTTP(rw, r)
		So(rw.Code, ShouldEqual, 200)
		expected := `{"Username":"ggicci","PostID":1024}` + "\n"
		So(rw.Body.String(), ShouldEqual, expected)
	})
}
