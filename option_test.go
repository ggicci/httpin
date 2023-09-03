package httpin

import (
	"net/http"
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
	assert.ErrorIs(t, err, ErrNilErrorHandler)
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
	assert.ErrorIs(t, err, ErrMaxMemoryTooSmall)
}
