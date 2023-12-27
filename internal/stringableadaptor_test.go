package internal

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type YesNo bool

func (yn YesNo) ToString() (string, error) {
	if yn {
		return "yes", nil
	} else {
		return "no", nil
	}
}

func (yn *YesNo) FromString(s string) error {
	switch strings.ToLower(s) {
	case "yes":
		*yn = true
	case "no":
		*yn = false
	default:
		return errors.New("invalid value")
	}
	return nil
}

func TestToAnyStringableAdaptor(t *testing.T) {
	adaptor := NewAnyStringableAdaptor[bool](func(b *bool) (Stringable, error) {
		return (*YesNo)(b), nil
	})

	var validBoolean bool = true
	stringable, err := adaptor(&validBoolean)
	assert.NoError(t, err)
	v, err := stringable.ToString()
	assert.NoError(t, err)
	assert.Equal(t, "yes", v)
	assert.NoError(t, stringable.FromString("no"))
	assert.False(t, validBoolean)

	var invalidType int = 0
	stringable, err = adaptor(&invalidType)
	assert.ErrorIs(t, err, ErrTypeMismatch)
	assert.Nil(t, stringable)
	assert.ErrorContains(t, err, "cannot convert *int to *bool")
}
