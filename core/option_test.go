package core

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithErrorHandler(t *testing.T) {
	// Use the default error handler.
	co, _ := New(ProductQuery{})
	assert.True(t, equalFuncs(globalCustomErrorHandler, co.GetErrorHandler()))

	// Override the default error handler.
	myErrorHandler := func(rw http.ResponseWriter, r *http.Request, err error) {}
	co, _ = New(ProductQuery{}, WithErrorHandler(myErrorHandler))
	assert.True(t, equalFuncs(myErrorHandler, co.GetErrorHandler()))

	// Fail on nil error handler.
	_, err := New(ProductQuery{}, WithErrorHandler(nil))
	assert.ErrorContains(t, err, "nil error handler")
}

func TestWithMaxMemory(t *testing.T) {
	// Use the default max memory.
	co, _ := New(ProductQuery{})
	assert.Equal(t, defaultMaxMemory, co.maxMemory)

	// Override the default max memory.
	co, _ = New(ProductQuery{}, WithMaxMemory(16<<20))
	assert.Equal(t, int64(16<<20), co.maxMemory)

	// Fail on too small max memory.
	_, err := New(ProductQuery{}, WithMaxMemory(100))
	assert.ErrorContains(t, err, "max memory too small")
}

func TestWithNestedDirectivesEnabled(t *testing.T) {
	// Override the default nested directives flag.
	co, _ := New(ProductQuery{}, WithNestedDirectivesEnabled(true))
	assert.Equal(t, true, co.enableNestedDirectives)
	co, _ = New(ProductQuery{}, WithNestedDirectivesEnabled(false))
	assert.Equal(t, false, co.enableNestedDirectives)
}

func TestEnableNestedDirectives(t *testing.T) {
	// Use the default nested directives flag.
	EnableNestedDirectives(false)
	co, _ := New(ProductQuery{})
	assert.Equal(t, false, co.enableNestedDirectives)

	EnableNestedDirectives(true)
	co, _ = New(ProductQuery{})
	assert.Equal(t, true, co.enableNestedDirectives)
}

func equalFuncs(expected, actual any) bool {
	return reflect.ValueOf(expected).Pointer() == reflect.ValueOf(actual).Pointer()
}
