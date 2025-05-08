package integration_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ggicci/httpin"
	httpin_integration "github.com/ggicci/httpin/integration"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestUseEchoMux(t *testing.T) {
	e := echo.New()
	// NOTE: I removed the API UseEchoPathRouter because it introduces minimal benefit
	// but adds surface area and maintenance cost.
	httpin_integration.UseEchoRouter("path", e)

	req := httptest.NewRequest(http.MethodGet, "/users/ggicci/posts/123", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)

	handler := func(ctx echo.Context) error {
		param := &GetPostOfUserInput{}
		core, err := httpin.New(param)
		if err != nil {
			return err
		}
		v, err := core.Decode(ctx.Request())
		if err != nil {
			return err
		}
		fmt.Println(param)
		return c.JSON(http.StatusOK, v)
	}
	e.GET("/users/:username/posts/:pid", handler)
	err := handler(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, `{"Username":"ggicci","PostID":123}`, strings.TrimSpace(rec.Body.String()))
}
