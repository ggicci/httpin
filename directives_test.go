package httpin

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var noopDirective = DirectiveExecutorFunc(nil)

func TestRegisterDirectiveExecutor(t *testing.T) {
	assert.NotPanics(t, func() {
		RegisterDirectiveExecutor("noop_TestRegisterDirectiveExecutor", noopDirective)
	})

	assert.Panics(t, func() {
		RegisterDirectiveExecutor("noop_TestRegisterDirectiveExecutor", noopDirective)
	}, "should panic on duplicate name")

	assert.Panics(t, func() {
		RegisterDirectiveExecutor("nil_TestRegisterDirectiveExecutor", nil)
	}, "should panic on nil executor")

	assert.Panics(t, func() {
		RegisterDirectiveExecutor("decoder", noopDirective)
	}, "should panic on reserved name")
}

func TestReplaceDirectiveExecutor(t *testing.T) {
	assert.NotPanics(t, func() {
		ReplaceDirectiveExecutor("noop_TestReplaceDirectiveExecutor", noopDirective)
	})

	assert.NotPanics(t, func() {
		ReplaceDirectiveExecutor("noop_TestReplaceDirectiveExecutor", noopDirective)
	}, "should not panic on duplicate name")

	assert.Panics(t, func() {
		ReplaceDirectiveExecutor("nil_TestReplaceDirectiveExecutor", nil)
	}, "should panic on nil executor")

	assert.Panics(t, func() {
		ReplaceDirectiveExecutor("decoder", noopDirective)
	}, "should panic on reserved name")
}
