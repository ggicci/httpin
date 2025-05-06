package integration_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ggicci/httpin"
	httpinIntegration "github.com/ggicci/httpin/integration"
	"github.com/stretchr/testify/assert"
)

type HttpMuxPathInput struct {
	Username string `in:"path=username"`
	PostID   int64  `in:"path=pid"`
}

func TestUseHttpMux(t *testing.T) {
	httpinIntegration.UseHttpMux("path")

	srv := http.NewServeMux()
	srv.HandleFunc("/users/{username}/posts/{pid}", func(w http.ResponseWriter, r *http.Request) {
		param := &HttpMuxPathInput{}
		core, err := httpin.New(param)
		if err != nil {
			t.Fatal(err)
			return
		}
		v, err := core.Decode(r)
		if err != nil {
			t.Fatal(err)
			return
		}
		jsonBytes, err := json.Marshal(v)
		if err != nil {
			t.Fatal(err)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(jsonBytes)
	})
	ts := httptest.NewServer(srv)
	defer ts.Close()

	resp, err := http.DefaultClient.Get(ts.URL + "/users/chriss-de/posts/456")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	bodyBytes, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	bodyString := string(bodyBytes)

	assert.Equal(t, `{"Username":"chriss-de","PostID":456}`, strings.TrimSpace(bodyString))
}
