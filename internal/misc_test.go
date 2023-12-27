package internal

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsNil(t *testing.T) {
	assert.False(t, IsNil(reflect.ValueOf("hello")))
	assert.True(t, IsNil(reflect.ValueOf((*string)(nil))))
}

func TestPanicOnError(t *testing.T) {
	PanicOnError(nil)

	assert.PanicsWithError(t, "httpin: "+assert.AnError.Error(), func() {
		PanicOnError(assert.AnError)
	})
}

func TestTypeOf(t *testing.T) {
	assert.Equal(t, reflect.TypeOf(0), TypeOf[int]())
}

func TestPointerize(t *testing.T) {
	assert.Equal(t, 102, *Pointerize[int](102))
}
