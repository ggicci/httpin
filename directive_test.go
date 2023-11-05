package httpin

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegisterDirectiveExecutor(t *testing.T) {
	assert.NotPanics(t, func() {
		Customizer().RegisterDirective("noop_TestRegisterDirectiveExecutor", noopDirective)
	})

	assert.Panics(t, func() {
		Customizer().RegisterDirective("noop_TestRegisterDirectiveExecutor", noopDirective)
	}, "should panic on duplicate name")

	assert.Panics(t, func() {
		Customizer().RegisterDirective("nil_TestRegisterDirectiveExecutor", nil)
	}, "should panic on nil executor")

	assert.Panics(t, func() {
		Customizer().RegisterDirective("decoder", noopDirective)
	}, "should panic on reserved name")
}

func TestRegisterDirectiveExecutor_forceReplace(t *testing.T) {
	assert.NotPanics(t, func() {
		Customizer().RegisterDirective("noop_TestRegisterDirectiveExecutor_forceReplace", noopDirective, true)
	})

	assert.NotPanics(t, func() {
		Customizer().RegisterDirective("noop_TestRegisterDirectiveExecutor_forceReplace", noopDirective, true)
	}, "should not panic on duplicate name")

	assert.Panics(t, func() {
		Customizer().RegisterDirective("nil_TestRegisterDirectiveExecutor_forceReplace", nil, true)
	}, "should panic on nil executor")

	assert.Panics(t, func() {
		Customizer().RegisterDirective("decoder", noopDirective, true)
	}, "should panic on reserved name")
}
