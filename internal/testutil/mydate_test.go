package testutil

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMyDate_ToString(t *testing.T) {
	tm := MyDate(time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC))
	str, err := tm.ToString()
	assert.NoError(t, err)
	assert.Equal(t, "2023-10-01", str)
}

func TestMyDate_FromString(t *testing.T) {
	tm := MyDate{}
	err := tm.FromString("2023-10-01")
	assert.NoError(t, err)
	assert.Equal(t, MyDate(time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC)), tm)

	err = tm.FromString("invalid-date")
	assert.Error(t, err)
	var invalidDateErr *InvalidDate
	assert.ErrorAs(t, err, &invalidDateErr)
	assert.Equal(t, "invalid-date", invalidDateErr.Value)
}
