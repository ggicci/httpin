package httpin

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegisterDirectiveExecutor(t *testing.T) {
	assert.NotPanics(t, func() {
		RegisterDirectiveExecutor("noop_TestRegisterDirectiveExecutor", noopDirective, noopDirective)
	})

	assert.Panics(t, func() {
		RegisterDirectiveExecutor("noop_TestRegisterDirectiveExecutor", noopDirective, noopDirective)
	}, "should panic on duplicate name")

	assert.Panics(t, func() {
		RegisterDirectiveExecutor("nil_TestRegisterDirectiveExecutor", nil, nil)
	}, "should panic on nil executor")

	assert.Panics(t, func() {
		RegisterDirectiveExecutor("decoder", noopDirective, noopDirective)
	}, "should panic on reserved name")
}

func TestRegisterDirectiveExecutor_forceReplace(t *testing.T) {
	assert.NotPanics(t, func() {
		RegisterDirectiveExecutor("noop_TestRegisterDirectiveExecutor_forceReplace", noopDirective, noopDirective, true)
	})

	assert.NotPanics(t, func() {
		RegisterDirectiveExecutor("noop_TestRegisterDirectiveExecutor_forceReplace", noopDirective, noopDirective, true)
	}, "should not panic on duplicate name")

	assert.Panics(t, func() {
		RegisterDirectiveExecutor("nil_TestRegisterDirectiveExecutor_forceReplace", nil, nil, true)
	}, "should panic on nil executor")

	assert.Panics(t, func() {
		RegisterDirectiveExecutor("decoder", noopDirective, noopDirective, true)
	}, "should panic on reserved name")
}
