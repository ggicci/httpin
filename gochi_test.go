package httpin

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	. "github.com/smartystreets/goconvey/convey"
)

type GetPostOfUserV2Input struct {
	Username string `in:"gochi=username"`
	PostID   int64  `in:"gochi=pid"`
}

func GetPostOfUserV2Handler(rw http.ResponseWriter, r *http.Request) {
	var input = r.Context().Value(Input).(*GetPostOfUserV2Input)
	json.NewEncoder(rw).Encode(input)
}

func TestGochiURLParam(t *testing.T) {
	UseGochiURLParam("gochi", chi.URLParam) // register the "gochi" executor

	Convey("Gochi: can extract URLParam", t, func() {
		rw := httptest.NewRecorder()
		r, err := http.NewRequest("GET", "/ggicci/posts/1024", nil)
		So(err, ShouldBeNil)

		router := chi.NewRouter()
		router.With(
			NewInput(GetPostOfUserV2Input{}),
		).Get("/{username}/posts/{pid}", GetPostOfUserV2Handler)

		router.ServeHTTP(rw, r)
		So(rw.Code, ShouldEqual, 200)
		expected := `{"Username":"ggicci","PostID":1024}` + "\n"
		So(rw.Body.String(), ShouldEqual, expected)
	})
}
