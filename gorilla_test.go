package httpin

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/stretchr/testify/assert"
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
	UseGorillaMux("path", mux.Vars) // register the "path" directive, usually in init()

	rw := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/ggicci/posts/1024", nil)
	assert.NoError(t, err)

	router := mux.NewRouter()
	router.Handle("/{username}/posts/{pid}", alice.New(
		NewInput(GetPostOfUserInput{}),
	).ThenFunc(GetPostOfUserHandler)).Methods("GET")
	router.ServeHTTP(rw, r)
	assert.Equal(t, 200, rw.Code)
	expected := `{"Username":"ggicci","PostID":1024}` + "\n"
	assert.Equal(t, expected, rw.Body.String())
}
