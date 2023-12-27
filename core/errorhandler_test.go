package core

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func myCustomErrorHandler(rw http.ResponseWriter, r *http.Request, err error) {
	var invalidFieldError *InvalidFieldError
	if errors.As(err, &invalidFieldError) {
		rw.WriteHeader(http.StatusBadRequest) // status: 400
		io.WriteString(rw, invalidFieldError.Error())
		return
	}
	http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError) // status: 500
}

func TestRegisterErrorHandler(t *testing.T) {
	// Nil handler should panic.
	assert.PanicsWithError(t, "httpin: nil error handler", func() {
		RegisterErrorHandler(nil)
	})

	RegisterErrorHandler(myCustomErrorHandler)
	assert.True(t, equalFuncs(globalCustomErrorHandler, myCustomErrorHandler))
}

func TestDefaultErrorHandler(t *testing.T) {
	r, _ := http.NewRequest("GET", "/", nil)
	rw := httptest.NewRecorder()

	// When met InvalidFieldError, it should return 422.
	defaultErrorHandler(rw, r, &InvalidFieldError{err: assert.AnError, ErrorMessage: assert.AnError.Error()})
	assert.Equal(t, 422, rw.Code)

	// When met other errors, it should return 500.
	rw = httptest.NewRecorder()
	defaultErrorHandler(rw, r, assert.AnError)
	assert.Equal(t, 500, rw.Code)
}
