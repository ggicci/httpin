package integration_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ggicci/httpin"
	httpin_integration "github.com/ggicci/httpin/integration"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

type GetArticleOfUserInput struct {
	Author    string `in:"gochi=author"` // equivalent to chi.URLParam("author")
	ArticleID int64  `in:"gochi=articleID"`
}

func GetArticleOfUser(rw http.ResponseWriter, r *http.Request) {
	var input = r.Context().Value(httpin.Input).(*GetArticleOfUserInput)
	json.NewEncoder(rw).Encode(input)
}

func TestUseGochiURLParam(t *testing.T) {
	// Register the "gochi" directive, usually in init().
	// In most cases, you register this as "path", here's just an example.
	// Which is in order to avoid test conflicts with other tests
	httpin_integration.UseGochiURLParam("gochi", chi.URLParam)

	rw := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/ggicci/articles/1024", nil)
	assert.NoError(t, err)

	router := chi.NewRouter()
	router.With(
		httpin.NewInput(GetArticleOfUserInput{}),
	).Get("/{author}/articles/{articleID}", GetArticleOfUser)

	router.ServeHTTP(rw, r)
	assert.Equal(t, 200, rw.Code)
	expected := `{"Author":"ggicci","ArticleID":1024}` + "\n"
	assert.Equal(t, expected, rw.Body.String())
}
