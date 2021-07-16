package httpin

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	. "github.com/smartystreets/goconvey/convey"
)

type GetArticleOfUserInput struct {
	Author    string `in:"gochi=author"` // equivalent to chi.URLParam("author")
	ArticleID int64  `in:"gochi=articleID"`
}

func GetArticleOfUser(rw http.ResponseWriter, r *http.Request) {
	var input = r.Context().Value(Input).(*GetArticleOfUserInput)
	json.NewEncoder(rw).Encode(input)
}

func TestGochiURLParam(t *testing.T) {
	UseGochiURLParam("gochi", chi.URLParam) // register the "gochi" executor

	Convey("Gochi: can extract URLParam", t, func() {
		rw := httptest.NewRecorder()
		r, err := http.NewRequest("GET", "/ggicci/articles/1024", nil)
		So(err, ShouldBeNil)

		router := chi.NewRouter()
		router.With(
			NewInput(GetArticleOfUserInput{}),
		).Get("/{author}/articles/{articleID}", GetArticleOfUser)

		router.ServeHTTP(rw, r)
		So(rw.Code, ShouldEqual, 200)
		expected := `{"Author":"ggicci","ArticleID":1024}` + "\n"
		So(rw.Body.String(), ShouldEqual, expected)
	})
}
