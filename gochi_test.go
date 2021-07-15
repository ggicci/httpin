package httpin

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGochiURLParam(t *testing.T) {
	UseGochiURLParam("path", chi.URLParam) // register the "path" executor

	Convey("Gochi: can extract URLParam", t, func() {
		rw := httptest.NewRecorder()
		r, err := http.NewRequest("GET", "/ggicci/posts/1024", nil)
		So(err, ShouldBeNil)

		router := chi.NewRouter()
		router.With(NewInput(GetPostOfUserInput{})).Get("/{username}/posts/{pid}", GetPostOfUserHandler)
		router.ServeHTTP(rw, r)
		So(rw.Code, ShouldEqual, 200)
		expected := `{"Username":"ggicci","PostID":1024}` + "\n"
		So(rw.Body.String(), ShouldEqual, expected)
	})
}
