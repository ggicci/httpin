package httpin

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithErrorHandler(t *testing.T) {
	// Use the default error handler.
	core, _ := New(ProductQuery{})
	assert.True(t, equalFuncs(globalCustomErrorHandler, core.getErrorHandler()))

	// Override the default error handler.
	myErrorHandler := func(rw http.ResponseWriter, r *http.Request, err error) {}
	core, _ = New(ProductQuery{}, WithErrorHandler(myErrorHandler))
	assert.True(t, equalFuncs(myErrorHandler, core.getErrorHandler()))

	// Fail on nil error handler.
	_, err := New(ProductQuery{}, WithErrorHandler(nil))
	assert.ErrorContains(t, err, "nil error handler")
}

func TestWithMaxMemory(t *testing.T) {
	// Use the default max memory.
	core, _ := New(ProductQuery{})
	assert.Equal(t, defaultMaxMemory, core.maxMemory)

	// Override the default max memory.
	core, _ = New(ProductQuery{}, WithMaxMemory(16<<20))
	assert.Equal(t, int64(16<<20), core.maxMemory)

	// Fail on too small max memory.
	_, err := New(ProductQuery{}, WithMaxMemory(100))
	assert.ErrorContains(t, err, "max memory too small")
}

func equalFuncs(expected, actual any) bool {
	return reflect.ValueOf(expected).Pointer() == reflect.ValueOf(actual).Pointer()
}
