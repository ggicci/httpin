package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPriorityPair_ErrDuplicateType(t *testing.T) {
	pair := NewPriorityPair()
	err := pair.SetPair(TypeOf[int](), 1, 2, false)
	assert.NoError(t, err)
	assert.Equal(t, 1, pair.GetOne(TypeOf[int]()))

	err = pair.SetPair(TypeOf[int](), 3, 4, false)
	assert.ErrorContains(t, err, "duplicate type")
	assert.Equal(t, 1, pair.GetPrimary(TypeOf[int]()))
	assert.Equal(t, 2, pair.GetSecondary(TypeOf[int]()), "secondary value shouldn't be changed when there's a conflict")
	assert.Equal(t, 1, pair.GetOne(TypeOf[int]()))
}

func TestPriorityPair_ignoreConflict(t *testing.T) {
	pair := NewPriorityPair()
	err := pair.SetPair(TypeOf[int](), 1, 2, false)
	assert.NoError(t, err)
	assert.Equal(t, 1, pair.GetPrimary(TypeOf[int]()))
	assert.Equal(t, 2, pair.GetSecondary(TypeOf[int]()))
	assert.Equal(t, 1, pair.GetOne(TypeOf[int]()))

	err = pair.SetPair(TypeOf[int](), 3, 4, true)
	assert.NoError(t, err)
	assert.Equal(t, 3, pair.GetPrimary(TypeOf[int]()), "primary value can be changed if we ignore conflict")
	assert.Equal(t, 4, pair.GetSecondary(TypeOf[int]()), "secondary value can be changed as long as no error returns")
	assert.Equal(t, 3, pair.GetOne(TypeOf[int]()))
}

func TestPriorityPair_Get_withoutExistingKeys(t *testing.T) {
	pair := NewPriorityPair()
	typ := TypeOf[int]()
	assert.Nil(t, pair.GetOne(typ))
	assert.Nil(t, pair.GetPrimary(typ))
	assert.Nil(t, pair.GetSecondary(typ))
}

func TestPriorityPair_GetOne(t *testing.T) {
	pair := NewPriorityPair()
	typ := TypeOf[int]()
	pair.SetPair(typ, nil, 2, false)
	assert.Equal(t, 2, pair.GetOne(typ), "secondary value should be returned if primary is nil")
}
