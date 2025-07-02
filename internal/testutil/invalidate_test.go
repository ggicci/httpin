package testutil

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInvalidDate(t *testing.T) {
	underlyingErr := errors.New("world")
	err := &InvalidDate{Value: "hello", Err: underlyingErr}
	assert.Error(t, err)
	assert.ErrorContains(t, err, "invalid date:")
	assert.ErrorIs(t, err, underlyingErr)
}
