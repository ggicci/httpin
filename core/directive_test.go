package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegisterDirectiveExecutor(t *testing.T) {
	assert.NotPanics(t, func() {
		RegisterDirective("noop_TestRegisterDirectiveExecutor", &DirectiveNoop{})
	})

	assert.Panics(t, func() {
		RegisterDirective("noop_TestRegisterDirectiveExecutor", &DirectiveNoop{})
	}, "should panic on duplicate name")

	assert.Panics(t, func() {
		RegisterDirective("nil_TestRegisterDirectiveExecutor", nil)
	}, "should panic on nil executor")

	assert.Panics(t, func() {
		RegisterDirective("decoder", &DirectiveNoop{})
	}, "should panic on reserved name")
}

func TestRegisterDirectiveExecutor_ForceReplace(t *testing.T) {
	assert.NotPanics(t, func() {
		RegisterDirective("noop_TestRegisterDirectiveExecutor_forceReplace", &DirectiveNoop{}, true)
	})

	assert.NotPanics(t, func() {
		RegisterDirective("noop_TestRegisterDirectiveExecutor_forceReplace", &DirectiveNoop{}, true)
	}, "should not panic on duplicate name")

	assert.Panics(t, func() {
		RegisterDirective("nil_TestRegisterDirectiveExecutor_forceReplace", nil, true)
	}, "should panic on nil executor")

	assert.Panics(t, func() {
		RegisterDirective("decoder", &DirectiveNoop{}, true)
	}, "should panic on reserved name")
}
