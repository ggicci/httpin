package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDirectiveNoop(t *testing.T) {
	var noop DirectiveNoop
	assert.NoError(t, noop.Encode(nil), "Encode should not return an error")
	assert.NoError(t, noop.Decode(nil), "Decode should not return an error")
}
